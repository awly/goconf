package goconf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

const default_refresh_interval = 5

var (
	c = &config{
		data: make(map[string]interface{}),
		path: "config.json",
	}
)

type config struct {
	sync.Mutex
	data map[string]interface{}
	path string
}

func init() {
	go func() {
		for {
			LoadConfig(c.path)
			ri, err := Get("refresh_interval")
			if err != nil {
				log.Println("'refresh' not set, unsing default :", default_refresh_interval)
				ri = default_refresh_interval
			}
			ref_interval, ok := ri.(int)
			if !ok {
				log.Println("invalid 'refresh', unsing default :", default_refresh_interval)
				ref_interval = default_refresh_interval
			}
			time.Sleep(time.Duration(ref_interval) * time.Second)
		}
	}()
}

func Get(keys ...string) (res interface{}, err error) {
	c.Lock()
	defer c.Unlock()
	res = c.data
	for i := 0; i < len(keys); i++ {
        t, ok := map[string]interface{}{}, false
		if t, ok = res.(map[string]interface{}); !ok {
			return nil, errors.New(fmt.Sprint("can't get vale for", keys))
		}
		if res, ok = t[keys[i]]; !ok {
			return nil, errors.New(fmt.Sprint("can't get vale for", keys))
		}
	}
	return
}

func LoadConfig(path string) error {
	if path == "" {
		path = c.path
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
        log.Println("failed to read config file :", err)
		return err
	}
	c.Lock()
	defer c.Unlock()
	c.path = path
	return json.Unmarshal(data, &c.data)
}
