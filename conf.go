package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
	"time"
)

// Value used as config refresh interval (reading from file) if it was
// not specified in config.
const default_refresh_interval = 3

var (
	c     = &config{data: make(map[string]interface{})}
	cPath = "config.json"
)

func init() {
	go func() {
		for {
			LoadConfig(cPath)
			ri, err := Get("refresh_interval")
			if err != nil {
				ri = default_refresh_interval
			}
            ref_interval, ok := ri.(float64)
            if !ok {
                ref_interval = default_refresh_interval
            }
			time.Sleep(time.Duration(ref_interval) * time.Second)
		}
	}()
}

type config struct {
	sync.Mutex
	data map[string]interface{}
}

func Get(key string) (interface{}, error) {
	c.Lock()
	defer c.Unlock()
	if val, ok := c.data[key]; ok {
		return val, nil
	}
	return nil, errors.New("no such value: " + key)
}

func Put(key string, val interface{}) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val
	data, err := json.Marshal(c.data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cPath, data, 0666)
	return err
}

func Del(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.data, key)
}

func LoadConfig(path string) error {
	if path == "" {
		path = cPath
	} else {
		cPath = path
	}
	c.Lock()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	defer c.Unlock()
	if err := json.Unmarshal(data, &c.data); err != nil {
		return err
	}
	return nil
}
