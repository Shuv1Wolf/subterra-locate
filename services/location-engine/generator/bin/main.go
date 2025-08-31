package main

import (
	"context"
	"os"

	"github.com/Shuv1Wolf/subterra-locate/services/location-engine/generator/containers"
)

func main() {
	proc := containers.NewGenProcess()
	proc.SetConfigPath("./config/config.yml")
	proc.Run(context.Background(), os.Args)
}
