package cachedurls

import (
	"fmt"
	"net/url"

	"github.com/garyburd/redigo/redis"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type _Service struct {
	redisConn      redis.Conn
	redisNamespace string
	shortProtocol  string
	urls           *mgo.Collection
}

func newService(mongoDB *mgo.Database, redisConn redis.Conn, redisNamespace, shortProtocol string) *_Service {
	return &_Service{
		redisConn:      redisConn,
		redisNamespace: redisNamespace,
		shortProtocol:  shortProtocol,
		urls:           mongoDB.C("urls"),
	}
}

func (service *_Service) GetLongURL(host, token string) (string, error) {
	shortURL := &url.URL{
		Scheme: service.shortProtocol,
		Host:   host,
		Path:   token,
	}

	longURL, err := service.getLongURLFromCache(shortURL.String())
	if err != nil {
		return "", err
	}
	if longURL != "" {
		return longURL, nil
	}

	return service.getLongURLFromMongo(shortURL.String())
}

func (service *_Service) cacheLongURL(shortURL, longURL string) error {
	key := fmt.Sprintf("%v:%v", service.redisNamespace, shortURL)
	_, err := service.redisConn.Do("SET", key, longURL)
	return err
}

func (service *_Service) getLongURLFromCache(shortURL string) (string, error) {
	key := fmt.Sprintf("%v:%v", service.redisNamespace, shortURL)
	longURL, err := redis.String(service.redisConn.Do("GET", key))
	if err != nil && err.Error() == "redigo: nil returned" {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func (service *_Service) getLongURLFromMongo(shortURL string) (string, error) {
	shorterURL := &_ShorterURL{}
	err := service.urls.Find(bson.M{"shortUrl": shortURL}).One(shorterURL)
	if err != nil {
		return "", err
	}

	err = service.cacheLongURL(shortURL, shorterURL.LongURL)
	if err != nil {
		return "", err
	}

	return shorterURL.LongURL, nil
}

type _ShorterURL struct {
	LongURL string `bson:"longUrl"`
}
