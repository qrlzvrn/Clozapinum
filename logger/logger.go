package logger

import (
	"github.com/qrlzvrn/Clozapinum/erro"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// FatalMe - логиурет ошибки
func FatalMe(e erro.Err) {
	err, action := e.Erro()
	handlerLogger := log.WithFields(logrus.Fields{
		"action": action,
		"err":    err,
	})
	handlerLogger.Fatalf("Error: %+v", err)
}

func BotSendFatal(err error, obj string) {
	action := "sendMessage"
	handlerLogger := log.WithFields(logrus.Fields{
		"action": action,
		"obj":    obj,
		"err":    err,
	})
	handlerLogger.Fatalf("Error: %+v", err)
}

func Infof(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}
