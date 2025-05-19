package plugin

import (
	"os"

	"aiexec-plugin/internal/utils/log"
	"aiexec-plugin/pkg/plugin_packager/decoder"
)

func CalculateChecksum(pluginPath string) {
	var pluginDecoder decoder.PluginDecoder
	if stat, err := os.Stat(pluginPath); err == nil {
		if stat.IsDir() {
			pluginDecoder, err = decoder.NewFSPluginDecoder(pluginPath)
			if err != nil {
				log.Error("failed to create plugin decoder, plugin path: %s, error: %v", pluginPath, err)
				return
			}
		} else {
			bytes, err := os.ReadFile(pluginPath)
			if err != nil {
				log.Error("failed to read plugin file, plugin path: %s, error: %v", pluginPath, err)
				return
			}

			pluginDecoder, err = decoder.NewZipPluginDecoder(bytes)
			if err != nil {
				log.Error("failed to create plugin decoder, plugin path: %s, error: %v", pluginPath, err)
				return
			}
		}
	} else {
		log.Error("failed to get plugin file info, plugin path: %s, error: %v", pluginPath, err)
		return
	}

	checksum, err := pluginDecoder.Checksum()
	if err != nil {
		log.Error("failed to calculate checksum, plugin path: %s, error: %v", pluginPath, err)
		return
	}

	log.Info("plugin checksum: %s", checksum)
}
