package main

import (
	"context"
	"os"

	"github.com/Shuv1Wolf/subterra-locate/facades/geo-facade/container"
)

func main() {
	proc := container.NewFacadeProcess()
	proc.SetConfigPath("./config/config.yml")
	proc.Run(context.Background(), os.Args)
}

