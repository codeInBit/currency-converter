package rates

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
)

var myCache CacheItf

type Rate struct {
	ConversionRate   float64 `json:"conversion_rate"`
	ConversionResult float64 `json:"conversion_result"`
}

type ForCache struct {
	Value float64 `json:"value"`
}

type CacheItf interface {
	Set(key string, data interface{}, expiration time.Duration) error
	Get(key string) ([]byte, error)
}

func fetchRate(url string, cacheKey string) (float64, float64) {
	var err error
	var result Rate
	var forCache ForCache

	start := time.Now()

	//fetch from cache if data exist
	cacheValue, err := myCache.Get(cacheKey)
	if err != nil {
		log.Fatal(err)
	}

	if cacheValue != nil { // cache exist
		err := json.Unmarshal(cacheValue, &forCache)
		if err != nil {
			log.Fatal(err)
		}

		elapsed := time.Since(start).Seconds()

		return forCache.Value, elapsed
	}

	err = godotenv.Load()

	//Fetch from external API, and load into cache
	response, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	conversionRate := result.ConversionRate

	// Store into local cache
	// myCache.Set(cacheKey, conversionRate, 1*time.Minute)

	elapsed := time.Since(start).Seconds()

	return conversionRate, elapsed
}

type AppCache struct {
	client *cache.Cache
}

func (r *AppCache) Set(key string, data interface{}, expiration time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	r.client.Set(key, b, expiration)
	return nil
}

func (r *AppCache) Get(key string) ([]byte, error) {
	res, exist := r.client.Get(key)
	if !exist {
		return nil, nil
	}

	resByte, ok := res.([]byte)
	if !ok {
		return nil, errors.New("Format is not arr of bytes")
	}

	return resByte, nil
}

func InitCache() {
	myCache = &AppCache{
		client: cache.New(5*time.Minute, 10*time.Minute),
	}
}
