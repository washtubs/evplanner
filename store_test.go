package evplanner

import (
	"sync"
	"testing"
)

func TestMutex(t *testing.T) {
	m := new(sync.Mutex)
	m.Lock()
	t.Log("locked once")
	m.Unlock()
	//m.Unlock()
}
