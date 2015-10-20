package base

import (
	"testing"
	"time"
)

func TestItem(t *testing.T) {
	t.Log("test item")
}

func BenchmarkItem(b *testing.B) {
	customTimerTag := false
	if customTimerTag {
		b.Log("stop timer")
		b.StopTimer()
	}

	time.Sleep(time.Second)
	if customTimerTag {
		b.Log("start timer")
		b.StartTimer()
	}
}
