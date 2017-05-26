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
	redisPool      *redis.Pool
	redisNamespace string
	shortProtocol  string
	mongoSession   *mgo.Session
}

func newService(redisPool *redis.Pool, mongoSession *mgo.Session, redisNamespace, shortProtocol string) *_Service {
	return &_Service{
		redisPool:      redisPool,
		redisNamespace: redisNamespace,
		shortProtocol:  shortProtocol,
		mongoSession:   mongoSession,
	}
}

func (service *_Service) Create(longURL, shortURLOverride, host string) (*_ShorterURL, error) {
	for i := 0; i < 10; i++ {
		shortURL := service.generateShortURL(host)
		if shortURLOverride != "" {
			shortURL = shortURLOverride
		}

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

	cache := service.redisPool.Get()
	defer cache.Close()
	_, err := cache.Do("DEL", key)
	if err != nil && err.Error() != "redigo: nil returned" {
		return nil
	}

	session := service.mongoSession.Copy()
	urls := session.DB("").C("urls")
	err = urls.Remove(bson.M{"shortUrl": shortURL})
	session.Close()

	return err
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
	cache := service.redisPool.Get()
	defer cache.Close()
	_, err := cache.Do("SET", key, longURL)
	return err
}

func (service *_Service) generateShortURL(host string) string {
	token := strings.ToLower(randomdata.Letters(4))

	shortURL := &url.URL{Scheme: service.shortProtocol, Host: host, Path: token}

	return shortURL.String()
}

func (service *_Service) getLongURLFromCache(shortURL string) (string, error) {
	key := fmt.Sprintf("%v:%v", service.redisNamespace, shortURL)
	cache := service.redisPool.Get()
	defer cache.Close()
	longURL, err := redis.String(cache.Do("GET", key))
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
	ShortURL string `bson:"shortUrl"`
	LongURL  string `bson:"longUrl"`
}
