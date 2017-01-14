package goconf

import "testing"

func TestFromFile(t *testing.T) {
	type testStringConfig struct {
		SomeHostname string  `config:"test.hostname"`
		Port         float64 `config:"test.port"`
	}
	c := new(testStringConfig)
	FromFile("config.json", "APP")
	GetConfig(c)
	if c.SomeHostname != "127.0.0.1" {
		t.Errorf("c.SomeHostname == %q, want %q", c.SomeHostname, "127.0.0.1")
	}
	if c.Port != 3126 {
		t.Errorf("c.Port == %d, want %d", c.Port, 3126)
	}
}
