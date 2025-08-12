package main

import (
	"context"
	"os"

	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/containers"
)

func main() {
	proc := containers.NewLocationEngineProcess()
	proc.SetConfigPath("./config/config.yml")
	proc.Run(context.Background(), os.Args)
}
