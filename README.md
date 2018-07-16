# go-config

⚠️ This package is not mature and has not been extensively tested or used. This package was developed quickly for use in my projects. Bugs are likely.

## Introduction
This is a go package intended to support hierarchical configurations for application deployments. It is modeled after [lorenwest/node-config](https://github.com/lorenwest/node-config).

## Quickstart

1. Create your configuration files.

```
# config/default.yaml
plumber: mario
tool: wrench

# config/development.yaml
tool: sledgehammer

# config/production.yaml
plumber: luigi
```

2. Create a struct for your configuration.

```go
type AppConfig struct {
    Plumber string
    Tool string
}
```

3. Load the configs from your app

```go
import (
    "os"
    "github.com/alex-whitney/go-config"
)

// Example
func getConfig() (*AppConfig, error) {
    env := os.Getenv("ENV")
    if env == "" {
        env = "development"
    }

    appConfig := &AppConfig{}
    err := config.Load(&config.Options{
        Environment: env,
        Directory: "config",
    }, appConfig)
    
    return appConfig, err
}
```

## Environment Variables
Overrides can be set using Environment Variables with the `custom-environment-variables` file. These overrides are applied last.

```
# app/config/custom-environment-variables.yaml
foo: ENV_FOO
bar: BAR
```

```go
// app/main.go

package main

import (
    "fmt"
    "github.com/alex-whitney/go-config"
)

type AppConfig struct {
    Foo string
    Bar string
}

func main() {
    appConfig := &AppConfig{}
    config.Load(&config.Options{
        Directory: "config",
    }, appConfig)

    fmt.Println(appConfig.Foo)
    fmt.Println(appConfig.Bar)
}
```

```
$ go install
$ ENV_FOO=abc BAR=def app
abc
def
```

## Load order
Configurations are applied in the following order.

```
default.yaml
{environment}.yaml
local.yaml
local-{environment}.yaml
custom-environment-variables.yaml
```
