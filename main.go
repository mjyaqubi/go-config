package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var configs map[string]interface{}

func main() {

	loadConfigs([]string{"config/config.json", "config/config.local.json"})

	fmt.Println(configs)

	fmt.Println(get("database.host"))
	fmt.Println(get("something"))
	fmt.Println(get("something.ABC"))
	fmt.Println(get("something.1"))
	fmt.Println(get("something.2"))

}

func loadConfigs(paths []string) {
	for _, s := range paths {
		loadConfig(s)
	}
}

func loadConfig(path string) {
	var config map[string]interface{}

	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	configs = mergeConfig(configs, config)
}

func mergeConfig[K comparable, V any](maps ...map[K]V) map[K]V {
	out := map[K]V{}

	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}

	return out
}

func get(key string) interface{} {
	var keySplited = strings.Split(key, ".")

	return getSub(keySplited, configs)
}

func getSub(keys []string, subConfig map[string]interface{}) interface{} {
	if len(keys) < 2 {
		return subConfig[keys[0]]
	}

	if val, ok := subConfig[keys[0]]; ok {
		return getSub(keys[1:], val.(map[string]interface{}))
	}

	return nil
}
