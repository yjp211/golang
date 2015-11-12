package misc

import (
	"testing"
)

func TestConfig(t *testing.T) {
	configPath := "/tmp/config"
	key := "Hello"
	Config, err := NewConfig(configPath)
	if err != nil {
		t.Errorf("load config<%s> error: %s\n", configPath, err)
	}
	val := Config.GetString(key, "")
	if val == "" {
		t.Errorf("no key for <%s>\n", key)
	}
	ints, _ := Config.GetStrings("WC", []string{})
	if len(ints) == 0 {
		t.Errorf("no key for <%s>\n", "WC")
	}
	t.Logf("wc:%v", ints)
}
