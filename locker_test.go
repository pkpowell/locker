package locker

import "testing"

func TestLock(t *testing.T) {
	l := New("test", "/tmp")

	t.Log("init 1")
	err := l.Init()
	if err != nil {
		t.Errorf("Init() error = %v", err)
	}

	t.Log("init 2")
	err = l.Init()
	if err != nil {
		t.Errorf("Init() error = %v", err)
	}

	l.Remove()

	t.Log("init 3 after remove()")
	err = l.Init()
	if err != nil {
		t.Errorf("Lock() error = %v", err)
	}
}
