package commands

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/stevvooe/continuity"
)

func readManifest(path string) (*continuity.Manifest, error) {
	p, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading manifest: %v", err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	mediaType := continuity.MediaTypeManifestV0Protobuf
	if ext == ".json" {
		mediaType = continuity.MediaTypeManifestV0JSON
	}

	m, err := continuity.Unmarshal(p, mediaType)
	if err != nil {
		return m, fmt.Errorf("error unmarshalling manifest: %v", err)
	}
	return m, nil
}
