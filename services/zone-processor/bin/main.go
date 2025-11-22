package main

import (
	"context"
	"os"

	"github.com/Shuv1Wolf/subterra-locate/services/zone-processor/containers"
)

func main() {
	proc := containers.NewZoneProcess()
	proc.SetConfigPath("../config/config.yml")
	proc.Run(context.Background(), os.Args)
}
