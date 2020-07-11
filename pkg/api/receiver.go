package api

type ReverseRequest struct {
	Word string `query:"word"`
}

type ReverseResponse struct {
	Word string `json:"word"`
}
