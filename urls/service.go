package urls

import (
	"net/url"
	"strings"

	randomdata "github.com/Pallinder/go-randomdata"
	mgo "gopkg.in/mgo.v2"
)

type _Service struct {
	shortProtocol string
	urls          *mgo.Collection
}

func newService(mongoDB *mgo.Database, shortProtocol string) *_Service {
	return &_Service{urls: mongoDB.C("urls"), shortProtocol: shortProtocol}
}

func (service *_Service) Create(longURL, scheme, host string) (*_ShorterURL, error) {
	shortURL, err := service.generateShortURL(service.shortProtocol, host)
	if err != nil {
		return nil, err
	}

	shorterURL := &_ShorterURL{
		LongURL:  longURL,
		ShortURL: shortURL,
	}

	err = service.urls.Insert(shorterURL)
	return shorterURL, err
}

func (service *_Service) generateShortURL(scheme, host string) (string, error) {
	token := strings.ToLower(randomdata.Letters(4))

	shortURL := &url.URL{Scheme: scheme, Host: host, Path: token}

	return shortURL.String(), nil
}

type _ShorterURL struct {
	ShortURL string `bson:"shortUrl"`
	LongURL  string `bson:"longUrl"`
}
