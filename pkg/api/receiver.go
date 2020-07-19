package api

type ReverseRequest struct {
	Word string `query:"word"`
}

type ReverseResponse struct {
	Word string `json:"word"`
	Uses uint64 `json:"uses"`
}
