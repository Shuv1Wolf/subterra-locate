package main

import (
	"context"
	"os"

	"github.com/Shuv1Wolf/subterra-locate/services/geo-renderer/containers"
)

func main() {
	proc := containers.NewMapProcess()
	proc.SetConfigPath("../config/config.yml")
	proc.Run(context.Background(), os.Args)
}
