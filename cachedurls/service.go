package cachedurls

import (
	"net/url"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type _Service struct {
	shortProtocol string
	urls          *mgo.Collection
}

func newService(mongoDB *mgo.Database, shortProtocol string) *_Service {
	return &_Service{
		urls:          mongoDB.C("urls"),
		shortProtocol: shortProtocol,
	}
}

func (service *_Service) GetLongURL(host, token string) (string, error) {
	shortURL := &url.URL{
		Scheme: service.shortProtocol,
		Host:   host,
		Path:   token,
	}

	shorterURL := &_ShorterURL{}
	err := service.urls.Find(bson.M{"shortUrl": shortURL.String()}).One(shorterURL)
	if err != nil {
		return "", err
	}

	return shorterURL.LongURL, nil
}

type _ShorterURL struct {
	LongURL string `bson:"longUrl"`
}
