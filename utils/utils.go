package utils

import (
	"github.com/RichardKnop/machinery/v2"
	"honnef.co/go/tools/config"
)

func GetMachineryServer() *machinery.Server {
	Logger.Info("initing task server")

	taskserver, err := machinery.NewServer(&config.Config{
		Broker:        "redis://localhost:6379",
		ResultBackend: "redis://localhost:6379",
	})
	if err != nil {
		Logger.Fatal(err.Error())
	}

	taskserver.RegisterTasks(map[string]interface{}{
		"send_email": tasks.SendMail,
	})

	return taskserver
}
