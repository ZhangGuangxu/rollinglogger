package rollinglogger

import (
	"testing"

	"github.com/ZhangGuangxu/misc"
)

func TestNewLogger(t *testing.T) {
	filename := "./warn.log"
	opts := Options{
		Level:    Warn,
		Filename: filename,
		MaxSize:  1,
		MaxAge:   7,
	}
	logger := NewLogger(opts)
	if logger == nil {
		t.Errorf("create logger failed, opts [%v]\n", opts)
		return
	}
	logger.Warn("some warning")
	if err := logger.Sync(); err != nil {
		t.Errorf("logger.Sync() error [%v]\n", err)
		return
	}
	if !misc.IsFileExist(filename) {
		t.Errorf("create logger failed, logger file [%s] does not exist\n", filename)
		return
	}
}
