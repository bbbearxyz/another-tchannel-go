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

package main

import (
	"github.com/bbbearxyz/another-tchannel-go/examples/protobuf/test"
	"log"
	"net"
	"time"

	"github.com/bbbearxyz/another-tchannel-go"
	"github.com/bbbearxyz/another-tchannel-go/pb"
)

func main() {
	var (
		listener net.Listener
		err      error
	)

	if listener, err = setupServer(); err != nil {
		log.Fatalf("setupServer failed: %v", err)
	}

	if err := runClient1("server", listener.Addr()); err != nil {
		log.Fatalf("runClient1 failed: %v", err)
	}

	if err := runClient2("server", listener.Addr()); err != nil {
		log.Fatalf("runClient1 failed: %v", err)
	}

	// Run for 10 seconds, then stop
	time.Sleep(time.Second * 10)
}

func setupServer() (net.Listener, error) {
	tchan, err := tchannel.NewChannel("server", nil)
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		return nil, err
	}

	server := pb.NewServer(tchan)
	server.Register(test.NewEchoServer(&echoHandler{}))

	// Serve will set the local peer info, and start accepting sockets in a separate goroutine.
	tchan.Serve(listener)
	return listener, nil
}

func runClient1(service string, addr net.Addr) error {
	tchan, err := tchannel.NewChannel("client1", nil)
	if err != nil {
		return err
	}
	tchan.Peers().Add(addr.String())
	tclient := pb.NewClient(tchan, service, nil)
	client := test.NewEchoClient(tclient)

	ctx, cancel := pb.NewContext(time.Second)
	resp, err := client.Send(ctx, &test.Request{Field1: "xxx"})
	if err != nil {
		println(err.Error())
	}
	println(resp.Msg)
	cancel()

	return nil
}

func runClient2(service string, addr net.Addr) error {
	tchan, err := tchannel.NewChannel("client2", nil)
	if err != nil {
		return err
	}
	tchan.Peers().Add(addr.String())
	tclient := pb.NewClient(tchan, service, nil)
	client := test.NewEchoClient(tclient)

	ctx, cancel := pb.NewContext(time.Minute)
	server, err := client.StreamTest(ctx)
	if err != nil {
		println(err.Error())
	}
	server.Send(&test.Request{Field1: "from to client"})
	resp1, _ := server.Recv()
	println(resp1.Msg)
	resp2, _ := server.Recv()
	println(resp2.Msg)
	server.Close()
	cancel()
	return nil
}

type echoHandler struct{
}

func (h *echoHandler) StreamTest(server test.Echo_StreamTest_Server) error {
	server.Send(&test.Response{Msg: "hello, world1"})
	server.Send(&test.Response{Msg: "hello, world2"})
	req, _ := server.Recv()
	println(req.Field1)
	server.Close()
	return nil
}

func (h *echoHandler) Send(ctx pb.Context, arg *test.Request) (*test.Response, error) {
	resp := test.Response{
		Msg: "yyy",
	}
	return &resp, nil
}

