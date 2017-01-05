# Golang configuration library

This library uses [spf13/viper](http://github.com/spf13/viper) configuration library 
to reading configuration from files dan environmnet variables.

## Usage

- Create configuration file
```yaml
# ./config.yaml
mysql:
    username: "a user"
    password: "a password"
    connections:
        - {name: "default", port: 3306}
        - {name: "test", port: 13306}
```

- Import the module
```golang
package main

import (
    "fmt"
    "github.com/salestock/ssource/shared-libs/go/config"
)
```

- Create struct with `config` tag
```golang
type Connection struct {
    Name string `config:"name"` // "name" config for every mysql.connections element
    Port int    `config:"port"` // "port" value for every mysql.connections element
}

type MySQL struct {
    Username    string       `config:"mysql.hostname"`
    Password    string       `config:"mysql.hort"`
    Connections []Connection `config:"mysql.connections"`
}
```

- Explicitly call `GetConfig` method to get configuration values
```golang
func main() {
    mysql := new(MySQL)
    config.GetConfig(mysql)
    fmt.Println(mysql.Username)
    fmt.Println(mysql.Password)
    fmt.Println(mysql.Connections)
    // Outputs
    // a user
    // a password
    // [{default, 1306}, {test, 13306}]
}
```

- The module will search for environment variables with prefix `SS_` followed by configuration path (all uppercased and any `.` in path must be replaced by `_`)

```sh
SS_MYSQL_USERNAME="test" go run main.go
```

will give output:

```
test
a password
[{default, 3306}, {test, 13306}]
```

## Test

Run `glide up` in the module directory to install depedency. Then run `go test` to test the module
