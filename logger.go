package log

type MultiHandlerLogger interface {
	AddHandler(h Handler)
	RemoveHandler(h Handler)
	Handlers() []Handler
}

type Logger interface {
	MultiHandlerLogger
	LevelLogger
	DebugLogger
	PrintLogger
	InfoLogger
	WarnLogger
	FatalLogger
	SetRPCID(rpcID string)
	SetRequestID(requestID string)
}

type LevelLogger interface {
	Level() LevelType
	SetLevel(lv LevelType)
}

type DebugLogger interface {
	// Debug APIs
	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
}

type PrintLogger interface {
	// Print APIs
	Print(a ...interface{})
	Println(a ...interface{})
	Printf(f string, a ...interface{})
}

type InfoLogger interface {
	// Info APIs
	Info(a ...interface{})
	Infof(f string, a ...interface{})
}

type WarnLogger interface {
	// Warn APIs
	Warn(a ...interface{})
	Warnf(f string, a ...interface{})
}

type FatalLogger interface {
	// Fatal APIs
	Fatal(a ...interface{})
	Fatalf(f string, a ...interface{})
}
