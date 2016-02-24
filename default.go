package log

const (
	defaultTpl = "{{level}} {{date}} {{time}} {{name}} {{}}"
	syslogTpl  = "[{{app_id}} {{rpc_id}} {{request_id}}] ## {{}}"
)

var defaultLevel = INFO
var defaultLogger = New("")

var (
	Level    = defaultLogger.Level
	SetLevel = defaultLogger.SetLevel
	Print    = defaultLogger.Print
	Printf   = defaultLogger.Printf
	Println  = defaultLogger.Println
	Fatal    = defaultLogger.Fatal
	Fatalf   = defaultLogger.Fatalf
)
