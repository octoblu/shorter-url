package cachedurls

import (
	"fmt"
	"net/url"

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
	_, err := service.cache.Do("SET", key, longURL)
	return err
}

func (service *_Service) getLongURLFromCache(shortURL string) (string, error) {
	key := fmt.Sprintf("%v:%v", service.redisNamespace, shortURL)
	longURL, err := redis.String(service.cache.Do("GET", key))
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
	session := service.mongoSession.Copy()
	urls := session.DB("").C("urls")
	err := urls.Find(bson.M{"shortUrl": shortURL}).One(shorterURL)
	session.Close()
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
