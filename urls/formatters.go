package urls

import "encoding/json"

func formatCreateResponse(shorterURL *_ShorterURL) ([]byte, error) {
	return json.MarshalIndent(struct {
		ShortURL string `json:"shortUrl"`
		LongURL  string `json:"longUrl"`
	}{
		ShortURL: shorterURL.ShortURL,
		LongURL:  shorterURL.LongURL,
	}, "", "  ")
}
