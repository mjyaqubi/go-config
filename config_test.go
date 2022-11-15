package config

import (
	"strings"
	"testing"
)

var jsonExample1 = `{
	"environment": "production",
	"database": {
	  "host": "mysql",
	  "port": 3306,
	  "username": "gholi",
	  "password": "gholi"
	},
	"cache": {
	  "redis": {
		"host": "redis",
		"port": 6379
	  }
	}
}`

var jsonExample2 = `{
	"environment": "development",
	"database": {
	  "host": "127.0.0.1",
	  "port": 3306,
	  "username": "gholi",
	  "password": "gholi"
	},
	"cache": {
	  "redis": {
		"host": "127.0.0.1",
		"port": 6379
	  }
	}
}`

var configExample = map[string]interface{}{
	"a": map[string]interface{}{"aa": "aa"},
	"b": map[string]interface{}{"bb": "bb"},
}

var configEmptyExample = make(map[string]interface{})

func TestـNew(t *testing.T) {
	t.Run("create new config instance", func(t *testing.T) {
		config := New()

		if config == nil {
			t.Fatal(`the config expected not to be {nil}`)
		}
	})
}

func Test_AppendFiles(t *testing.T) {
	t.Run("nil as argument", func(t *testing.T) {
		config := New()

		config.AppendFiles()

		if config.configs != nil {
			t.Fatal(`the config.configs expected to be {nil}`)
		}
	})

	t.Run("empty string as argument", func(t *testing.T) {
		config := New()

		config.AppendFiles("")

		if config.configs != nil {
			t.Fatal(`the config.configs expected to be {nil}`)
		}
	})

	t.Run("invalid file path as argument", func(t *testing.T) {
		config := New()

		config.AppendFiles("notExistConfig.json")

		if config.configs != nil {
			t.Fatal(`the config.configs expected to be {nil}`)
		}
	})
}

func Test_AppendFile(t *testing.T) {
	t.Run("empty string as argument", func(t *testing.T) {
		config := New()

		config.AppendFile("")

		if config.configs != nil {
			t.Fatal(`the config.configs expected to be {nil}`)
		}
	})

	t.Run("invalid file path as argument", func(t *testing.T) {
		config := New()

		config.AppendFile("notExistConfig.json")

		if config.configs != nil {
			t.Fatal(`the config.configs expected to be {nil}`)
		}
	})

	t.Run("appending valid json data", func(t *testing.T) {
		config := New()

		newConfig, _ := readAndParseFile(strings.NewReader(jsonExample1))
		config.configs = mergeConfigs(config.configs, newConfig.(map[string]interface{}))

		if config.configs == nil {
			t.Fatal(`the config.configs expected not to be {nil}`)
		}
	})
}

func Test_Get(t *testing.T) {
	t.Run("empty config example and empty key", func(t *testing.T) {
		config := New()
		config.configs = configEmptyExample

		value := config.Get("")

		if value != nil {
			t.Fatal(`the value expected to be {nil} but it isn't`)
		}
	})

	t.Run("config example and empty key", func(t *testing.T) {
		config := New()
		config.configs = configExample

		value := config.Get("")

		if value != nil {
			t.Fatal(`the value expected to be {nil} but it isn't`)
		}
	})

	t.Run("config example and a key", func(t *testing.T) {
		config := New()
		config.configs = configExample

		value := config.Get("a")

		if value == nil {
			t.Fatal(`the value expected not to be {nil} but it is`)
		}
	})

	t.Run("config example and a nested key", func(t *testing.T) {
		config := New()
		config.configs = configExample

		value := config.Get("a.aa")

		if value != "aa" {
			t.Fatalf(`the value expected to be {aa} but actual value is {%s}`, value)
		}
	})
}

func Test_Keys(t *testing.T) {
	t.Run("empty config example and empty key", func(t *testing.T) {
		config := New()
		config.configs = configEmptyExample

		keys := config.Keys("")

		if len(keys) > 0 {
			t.Fatalf(`the keys length expected to be {0} but actual length is {%d}`, len(keys))
		}
	})

	t.Run("config example and empty key", func(t *testing.T) {
		config := New()
		config.configs = configExample

		keys := config.Keys("")

		if len(keys) != 2 {
			t.Fatalf(`the keys length expected to be {2} but actual length is {%d}`, len(keys))
		}
	})

	t.Run("config example and a key", func(t *testing.T) {
		config := New()
		config.configs = configExample

		keys := config.Keys("a")

		if len(keys) != 1 {
			t.Fatalf(`the keys length expected to be {1} but actual length is {%d}`, len(keys))
		}
	})

	t.Run("config example and a nested key", func(t *testing.T) {
		config := New()
		config.configs = configExample

		keys := config.Keys("a.aa")

		if len(keys) != 0 {
			t.Fatalf(`the keys length expected to be {0} but actual length is {%d}`, len(keys))
		}
	})
}

func Test_getKeys(t *testing.T) {
	t.Run("empty config example", func(t *testing.T) {
		keys := getKeys(configEmptyExample)

		if len(keys) > 0 {
			t.Fatalf(`the keys length expected to be {0} but actual length is {%d}`, len(keys))
		}
	})

	t.Run("config example", func(t *testing.T) {
		keys := getKeys(configExample)

		if len(keys) != 2 {
			t.Fatalf(`the keys length expected to be {2} but actual length is {%d}`, len(keys))
		}
	})
}

func Test_getValue(t *testing.T) {
	t.Run("nil keys", func(t *testing.T) {
		value := getValue(nil, configEmptyExample)

		if value != nil {
			t.Fatalf(`the returned value expected to be {nil} but actual value is {%s}`, value.(string))
		}
	})

	t.Run("valid keys with empty config", func(t *testing.T) {
		value := getValue(nil, configEmptyExample)

		if value != nil {
			t.Fatalf(`the returned value expected to be {nil} but actual value is {%s}`, value.(string))
		}
	})

	t.Run("existing key in the config", func(t *testing.T) {
		value := getValue([]string{"a", "aa"}, configExample)

		if value != "aa" {
			t.Fatalf(`the returned value expected to be {aa} but actual value is {%s}`, value.(string))
		}
	})
}

func Test_mergeConfigs(t *testing.T) {
	t.Run("one map as argument", func(t *testing.T) {
		config := make(map[string]interface{})
		config["a"] = "b"

		config = mergeConfigs(config)

		if config["a"] != "b" {
			t.Fatalf(`the config["a"] expected to be {b} but actual value is {%s}`, config["a"].(string))
		}
	})

	t.Run("a map as argument multiple time", func(t *testing.T) {
		config := make(map[string]interface{})
		config["a"] = "b"

		config = mergeConfigs(config, config, config, config)

		if config["a"] != "b" {
			t.Fatalf(`the config["a"] expected to be {b} but actual value is {%s}`, config["a"].(string))
		}
	})

	t.Run("two map as argument", func(t *testing.T) {
		config := make(map[string]interface{})
		config["a"] = "b"

		secondConfig := make(map[string]interface{})
		secondConfig["c"] = "d"

		config = mergeConfigs(config, secondConfig)

		if config["a"] != "b" || config["c"] != "d" {
			t.Fatalf(`the config["a"] expected to be {b} or config["c"] expected to be {d} but actual value is {%s} and {%s}`, config["a"].(string), config["b"].(string))
		}
	})
}

func Test_readAndParseFile(t *testing.T) {
	t.Run("invalid json format", func(t *testing.T) {
		_, err := readAndParseFile(strings.NewReader("This is not a valid JSON file"))
		if err == nil {
			t.Fatal(`the error expected not to be {nil} but it is`)
		}
	})

	t.Run("valid json format", func(t *testing.T) {
		config, err := readAndParseFile(strings.NewReader(jsonExample1))
		if err != nil {
			t.Fatal(`the error expected to be {nil} but it isn't`)
		}

		if config == nil {
			t.Fatal(`the config expected not to be {nil} but it is`)
		}
	})
}
