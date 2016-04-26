package log

import (
	"os"
)

// RPC implements the RPCLogger.
type RPC struct {
	*Logger
	rpcID     string
	requestID string
}

// NewRPCLogger creates a RPC with given name as a RPCLogger.
func NewRPCLogger(name string) RPCLogger {
	l := NewWithWriter(name, os.Stdout)
	rpc := &RPC{Logger: l}
	l.recordFactory = NewRPCRecordFactory(rpc.rpcID, rpc.requestID)
	return rpc
}

// RPCID returns the RPCID of RPC logger.
func (l *RPC) RPCID() string {
	return l.rpcID
}

// RequestID returns the requestID of RPC logger.
func (l *RPC) RequestID() string {
	return l.requestID
}

// WithRPCID set rpcID on logger.
func (l *RPC) WithRPCID(rpcID string) RPCLogger {
	l.recordFactory = NewRPCRecordFactory(rpcID, l.requestID)
	return &RPC{
		Logger:    l.Logger,
		rpcID:     rpcID,
		requestID: l.requestID,
	}
}

// WithRequestID set requestID on logger.
func (l *RPC) WithRequestID(requestID string) RPCLogger {
	l.recordFactory = NewRPCRecordFactory(l.rpcID, requestID)
	return &RPC{
		Logger:    l.Logger,
		rpcID:     l.rpcID,
		requestID: requestID,
	}
}
