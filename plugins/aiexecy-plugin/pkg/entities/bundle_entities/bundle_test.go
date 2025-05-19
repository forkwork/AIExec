package bundle_entities

import (
	"testing"

	"aiexec-plugin/internal/utils/parser"
)

func TestAvoidNilTags(t *testing.T) {
	yaml := `name: test
`
	bundle, err := parser.UnmarshalYamlBytes[Bundle]([]byte(yaml))
	if err != nil {
		t.Fatal(err)
	}

	if bundle.Tags == nil {
		t.Fatal("tags should not be nil")
	}
}
