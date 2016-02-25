package log

type Namer interface {
	Name() string
}

type Leveler interface {
	Level() LevelType
	SetLevel(lv LevelType)
}

type NamedLeveler interface {
	Namer
	Leveler
}

type MultiHandler interface {
	AddHandler(h Handler)
	RemoveHandler(h Handler)
	Handlers() []Handler
}

type Logger interface {
	// Basic
	NamedLeveler

	// multiple handlers
	MultiHandler

	// level APIs
	Debugger
	Printer
	Infoer
	Warner
	Fataler
}

type RPCLogger interface {
	Logger
	// RPC APIs
	RPCID() string
	RequestID() string
	SetRPCID(rpcID string)
	SetRequestID(requestID string)
}

type Debugger interface {
	// Debug APIs
	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
}

type Printer interface {
	// Print APIs
	Print(a ...interface{})
	Println(a ...interface{})
	Printf(f string, a ...interface{})
}

type Infoer interface {
	// Info APIs
	Info(a ...interface{})
	Infof(f string, a ...interface{})
}

type Warner interface {
	// Warn APIs
	Warn(a ...interface{})
	Warnf(f string, a ...interface{})
}

type Fataler interface {
	// Fatal APIs
	Fatal(a ...interface{})
	Fatalf(f string, a ...interface{})
}
