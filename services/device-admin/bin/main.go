package main

import (
	"context"
	"os"

	"github.com/Shuv1Wolf/subterra-locate/services/device-admin/containers"
)

func main() {
	proc := containers.NewDeviceProcess()
	proc.SetConfigPath("../config/config.yml")
	proc.Run(context.Background(), os.Args)
}
