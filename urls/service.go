package urls

import (
	"fmt"
	"net/url"
	"strings"

	randomdata "github.com/Pallinder/go-randomdata"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type _Service struct {
	shortProtocol string
	urls          *mgo.Collection
}

func newService(mongoDB *mgo.Database, shortProtocol string) *_Service {
	return &_Service{urls: mongoDB.C("urls"), shortProtocol: shortProtocol}
}

func (service *_Service) Create(longURL, scheme, host string) (*_ShorterURL, error) {
	shortURL, err := service.generateUniqueShortURL(host)
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

func (service *_Service) generateUniqueShortURL(host string) (string, error) {
	for i := 0; i < 10; i++ {
		shortURL := service.generateShortURL(host)
		n, err := service.urls.Find(bson.M{"shortUrl": shortURL}).Count()
		if err != nil {
			return "", err
		}

		if n == 0 {
			return shortURL, nil
		}
	}
	return "", fmt.Errorf("Failed to generate a unique url in 10 tries")
}

func (service *_Service) generateShortURL(host string) string {
	token := strings.ToLower(randomdata.Letters(4))

	shortURL := &url.URL{Scheme: service.shortProtocol, Host: host, Path: token}

	return shortURL.String()
}

type _ShorterURL struct {
	ShortURL string `bson:"shortUrl"`
	LongURL  string `bson:"longUrl"`
}
