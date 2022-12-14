package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type Config struct {
	configs map[string]interface{}
}

var config *Config

var once sync.Once

// New create a new or get the existing config instance (Singleton)
func New() *Config {
	// Check for the config if not initiated
	if config == nil {
		// Using once from sync library to make sure the config initiated once
		once.Do(func() { config = &Config{} })
	}

	return config
}

// AppendFiles append multiple config files data to the configs
func (config *Config) AppendFiles(paths ...string) {
	for _, path := range paths {
		config.AppendFile(path)
	}
}

// AppendFile append config file data to the configs
func (config *Config) AppendFile(path string) {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return
	}
	defer jsonFile.Close()

	newConfigs, err := readAndParseFile(jsonFile)
	if err != nil {
		log.Printf("read file %s: %s", path, err)
		return
	}

	config.configs = mergeConfigs(config.configs, newConfigs.(map[string]interface{}))
}

// Get get the config by key
func (config *Config) Get(key string) interface{} {
	if config.configs == nil {
		return nil
	}

	var keys = strings.Split(key, ".")
	return getValue(keys, config.configs)
}

// Keys get all the keys by specifying a key or empty string to get the keys of root
func (config *Config) Keys(key string) []string {
	if config.configs == nil {
		return nil
	}

	var configs interface{}

	if key != "" {
		var keys = strings.Split(key, ".")
		configs = getValue(keys, config.configs)
	} else {
		configs = config.configs
	}

	if fmt.Sprintf("%T", configs) != "map[string]interface {}" {
		return []string{}
	}

	return getKeys(configs.(map[string]interface{}))
}

// Get the keys of a map
func getKeys[K comparable, V any](in map[K]V) []K {
	result := make([]K, 0, len(in))

	for k := range in {
		result = append(result, k)
	}

	return result
}

// Get the config by array of the key that shows the depth
func getValue(keys []string, configs map[string]interface{}) interface{} {
	if keys == nil {
		return nil
	}

	if len(keys) < 2 {
		return configs[keys[0]]
	}

	if subConfigs, ok := configs[keys[0]]; ok {
		return getValue(keys[1:], subConfigs.(map[string]interface{}))
	}

	return nil
}

// Merging the existing config with the new configs
func mergeConfigs[K comparable, V any](maps ...map[K]V) map[K]V {
	result := map[K]V{}

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

func readAndParseFile(file io.Reader) (interface{}, error) {
	var newConfigs map[string]interface{}

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while reading the file")
	}

	err = json.Unmarshal(byteValue, &newConfigs)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration file format or structure")
	}

	return newConfigs, nil
}
