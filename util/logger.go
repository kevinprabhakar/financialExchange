package util

import (
"log"
"io"
"runtime/debug"
)

type Logger struct{
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

func NewLogger(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) (*Logger) {
	var newLogger Logger

	newLogger.Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	newLogger.Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	newLogger.Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	newLogger.Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return &newLogger
}

func (self *Logger)Debug(inStr string){
	self.Info.Println(inStr)
}

func (self *Logger)Warn(inStr string){
	self.Warning.Println(inStr)
}

func (self *Logger)ErrorMsg(inStr string){
	self.Error.Println(inStr)
	debug.PrintStack()
}
