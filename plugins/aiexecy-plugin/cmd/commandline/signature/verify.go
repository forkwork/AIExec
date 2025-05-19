package signature

import (
	"crypto/rsa"
	"os"

	"aiexec-plugin/internal/utils/encryption"
	"aiexec-plugin/internal/utils/log"
	"aiexec-plugin/pkg/plugin_packager/decoder"
)

func Verify(aiexecpkgPath string, publicKeyPath string) error {
	// read the plugin
	plugin, err := os.ReadFile(aiexecpkgPath)
	if err != nil {
		log.Error("Failed to read plugin file: %v", err)
		return err
	}

	decoderInstance, err := decoder.NewZipPluginDecoder(plugin)
	if err != nil {
		log.Error("Failed to create plugin decoder, plugin path: %s, error: %v", aiexecpkgPath, err)
		return err
	}

	if publicKeyPath == "" {
		// verify the plugin with the official (bundled) public key
		err = decoder.VerifyPlugin(decoderInstance)
		if err != nil {
			log.Error("Failed to verify plugin with official public key: %v", err)
			return err
		}
	} else {
		// read the public key
		publicKeyBytes, err := os.ReadFile(publicKeyPath)
		if err != nil {
			log.Error("Failed to read public key file: %v", err)
			return err
		}

		publicKey, err := encryption.LoadPublicKey(publicKeyBytes)
		if err != nil {
			log.Error("Failed to load public key: %v", err)
			return err
		}

		// verify the plugin
		err = decoder.VerifyPluginWithPublicKeys(decoderInstance, []*rsa.PublicKey{publicKey})
		if err != nil {
			log.Error("Failed to verify plugin with provided public key: %v", err)
			return err
		}
	}

	log.Info("Plugin verified successfully")
	return nil
}
