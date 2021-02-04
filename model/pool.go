package model

// Pool :
type Pool struct {
	UUID          string `json:"uuid"`
	Size          string `json:"size"`
	Free          string `json:"free"`
	Capacity      string `json:"capacity"`
	Health        string `json:"health"`
	Name          string `json:"name"`
	AvailableSize string `json:"availablesize"`
	Action        string `json:"action"`
	Used          string `json:"used"`
}
