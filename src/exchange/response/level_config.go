package response

type LevelConfigResponse struct {
	Data LevelConfigResponseData
}

type LevelConfigResponseData struct {
	Version string      `json:"version"`
	Config  interface{} `json:"config"`
}
