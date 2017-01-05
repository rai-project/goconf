package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var testYamlFile = "config.yaml"

var testJsonFile = "config.json"

var testYaml = `test:
    hostname: "127.0.0.1"
    port: 3126
    overridden: ""
    map:
        - { name: "example", port: 3127 }
`

var testJson = `{
    "test": {
	"hostname": "127.0.0.1",
	"port": 3126,
	"overridden": "",
	"map": [
	    { "name": "example", "port": 3127 }
	]
    }
}`

func TestMain(m *testing.M) {
	os.Setenv("CONF_TEST_OVERRIDDEN", "test")
	bytes := []byte(testYaml)
	err := ioutil.WriteFile(testYamlFile, bytes, 0777)
	if err != nil {
		panic("cannot create temp yaml config file needed for testing")
	}
	bytes = []byte(testJson)
	err = ioutil.WriteFile(testJsonFile, bytes, 0777)
	if err != nil {
		panic("cannot create temp json config file needed for testing")
	}
	m.Run()
	os.Remove(testYamlFile)
	os.Remove(testJsonFile)
}

func ExampleGetConfig() {
	// ./config.yaml
	// test:
	//    hostname: "127.0.0.1"
	//    port: 3126
	//    map:
	//        - { name: "example", port: 3127 }
	type SubConfig struct {
		Name    string `config:"name"`
		Port    int    `config:"port"`
		Invalid string `config:"invalid"`
	}
	type StringConfig struct {
		SomeHostname string      `config:"test.hostname"`
		Port         int         `config:"test.port"`
		Overriden    string      `config:"test.overridden"`
		SomeMap      []SubConfig `config:"test.map"`
	}
	type InvalidConfig struct {
		Invalid string `config:"some.invalid.config"`
	}
	type InvalidMapConfig struct{}
	c := StringConfig{SomeMap: make([]SubConfig, 1)}
	Setup("yaml", "CONF") // Or Setup("json")
	GetConfig(&c)
	fmt.Println(c.SomeHostname)
	fmt.Println(c.Port)
	fmt.Println(c.Overriden)
	fmt.Println(c.SomeMap)
	fmt.Println(c.SomeMap[0].Invalid)
	i := InvalidConfig{}
	GetConfig(&i)
	fmt.Printf(i.Invalid)
	// Output:
	// 127.0.0.1
	// 3126
	// test
	// [{example 3127 }]
}
