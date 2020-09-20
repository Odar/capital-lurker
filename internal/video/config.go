package video

// Config of client secret for uploading video on YouTube
type YouTubeClientSecretConfig struct {
	ClientID                string
	ProjectID               string
	AuthUri                 string
	TokenUri                string
	AuthProviderX509CertUrl string
	ClientSecret            string
}

// Config of video's resource for uploading on YouTube
type YouTubeVideoResourceConfig struct {
	FileName    string
	Title       string
	Description string
	Category    string
	Keywords    string
	Privacy     string
}
