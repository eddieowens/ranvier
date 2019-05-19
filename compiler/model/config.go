package model

type ConfigManifest struct {
	Name    string   `json:"name"`
	Extends []string `json:"extends"`
	Config  []byte   `json:"config"`
}
