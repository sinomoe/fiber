package config

import (
	"errors"
	"os"
	"reflect"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	mu       sync.RWMutex
	registry = make(map[string]interface{})
)

func load(path string, config interface{}) error {
	if reflect.ValueOf(config).Kind() != reflect.Ptr {
		return errors.New("config should be a ptr to struct")
	}
	var (
		file *os.File
		err  error
	)
	if file, err = os.Open(path); err != nil {
		return err
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(config); err != nil {
		return err
	}
	return nil
}

func Load(path string, config interface{}) error {
	if err := load(path, config); err != nil {
		return err
	}
	mu.Lock()
	defer mu.Unlock()

	registry[path] = config
	return nil
}

func Get(path string) (interface{}, error) {
	mu.RLock()
	defer mu.RUnlock()
	if v, ok := registry[path]; ok {
		return v, nil
	}
	return nil, errors.New("config not exist")
}
