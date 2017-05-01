package urls

import (
	"fmt"
	"net/url"
	"strings"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/garyburd/redigo/redis"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type _Service struct {
	cache          redis.Conn
	redisNamespace string
	shortProtocol  string
	mongoSession   *mgo.Session
}

func newService(cache redis.Conn, mongoSession *mgo.Session, redisNamespace, shortProtocol string) *_Service {
	return &_Service{
		cache:          cache,
		redisNamespace: redisNamespace,
		shortProtocol:  shortProtocol,
		mongoSession:   mongoSession,
	}
}

func (service *_Service) Create(longURL, host string) (*_ShorterURL, error) {
	for i := 0; i < 10; i++ {
		shortURL := service.generateShortURL(host)

		shorterURL := &_ShorterURL{
			LongURL:  longURL,
			ShortURL: shortURL,
		}

		session := service.mongoSession.Copy()
		urls := session.DB("").C("urls")

		changeInfo, err := urls.Upsert(bson.M{"shortUrl": shortURL}, bson.M{"$setOnInsert": shorterURL})
		session.Close()
		if err != nil {
			return nil, err
		}
		if changeInfo.Matched > 0 {
			continue
		}

		return shorterURL, nil
	}

	return nil, fmt.Errorf("Failed to generate a unique url in 10 tries")
}

func (service *_Service) Delete(host, token string) error {
	shortURL := (&url.URL{Scheme: service.shortProtocol, Host: host, Path: token}).String()
	key := fmt.Sprintf("%v:%v", service.redisNamespace, shortURL)

	_, err := service.cache.Do("DEL", key)
	if err != nil && err.Error() != "redigo: nil returned" {
		return nil
	}

	session := service.mongoSession.Copy()
	urls := session.DB("").C("urls")
	err = urls.Remove(bson.M{"shortUrl": shortURL})
	session.Close()

	return err
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
