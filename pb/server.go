// Copyright (c) 2015 Uber Technologies, Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package pb

import (
	"log"
	"strings"
	"sync"

	tchannel "github.com/bbbearxyz/another-tchannel-go"
	"github.com/bbbearxyz/another-tchannel-go/internal/argreader"

	"golang.org/x/net/context"
)

type handler struct {
	server         TChanServer
	postResponseCB PostResponseCB
}

// Server handles incoming TChannel calls and forwards them to the matching TChanServer.
type Server struct {
	sync.RWMutex
	ch          tchannel.Registrar
	log         tchannel.Logger
	handlers    map[string]handler
	ctxFn       func(ctx context.Context, method string, headers map[string]string) Context
	//cal unary or streaming
	typesMap    map[string]bool
}

// NewServer returns a server that can serve thrift services over TChannel.
func NewServer(registrar tchannel.Registrar) *Server {
	// 目前不支持meta handler
	server := &Server{
		ch:          registrar,
		log:         registrar.Logger(),
		handlers:    make(map[string]handler),
		ctxFn:       defaultContextFn,
		typesMap:    make(map[string]bool),
	}
	if ch, ok := registrar.(*tchannel.Channel); ok {
		// Register the meta endpoints on the "tchannel" service name.
		NewServer(ch.GetSubChannel("tchannel"))
	}
	return server
}

// Register registers the given TChanServer to be called on any incoming call for its' services.
// TODO(prashant): Replace Register call with this call.
func (s *Server) Register(svr TChanServer, opts ...RegisterOption) {
	types := svr.Types()
	service := svr.Service()
	handler := &handler{server: svr}
	for _, opt := range opts {
		opt.Apply(handler)
	}

	s.Lock()
	s.handlers[service] = *handler
	s.Unlock()

	for index, m := range svr.Methods() {
		name := service+"::"+m
		s.ch.Register(s, name)
		s.typesMap[name] = types[index]
	}
}

// SetContextFn sets the function used to convert a context.Context to a thrift.Context.
// Note: This API may change and is only intended to bridge different contexts.
func (s *Server) SetContextFn(f func(ctx context.Context, method string, headers map[string]string) Context) {
	s.ctxFn = f
}

func (s *Server) onError(call *tchannel.InboundCall, err error) {
	// TODO(prashant): Expose incoming call errors through options for NewServer.
	remotePeer := call.RemotePeer()
	logger := s.log.WithFields(
		tchannel.ErrField(err),
		tchannel.LogField{Key: "method", Value: call.MethodString()},
		tchannel.LogField{Key: "callerName", Value: call.CallerName()},

		// TODO: These are very similar to the connection fields, but we don't
		// have access to the connection's logger. Consider exposing the
		// connection through CurrentCall.
		tchannel.LogField{Key: "localAddr", Value: call.LocalPeer().HostPort},
		tchannel.LogField{Key: "remoteHostPort", Value: remotePeer.HostPort},
		tchannel.LogField{Key: "remoteIsEphemeral", Value: remotePeer.IsEphemeral},
		tchannel.LogField{Key: "remoteProcess", Value: remotePeer.ProcessName},
	)

	if tchannel.GetSystemErrorCode(err) == tchannel.ErrCodeTimeout {
		logger.Debug("Thrift server timeout.")
	} else {
		logger.Error("Thrift server error.")
	}
}

func defaultContextFn(ctx context.Context, method string, headers map[string]string) Context {
	return WithHeaders(ctx, headers)
}

func (s *Server) handleUnary(origCtx context.Context, handler handler, method string, call *tchannel.InboundCall) error {
	// 目前header没有任何参数
	reader, err := call.Arg2Reader()
	if err != nil {
		return err
	}
	headers, err := ReadHeaders(reader)
	if err != nil {
		return err
	}

	if err := argreader.EnsureEmpty(reader, "reading request headers"); err != nil {
		return err
	}

	if err := reader.Close(); err != nil {
		return err
	}

	reader, err = call.Arg3Reader()
	if err != nil {
		return err
	}

	tracer := tchannel.TracerFromRegistrar(s.ch)
	origCtx = tchannel.ExtractInboundSpan(origCtx, call, headers, tracer)
	ctx := s.ctxFn(origCtx, method, headers)

	success, resp, err := handler.server.Handle(ctx, method, reader, nil)
	if handler.postResponseCB != nil {
		defer handler.postResponseCB(ctx, method, resp)
	}

	if err := argreader.EnsureEmpty(reader, "reading request body"); err != nil {
		return err
	}
	if err := reader.Close(); err != nil {
		return err
	}

	if !success {
		call.Response().SetApplicationError()
	}

	writer, err := call.Response().Arg2Writer()
	if err != nil {
		return err
	}

	// 目前header不写入任何东西 只有一个len
	if err := WriteHeaders(writer, ctx.ResponseHeaders()); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	writer, err = call.Response().Arg3Writer()
	if err := WriteStruct(writer, resp); err != nil {
		call.Response().SendSystemError(err)
		return err
	}

	return writer.Close()
}

func (s *Server) handleStreaming(origCtx context.Context, handler handler, method string, call *tchannel.InboundCall) error {

	var arg2 []byte
	tchannel.NewArgReader(call.Arg2Reader()).Read(&arg2)

	tchannel.NewArgWriter(call.Response().Arg2Writer()).Write(nil)

	reader, err := call.Arg3Reader()
	if err != nil {
		return err
	}

	writer, _ := call.Response().Arg3Writer()

	success, _, err := handler.server.Handle(Wrap(origCtx), method, reader, writer)
	if !success {
		call.Response().SetApplicationError()
	}
	return err
}


// unary service::method
func getServiceMethod(method string) (string, string, bool) {
	s := string(method)
	sep := strings.Index(s, "::")
	if sep == -1 {
		return "", "", false
	}
	return s[:sep], s[sep+2:], true
}

// Handle handles an incoming TChannel call and forwards it to the correct handler.
func (s *Server) Handle(ctx context.Context, call *tchannel.InboundCall) {
	op := call.MethodString()

	service, method, ok := getServiceMethod(op)
	if !ok {
		log.Fatalf("Handle got call for %s which does not match the expected call format", op)
	}

	s.RLock()
	handler, ok := s.handlers[service]
	s.RUnlock()
	if !ok {
		log.Fatalf("Handle got call for service %v which is not registered", service)
	}

	// 根据method类型分发handleUnary还是handleStreaming
	// true means unary, false means stream
	if s.typesMap[op] {
		if err := s.handleUnary(ctx, handler, method, call); err != nil {
			s.onError(call, err)
		}
	} else {
		if err := s.handleStreaming(ctx, handler, method, call); err != nil {
			s.onError(call, err)
		}
	}
}
