package rds

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis"
	"github.com/kevin09002/gin-kit/log"
)

const (
	// GzipMinSize gzip min size
	GzipMinSize = 1024
	// CacheFormatRaw raw
	CacheFormatRaw = 0
	// CacheFormatRawGzip raw gzip
	CacheFormatRawGzip = 1
	// CacheFormatJSON json
	CacheFormatJSON = 10
	// CacheFormatJSONGzip json gzip
	CacheFormatJSONGzip = 11
)

// modelCacheItem model cache item
type modelCacheItem struct {
	Data []byte
	Flag uint32
}

// Config redis config
type Config struct {
	Address     string `yaml:"address" json:"address"`
	Password    string `yaml:"password" json:"password"`
	Timeout     int    `yaml:"timeout" json:"timeout"`
	MaxIdle     int    `yaml:"max_idle" json:"max_dile"`
	IdleTimeout int    `yaml:"idle_timeout" json:"idle_timeout"`
	RetryTimes  int    `yaml:"retry_times" json:"retry_times"`
}

type Client struct {
	RedisClient *redis.Client
	cfg         *Config
}

// NewClient new redis client
func NewClient(c *Config) *Client {
	return &Client{
		cfg: c,
		RedisClient: redis.NewClient(&redis.Options{
			Addr:         c.Address,
			Password:     c.Password,
			PoolSize:     c.MaxIdle,
			IdleTimeout:  time.Duration(c.IdleTimeout) * time.Second,
			DialTimeout:  time.Duration(c.Timeout) * time.Second,
			ReadTimeout:  time.Duration(c.Timeout) * time.Second,
			WriteTimeout: time.Duration(c.Timeout) * time.Second,
		}),
	}
}

// SetModelToCache save model to cache
func (c *Client) SetModelToCache(key string, model interface{}, ttl time.Duration) error {
	log.Tracef("[cache: set_model_to_cache]: key=%s", key)
	var (
		bs        []byte
		data      []byte
		err       error
		cacheFlag = CacheFormatJSON
		rc        = c.RedisClient
		gziped    bool
	)
	if gziped, data, err = toGzipJSON(model); err != nil {
		return err
	}
	if gziped {
		cacheFlag = CacheFormatJSONGzip
	}
	if bs, err = json.Marshal(modelCacheItem{
		Data: data,
		Flag: uint32(cacheFlag),
	}); err != nil {
		return err
	}
	if _, err = rc.Set(key, string(bs), ttl).Result(); err != nil {
		return err
	}
	return nil
}

// GetCacheToModel get cache to model
func (c *Client) GetCacheToModel(key string, model *interface{}) (bool, error) {
	log.Tracef("[cache: get_cache_to_model], key=%s", key)
	rc := c.RedisClient
	it, err := rc.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return true, nil
		}
		return false, err
	}
	cacheItem := modelCacheItem{}
	if err = json.Unmarshal([]byte(it), &cacheItem); err != nil {
		log.Errorf("[cache:%s] Unmarshal value error, %s", key, err)
		return false, err
	}
	switch cacheItem.Flag {
	case CacheFormatJSON:
		err = json.Unmarshal(cacheItem.Data, model)
	case CacheFormatJSONGzip:
		err = fromGzipJSON(cacheItem.Data, model)
	default:
		err = fmt.Errorf("invalid cache formate %d", cacheItem.Flag)
	}
	if err != nil {
		log.Errorf("[cache:%s] %s", key, err)
		return false, err
	}
	return false, nil
}

func toGzipJSON(obj interface{}) (gziped bool, data []byte, err error) {
	bs, err := json.Marshal(obj)
	if err != nil {
		return
	}
	if len(bs) <= GzipMinSize {
		return false, bs, nil
	}
	buf := &bytes.Buffer{}
	gzipWriter := gzip.NewWriter(buf)
	_, err = gzipWriter.Write(bs)
	gzipWriter.Close()
	if err != nil {
		return
	}
	return true, buf.Bytes(), nil
}

func fromGzipJSON(data []byte, obj interface{}) (err error) {
	buf := bytes.NewBuffer(data)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return
	}
	defer gzipReader.Close()
	bs, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return
	}
	return json.Unmarshal(bs, obj)
}
