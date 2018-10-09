package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	shopconfig "github.com/MYOB-Technology/shell-operator/pkg/config"
	"github.com/MYOB-Technology/shell-operator/pkg/metrics"
	"github.com/MYOB-Technology/shell-operator/pkg/shell"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = watcher.SetupWatches(ctx, mgr, shellConfig)

	if err != nil {
		glog.Fatal(err)
	}

	go func() { glog.Fatal(metrics.RunServer("8080")) }()

	// Run any boot commands
	for _, b := range shellConfig.Boot {
		ctx, bootCancel := context.WithTimeout(ctx, time.Duration(b.Timeout)*time.Second)
		cmd := shell.New(ctx, b.Command)
		shell.AddEnvironment(cmd, b.Environment)
		glog.V(1).Infof("Executing boot command:")
		err := shell.RunWithProgress(cmd, os.Stdout, os.Stderr)

		if err != nil {
			glog.Fatal(err)
		}

		bootCancel()
	}

	// Start the Cmd
	stopCh := signals.SetupSignalHandler()
	go func() {
		log.Fatal(mgr.Start(stopCh))
	}()

	<-stopCh
	glog.V(4).Infof("Recieved stop.")
}
