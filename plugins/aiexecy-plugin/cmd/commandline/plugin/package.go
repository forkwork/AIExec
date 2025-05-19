package plugin

import (
	"os"

	"aiexec-plugin/internal/utils/log"
	"aiexec-plugin/pkg/plugin_packager/decoder"
	"aiexec-plugin/pkg/plugin_packager/packager"
)

var (
	MaxPluginPackageSize = int64(52428800) // 50MB
)

func PackagePlugin(inputPath string, outputPath string) {
	decoder, err := decoder.NewFSPluginDecoder(inputPath)
	if err != nil {
		log.Error("failed to create plugin decoder , plugin path: %s, error: %v", inputPath, err)
		os.Exit(1)
		return
	}

	packager := packager.NewPackager(decoder)
	zipFile, err := packager.Pack(MaxPluginPackageSize)

	if err != nil {
		log.Error("failed to package plugin %v", err)
		os.Exit(1)
		return
	}

	err = os.WriteFile(outputPath, zipFile, 0644)
	if err != nil {
		log.Error("failed to write package file %v", err)
		os.Exit(1)
		return
	}

	log.Info("plugin packaged successfully, output path: %s", outputPath)
}
