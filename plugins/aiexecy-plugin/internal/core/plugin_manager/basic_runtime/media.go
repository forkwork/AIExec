package basic_runtime

import (
	"aiexec-plugin/internal/core/plugin_manager/media_transport"
	"aiexec-plugin/pkg/entities/plugin_entities"
)

type MediaTransport struct {
	mediaManager *media_transport.MediaBucket

	assetsIds []string
}

// RemapAssets will take the assets and remap them to a media id
func (r *MediaTransport) RemapAssets(
	declaration *plugin_entities.PluginDeclaration,
	assets map[string][]byte,
) error {
	assetsIds, err := r.mediaManager.RemapAssets(declaration, assets)
	if err != nil {
		return err
	}

	r.assetsIds = assetsIds
	return nil
}

func NewMediaTransport(mediaManager *media_transport.MediaBucket) MediaTransport {
	return MediaTransport{mediaManager: mediaManager}
}
