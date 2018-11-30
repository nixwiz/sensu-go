package main

import (
	_ "net/http/pprof"

	"github.com/sensu/sensu-go/backend"
	"github.com/sensu/sensu-go/backend/cmd"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithFields(logrus.Fields{
	"component": "backend",
})

func main() {
	if err := cmd.Execute(backend.Initialize); err != nil {
		logger.WithError(err).Fatal("error executing sensu-backend")
	}
}
