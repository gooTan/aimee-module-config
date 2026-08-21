package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RakuenSoftware/aimee-module-config/server-go/bus"
	handler "github.com/RakuenSoftware/aimee-module-config/server-go/modules/config"
)

func main() {
	moduleHandler, err := handler.NewDefaultHandler()
	if err != nil {
		fmt.Fprintf(os.Stderr, "aimee-module-config: %v\n", err)
		os.Exit(1)
	}
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s DAEMON_MODULE_BUS_SOCKET\n", os.Args[0])
		os.Exit(2)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config := bus.ModuleProcessConfig{
		SocketPath: os.Args[1], ModuleName: "config",
		PrincipalClass: 1, PrincipalRef: 2,
		Stages: []bus.ModuleStage{
			{EventKind: 4609, StageID: 1},
		},
		Handler: moduleHandler,
	}
	if err := bus.RunModuleProcess(ctx, config); err != nil {
		fmt.Fprintf(os.Stderr, "aimee-module-config: %v\n", err)
		os.Exit(1)
	}
}
