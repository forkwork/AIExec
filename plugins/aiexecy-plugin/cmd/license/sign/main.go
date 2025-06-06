package main

import (
	"flag"
	"os"

	"aiexec-plugin/internal/utils/log"
	"aiexec-plugin/pkg/plugin_packager/signer"
)

func main() {
	var (
		in_path  string
		out_path string
		help     bool
	)

	flag.StringVar(&in_path, "in", "", "input plugin file path")
	flag.StringVar(&out_path, "out", "", "output plugin file path")
	flag.BoolVar(&help, "help", false, "show help")
	flag.Parse()

	if help || in_path == "" || out_path == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// read plugin
	plugin, err := os.ReadFile(in_path)
	if err != nil {
		log.Panic("failed to read plugin file %v", err)
	}

	// sign plugin
	pluginFile, err := signer.SignPlugin(plugin)
	if err != nil {
		log.Panic("failed to sign plugin %v", err)
	}

	// write signature
	err = os.WriteFile(out_path, pluginFile, 0644)
	if err != nil {
		log.Panic("failed to write plugin file %v", err)
	}

	log.Info("plugin signed successfully, output path: %s", out_path)
}
