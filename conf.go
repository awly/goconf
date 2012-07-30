package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
	"time"
)

const refresh_interval = 3

var (
	c     = &config{data: make(map[string]interface{})}
	cPath = "config.json"
)

func init() {
	// Refresh routine, reacts to changes in config
	go func() {
		t := time.Tick(refresh_interval * time.Second)
		for _ = range t {
			LoadConfig(cPath)
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
	err = ioutil.WriteFile(cPath, data, 0)
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
