package bundle

import (
	"aiexec-plugin/internal/utils/log"
	"aiexec-plugin/pkg/entities/manifest_entities"
)

func BumpVersion(bundlePath string, targetVersion string) {
	packager, err := loadBundlePackager(bundlePath)
	if err != nil {
		log.Error("Failed to load bundle packager: %v", err)
		return
	}

	targetVersionObject, err := manifest_entities.NewVersion(targetVersion)
	if err != nil {
		log.Error("Failed to parse target version: %v", err)
		return
	}

	packager.BumpVersion(targetVersionObject)
	if err := packager.Save(); err != nil {
		log.Error("Failed to save bundle packager: %v", err)
		return
	}
}
