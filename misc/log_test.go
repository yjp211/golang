package misc

import (
	"testing"
)

func TestLog(t *testing.T) {
	path := "/tmp/log"
	level := "DEBUG"
	Log, err := NewLog(path, level)
	if err != nil {
		t.Errorf("create log  at<%s> error: %s\n", path, err)
	}
	Log.Debug("aaaa")
	Log.Info("info")
	Log.Error("error")
}
