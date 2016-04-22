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

// NewRPCLogger creates a Logger with given name as a RPCLogger.
func NewRPCLogger(name string) RPCLogger {
	return &RPC{
		Logger: NewWithWriter(name, os.Stdout),
	}
}

// WithRPCID set rpcID on logger.
func (l *RPC) WithRPCID(rpcID string) *Record {
	r := NewRecord(l.Logger, 2)
	r.rpcID = rpcID
	return r
}

// WithRequestID set requestID on logger.
func (l *RPC) WithRequestID(requestID string) *Record {
	r := NewRecord(l.Logger, 2)
	r.requestID = requestID
	return r
}
