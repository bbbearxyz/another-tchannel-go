// @generated Code generated by thrift-gen. Do not modify.

// Package keyvalue is generated code used to make or handle TChannel calls using Thrift.
package keyvalue

import (
	"fmt"

	athrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/bbbearxyz/another-tchannel-go/thrift"
)

// Interfaces for the service and client for the services defined in the IDL.

// TChanAdmin is the interface that defines the server handler and client interface.
type TChanAdmin interface {
	TChanBaseService

	ClearAll(ctx thrift.Context) error
}

// TChanKeyValue is the interface that defines the server handler and client interface.
type TChanKeyValue interface {
	TChanBaseService

	Get(ctx thrift.Context, key string) (string, error)
	Set(ctx thrift.Context, key string, value string) error
}

// TChanBaseService is the interface that defines the server handler and client interface.
type TChanBaseService interface {
	HealthCheck(ctx thrift.Context) (string, error)
}

// Implementation of a client and service handler.

type tchanAdminClient struct {
	TChanBaseService

	thriftService string
	client        thrift.TChanClient
}

func NewTChanAdminInheritedClient(thriftService string, client thrift.TChanClient) *tchanAdminClient {
	return &tchanAdminClient{
		NewTChanBaseServiceInheritedClient(thriftService, client),
		thriftService,
		client,
	}
}

// NewTChanAdminClient creates a client that can be used to make remote calls.
func NewTChanAdminClient(client thrift.TChanClient) TChanAdmin {
	return NewTChanAdminInheritedClient("Admin", client)
}

func (c *tchanAdminClient) ClearAll(ctx thrift.Context) error {
	var resp AdminClearAllResult
	args := AdminClearAllArgs{}
	success, err := c.client.Call(ctx, c.thriftService, "clearAll", &args, &resp)
	if err == nil && !success {
		switch {
		case resp.NotAuthorized != nil:
			err = resp.NotAuthorized
		default:
			err = fmt.Errorf("received no result or unknown exception for clearAll")
		}
	}

	return err
}

type tchanAdminServer struct {
	thrift.TChanServer

	handler TChanAdmin
}

// NewTChanAdminServer wraps a handler for TChanAdmin so it can be
// registered with a thrift.Server.
func NewTChanAdminServer(handler TChanAdmin) thrift.TChanServer {
	return &tchanAdminServer{
		NewTChanBaseServiceServer(handler),
		handler,
	}
}

func (s *tchanAdminServer) Service() string {
	return "Admin"
}

func (s *tchanAdminServer) Methods() []string {
	return []string{
		"clearAll",

		"HealthCheck",
	}
}

func (s *tchanAdminServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	switch methodName {
	case "clearAll":
		return s.handleClearAll(ctx, protocol)

	case "HealthCheck":
		return s.TChanServer.Handle(ctx, methodName, protocol)
	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanAdminServer) handleClearAll(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req AdminClearAllArgs
	var res AdminClearAllResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	err :=
		s.handler.ClearAll(ctx)

	if err != nil {
		switch v := err.(type) {
		case *NotAuthorized:
			if v == nil {
				return false, nil, fmt.Errorf("Handler for notAuthorized returned non-nil error type *NotAuthorized but nil value")
			}
			res.NotAuthorized = v
		default:
			return false, nil, err
		}
	} else {
	}

	return err == nil, &res, nil
}

type tchanKeyValueClient struct {
	TChanBaseService

	thriftService string
	client        thrift.TChanClient
}

func NewTChanKeyValueInheritedClient(thriftService string, client thrift.TChanClient) *tchanKeyValueClient {
	return &tchanKeyValueClient{
		NewTChanBaseServiceInheritedClient(thriftService, client),
		thriftService,
		client,
	}
}

// NewTChanKeyValueClient creates a client that can be used to make remote calls.
func NewTChanKeyValueClient(client thrift.TChanClient) TChanKeyValue {
	return NewTChanKeyValueInheritedClient("KeyValue", client)
}

func (c *tchanKeyValueClient) Get(ctx thrift.Context, key string) (string, error) {
	var resp KeyValueGetResult
	args := KeyValueGetArgs{
		Key: key,
	}
	success, err := c.client.Call(ctx, c.thriftService, "Get", &args, &resp)
	if err == nil && !success {
		switch {
		case resp.NotFound != nil:
			err = resp.NotFound
		case resp.InvalidKey != nil:
			err = resp.InvalidKey
		default:
			err = fmt.Errorf("received no result or unknown exception for Get")
		}
	}

	return resp.GetSuccess(), err
}

func (c *tchanKeyValueClient) Set(ctx thrift.Context, key string, value string) error {
	var resp KeyValueSetResult
	args := KeyValueSetArgs{
		Key:   key,
		Value: value,
	}
	success, err := c.client.Call(ctx, c.thriftService, "Set", &args, &resp)
	if err == nil && !success {
		switch {
		case resp.InvalidKey != nil:
			err = resp.InvalidKey
		default:
			err = fmt.Errorf("received no result or unknown exception for Set")
		}
	}

	return err
}

type tchanKeyValueServer struct {
	thrift.TChanServer

	handler TChanKeyValue
}

// NewTChanKeyValueServer wraps a handler for TChanKeyValue so it can be
// registered with a thrift.Server.
func NewTChanKeyValueServer(handler TChanKeyValue) thrift.TChanServer {
	return &tchanKeyValueServer{
		NewTChanBaseServiceServer(handler),
		handler,
	}
}

func (s *tchanKeyValueServer) Service() string {
	return "KeyValue"
}

func (s *tchanKeyValueServer) Methods() []string {
	return []string{
		"Get",
		"Set",

		"HealthCheck",
	}
}

func (s *tchanKeyValueServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	switch methodName {
	case "Get":
		return s.handleGet(ctx, protocol)
	case "Set":
		return s.handleSet(ctx, protocol)

	case "HealthCheck":
		return s.TChanServer.Handle(ctx, methodName, protocol)
	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanKeyValueServer) handleGet(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req KeyValueGetArgs
	var res KeyValueGetResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	r, err :=
		s.handler.Get(ctx, req.Key)

	if err != nil {
		switch v := err.(type) {
		case *KeyNotFound:
			if v == nil {
				return false, nil, fmt.Errorf("Handler for notFound returned non-nil error type *KeyNotFound but nil value")
			}
			res.NotFound = v
		case *InvalidKey:
			if v == nil {
				return false, nil, fmt.Errorf("Handler for invalidKey returned non-nil error type *InvalidKey but nil value")
			}
			res.InvalidKey = v
		default:
			return false, nil, err
		}
	} else {
		res.Success = &r
	}

	return err == nil, &res, nil
}

func (s *tchanKeyValueServer) handleSet(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req KeyValueSetArgs
	var res KeyValueSetResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	err :=
		s.handler.Set(ctx, req.Key, req.Value)

	if err != nil {
		switch v := err.(type) {
		case *InvalidKey:
			if v == nil {
				return false, nil, fmt.Errorf("Handler for invalidKey returned non-nil error type *InvalidKey but nil value")
			}
			res.InvalidKey = v
		default:
			return false, nil, err
		}
	} else {
	}

	return err == nil, &res, nil
}

type tchanBaseServiceClient struct {
	thriftService string
	client        thrift.TChanClient
}

func NewTChanBaseServiceInheritedClient(thriftService string, client thrift.TChanClient) *tchanBaseServiceClient {
	return &tchanBaseServiceClient{
		thriftService,
		client,
	}
}

// NewTChanBaseServiceClient creates a client that can be used to make remote calls.
func NewTChanBaseServiceClient(client thrift.TChanClient) TChanBaseService {
	return NewTChanBaseServiceInheritedClient("baseService", client)
}

func (c *tchanBaseServiceClient) HealthCheck(ctx thrift.Context) (string, error) {
	var resp BaseServiceHealthCheckResult
	args := BaseServiceHealthCheckArgs{}
	success, err := c.client.Call(ctx, c.thriftService, "HealthCheck", &args, &resp)
	if err == nil && !success {
		switch {
		default:
			err = fmt.Errorf("received no result or unknown exception for HealthCheck")
		}
	}

	return resp.GetSuccess(), err
}

type tchanBaseServiceServer struct {
	handler TChanBaseService
}

// NewTChanBaseServiceServer wraps a handler for TChanBaseService so it can be
// registered with a thrift.Server.
func NewTChanBaseServiceServer(handler TChanBaseService) thrift.TChanServer {
	return &tchanBaseServiceServer{
		handler,
	}
}

func (s *tchanBaseServiceServer) Service() string {
	return "baseService"
}

func (s *tchanBaseServiceServer) Methods() []string {
	return []string{
		"HealthCheck",
	}
}

func (s *tchanBaseServiceServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	switch methodName {
	case "HealthCheck":
		return s.handleHealthCheck(ctx, protocol)

	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanBaseServiceServer) handleHealthCheck(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req BaseServiceHealthCheckArgs
	var res BaseServiceHealthCheckResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	r, err :=
		s.handler.HealthCheck(ctx)

	if err != nil {
		return false, nil, err
	} else {
		res.Success = &r
	}

	return err == nil, &res, nil
}
