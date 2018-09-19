package main

import (
	"flag"
	"log"
	"os"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var shellConfigPath string

func init() {
	flag.StringVar(&shellConfigPath, "shell-config", "/app/shell-config.yaml", "Provide the path to the shell config yaml file.")
}

func main() {
	flag.Parse()

	confIn, err := os.Open(shellConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	shellConfig, errs := ParseAndValidateConfig(confIn)

	if len(errs) > 0 {
		log.Fatal(errs)
	}

	cfg, err := config.GetConfig()

	if err != nil {
		log.Fatal(err)
	}

	mgr, err := manager.New(cfg, manager.Options{})

	if err != nil {
		log.Fatal(err)
	}

	err = SetupWatches(mgr, shellConfig)

	if err != nil {
		log.Fatal(err)
	}

	// Start the Cmd
	log.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
