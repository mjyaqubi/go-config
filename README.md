# go-config
The **go-config** is solution for Go application to load config files and access to them by config key. The config files will override as they appending to the existing configs.

## Install

```
go get github.com/mjyaqubi/go-config
```

**Note:** go-config uses [Go Modules](https://github.com/golang/go/wiki/Modules) to manage dependencies.


## Usage
### Create an instance of the config
```
import "github.com/mjyaqubi/go-config"

var config = config.New()
```

### Append a config file
```
config.AppendFile("config.json")
```
### Append multiple config files

```
config.AppendFiles("config.json", "config2.json", ...)
```
### Get the config by key
```
config.Get("database.host")
```
### Get the keys of config
```
config.Keys("") // root keys
config.Keys("database")
```

## Example
```
package main

import (
    "fmt"

    "github.com/mjyaqubi/go-config"
)

func main() {
	config := config.New()

	config.AppendFiles("fixtures/config.json", "fixtures/config.local.json")

	fmt.Println(config.Keys(""))
    fmt.Println(config.Keys("database"))

	fmt.Println(config.Get("database"))
	fmt.Println(config.Get("database.host"))
	fmt.Println(config.Get("something"))
}

```