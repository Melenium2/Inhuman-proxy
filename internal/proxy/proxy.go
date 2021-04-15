package proxy

// easyjson:json
type RequestProxy struct {
	Code    string   `json:"code"`
	Address []string `json:"address"`
}
