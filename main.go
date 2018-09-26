package main

import (
	"context"
	"flag"
	"log"
	"time"

	shopconfig "github.com/MYOB-Technology/shell-operator/pkg/config"
	"github.com/MYOB-Technology/shell-operator/pkg/executor"
	"github.com/MYOB-Technology/shell-operator/pkg/watcher"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/golang/glog"
)

var shellConfigPath string

func init() {
	flag.StringVar(&shellConfigPath, "shell-config", "/app/shell-config.yaml", "Provide the path to the shell config yaml file.")
}

func main() {
	flag.Parse()

	glog.V(1).Infof("Loading config from %s...", shellConfigPath)
	shellConfig, err := shopconfig.ParseFromFile(shellConfigPath)

	if err != nil {
		glog.Fatal(err)
	}

	cfg, err := config.GetConfig()

	if err != nil {
		glog.Fatal(err)
	}

	mgr, err := manager.New(cfg, manager.Options{})

	if err != nil {
		log.Fatal(err)
	}

	err = watcher.SetupWatches(mgr, shellConfig)

	if err != nil {
		glog.Fatal(err)
	}

	// Run any boot commands
	for _, b := range shellConfig.Boot {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(b.Timeout)*time.Second)
		defer cancel()
		cmd := executor.SetupShellCommand(ctx, b.Command, b.Environment)
		glog.V(1).Infof("Executing boot command:")
		out, err := cmd.CombinedOutput()

		if err != nil {
			glog.Fatal(err)
		}

		glog.V(1).Infof(string(out))
	}

	// Start the Cmd
	log.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
