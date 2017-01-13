package goconf

import (
	"testing"
)

func TestYamlStringConfig(t *testing.T) {
	type testStringConfig struct {
		SomeHostname string `config:"test.hostname"`
		Port         int    `config:"test.port"`
	}
	c := new(testStringConfig)
	Setup("yaml", "CONF")
	GetConfig(c)
	if c.SomeHostname != "127.0.0.1" {
		t.Errorf("c.SomeHostname == %q, want %q", c.SomeHostname, "127.0.0.1")
	}
	if c.Port != 3126 {
		t.Errorf("c.Port == %d, want %d", c.Port, 3126)
	}
}

func TestYamlMapConfig(t *testing.T) {
	type testMapConfig struct {
		SomeMap []interface{} `config:"test.map"`
	}
	c := new(testMapConfig)
	Setup("yaml", "CONF")
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
