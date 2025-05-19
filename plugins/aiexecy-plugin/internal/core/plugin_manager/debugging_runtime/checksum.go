package debugging_runtime

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"aiexec-plugin/internal/utils/parser"
)

func (m *RemotePluginRuntime) calculateChecksum() string {
	configuration := m.Configuration()
	// calculate using sha256
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, parser.MarshalJsonBytes(configuration))
	hash := sha256.New()
	hash.Write(append(buffer.Bytes(), []byte(m.tenantId)...))
	return hex.EncodeToString(hash.Sum(nil))
}
