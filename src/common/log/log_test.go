package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sync"
)

func TestNewDefaultLog(t *testing.T) {
	a := assert.New(t)
	l, err := NewDaliyRotateLog("/tmp/test.go")
	a.NoError(err, "new time handler should no error")
	a.NotNil(l, "NewDefault log should not nil")
	l.LogInfo("test info")
	l.LogDebug("test debug")
	l.LogWarn("test Warn")
	l.LogError("test Error")
	l.InfoJson(a)
	l.Close()
}

func TestLogMultiClose(t *testing.T) {
	l, err := NewDaliyRotateLog("/tmp/test_log_mul_close")
	assert.Nil(t, err)

	wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		go func() {
			wg.Add(1)
			l.Close()
			wg.Done()
		}()
	}
	wg.Wait()
}
