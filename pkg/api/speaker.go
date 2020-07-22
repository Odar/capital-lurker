package api

type GetSpeakerOnMainRequest struct {
    Limit int64 `query:"limit"`
}

type SpeakerOnMain struct {
    ID       uint64 `json:"id"`
    Name     string `json:"name"`
    Position uint64 `json:"position"`
    Img      string `json:"img"`
}

type GetSpeakerOnMainResponse struct {
    Speakers []SpeakerOnMain `json:"speakers"`
}
