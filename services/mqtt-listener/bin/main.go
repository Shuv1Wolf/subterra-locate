package main

import (
	"context"
	"os"

	"github.com/Shuv1Wolf/subterra-locate/services/mqtt-listener/containers"
)

func main() {
	proc := containers.NewMqttListenerProcess()
	proc.SetConfigPath("./config/config.yml")
	proc.Run(context.Background(), os.Args)
}
