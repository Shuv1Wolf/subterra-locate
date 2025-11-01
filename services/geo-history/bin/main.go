package main

import (
	"context"
	"os"

	"github.com/Shuv1Wolf/subterra-locate/services/geo-history/containers"
)

func main() {
	proc := containers.NewGeoHistoryProcess()
	proc.SetConfigPath("../config/config.yml")
	proc.Run(context.Background(), os.Args)
}
