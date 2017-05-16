package urls

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func parseCreateBody(body io.ReadCloser) (_CreateBody, error) {
	createBody := _CreateBody{}

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return createBody, err
	}

	err = json.Unmarshal(bodyBytes, &createBody)
	if err != nil {
		return createBody, err
	}

	if createBody.LongURL == "" {
		return createBody, fmt.Errorf("Missing required property 'longUrl'")
	}

	return createBody, nil
}

type _CreateBody struct {
	LongURL  string `json:"longUrl"`
	ShortURL string `json:"shortUrl"`
}
