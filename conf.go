/*
Package goconf implements simple concurency-safe configuration tool.
All configuration parameters are read from json encoded file.

Default file name is "config.json", but can be changed easily (see LoadConfig)
File is re-read periodically, with default interval of 5 seconds.
It can be configured via "refresh" parameter in the same config file.

Example config.json is located in the repository.

Created as a helper library for my other projects, not tested well,
use at your own risk, no guarantees!
*/
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
			ri, err := Get("refresh")
			if err != nil {
				log.Println("'refresh' not set, unsing default :", default_refresh_interval)
				ri = default_refresh_interval
			}
			ref_interval, ok := ri.(float64)
			if !ok {
				log.Println("invalid 'refresh', unsing default :", default_refresh_interval)
				ref_interval = default_refresh_interval
			}
			time.Sleep(time.Duration(ref_interval) * time.Second)
		}
	}()
}

//Function Get is used to retrieve values from config.
//Can recieve multiple args in which case it tries to get deeper
//into the config json structure.
func Get(keys ...string) (res interface{}, err error) {
	c.Lock()
	defer c.Unlock()
	res = c.data
	for i := 0; i < len(keys); i++ {
		t, ok := map[string]interface{}{}, false
		if t, ok = res.(map[string]interface{}); !ok {
			return nil, errors.New(fmt.Sprint("can't get value for ", keys))
		}
		if res, ok = t[keys[i]]; !ok {
			return nil, errors.New(fmt.Sprint("can't get value for ", keys))
		}
	}
	return
}

//Function LoadConfig reloads configuration data from file named by path
//path can be set to "" (empty string) in which case the last used file name
//is kept. The last successive file is always used for auto-update.
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
	if err := json.Unmarshal(data, &c.data); err != nil {
		return err
	}
	c.path = path
	return nil
}
