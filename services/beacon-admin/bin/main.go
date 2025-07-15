package main

import (
	"context"
	"os"

	"github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/containers"
)

func main() {
	proc := containers.NewBeaconsProcess()
	proc.SetConfigPath("../config/config.yml")
	proc.Run(context.Background(), os.Args)
}
