package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var testFileName = "config.yaml"

var testYaml = `test:
    hostname: "127.0.0.1"
    port: 3126
    overridden: ""
    map:
        - { name: "example", port: 3127 }
`

type testStringConfig struct {
	SomeHostname string `config:"test.hostname"`
	Port         int    `config:"test.port"`
}

type testMapConfig struct {
	SomeMap []interface{} `config:"test.map"`
}

func TestStringConfig(t *testing.T) {
	c := new(testStringConfig)
	GetConfig(c)
	if c.SomeHostname != "127.0.0.1" {
		t.Errorf("c.SomeHostname == %q, want %q", c.SomeHostname, "127.0.0.1")
	}
	if c.Port != 3126 {
		t.Errorf("c.Port == %d, want %d", c.Port, 3126)
	}
}

func TestMapConfig(t *testing.T) {
	c := new(testMapConfig)
	GetConfig(c)
	if len(c.SomeMap) == 0 {
		t.Errorf("Failed to get config value of `c.Somemap`")
	}
	for i, elmt := range c.SomeMap {
		switch obj := elmt.(type) {
		case map[interface{}]interface{}:
			if v, ok := obj["name"]; !ok || v != "example" {
				t.Errorf("map[%d][%q] = \"%q\", want %q", i, "name", v, "example")
			}
			if v, ok := obj["port"]; !ok || v != 3127 {
				t.Errorf("map[%d][%q] = %d, want %d", i, "port", v, 3127)
			}
		}
	}
}

func TestMain(m *testing.M) {
	bytes := []byte(testYaml)
	err := ioutil.WriteFile(testFileName, bytes, 0777)
	os.Setenv("SS_TEST_OVERRIDDEN", "test")
	if err == nil {
		result := m.Run()
		os.Remove(testFileName)
		os.Exit(result)
	} else {
		panic("cannot create temp config file needed for testing")
	}
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
