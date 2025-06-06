package signer

import (
	"aiexec-plugin/internal/core/license/private_key"
	"aiexec-plugin/internal/utils/encryption"
	"aiexec-plugin/pkg/plugin_packager/signer/withkey"
)

/*
	AiexecPlugin is a file type that represents a plugin, it's designed to based on zip file format.
	When signing a plugin, we use RSA-4096 to create a signature for the plugin and write the signature
	into comment field of the zip file.
*/

// SignPlugin is a function that signs a plugin
// It takes a plugin as a stream of bytes and signs it with RSA-4096 with a bundled private key
func SignPlugin(plugin []byte) ([]byte, error) {
	// load private key
	privateKey, err := encryption.LoadPrivateKey(private_key.PRIVATE_KEY)
	if err != nil {
		return nil, err
	}

	return withkey.SignPluginWithPrivateKey(plugin, privateKey)
}
