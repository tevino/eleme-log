package log

// ELogRecord represents a ELogRecord with rpcId and requestID.
type ELogRecord struct {
	*BaseRecord
	rpcID     string
	requestID string
}

// NewELogRecord create a ELogRecord with rpcID and requestID.
func NewELogRecord(name string, calldepth int, lv LevelType, msg string, rpcID string, requestID string) *ELogRecord {
	return &ELogRecord{
		BaseRecord: NewBaseRecord(name, calldepth, lv, msg),
		rpcID:      rpcID,
		requestID:  requestID,
	}
}

// NewELogRecordFactory return a record factory for RPCRecord.
func NewELogRecordFactory(rpcID string, requestID string) RecordFactory {
	return func(name string, calldepth int, lv LevelType, msg string) Record {
		return NewELogRecord(name, calldepth+1, lv, msg, rpcID, requestID)
	}
}
