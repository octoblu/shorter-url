package serverinfo

import "encoding/json"

func formatHealthcheck() ([]byte, error) {
	return json.MarshalIndent(struct {
		Online bool `json:"online"`
	}{
		Online: true,
	}, "", "  ")
}

func formatVersion(version string) ([]byte, error) {
	return json.MarshalIndent(struct {
		Version string `json:"version"`
	}{
		Version: version,
	}, "", "  ")
}
