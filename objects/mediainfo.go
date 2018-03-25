package objects

type MediaInfo struct {
	Type             string `json:"type"`
	Title            string `json:"title"`
	Container        string `json:"container"`
	Year             int    `json:"year"`
	ReleaseGroup     string `json:"release_group"`
	Format           string `json:"format"`
	ScreenSize       string `json:"screen_size"`
	AudioChannels    string `json:"audio_channels"`
	AudioCodecs      string `json:"audio_codecs"`
	Language         string `json:"language"`
	SubtitleLanguage string `json:"subtitle_language"`
	Size             string `json:"size"`
}
