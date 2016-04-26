package log

// RPCRecord represents a RPCRecord with rpcId and requestID.
type RPCRecord struct {
	*BaseRecord
	rpcID     string
	requestID string
}

// NewRPCRecord create a RPCRecord with rpcID and requestID.
func NewRPCRecord(name string, calldepth int, lv LevelType, msg string, rpcID string, requestID string) *RPCRecord {
	return &RPCRecord{
		BaseRecord: NewBaseRecord(name, calldepth, lv, msg),
		rpcID:      rpcID,
		requestID:  requestID,
	}
}

// NewRPCRecordFactory return a record factory for RPCRecord.
func NewRPCRecordFactory(rpcID string, requestID string) RecordFactory {
	return func(name string, calldepth int, lv LevelType, msg string) Record {
		return NewRPCRecord(name, calldepth+1, lv, msg, rpcID, requestID)
	}
}
