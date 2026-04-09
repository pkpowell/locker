package locker

import "testing"

func TestLock(t *testing.T) {
	l := New("locker.test", "/tmp")
	// defer func() {
	// 	err := l.Remove()
	// 	if err != nil {
	// 		t.Errorf("Remove() error = %s", err)
	// 	}
	// }()

	t.Log("init 1")
	err := l.Init()
	if err != nil {
		t.Errorf("Init() error = %s", err)
	}

	t.Log("init 2")
	err = l.Init()
	t.Log("init 2", err)
	if err != nil {
		t.Errorf("Init() error = %s", err)
	}

	l.Remove()

	t.Log("init 3 after remove()")
	err = l.Init()
	if err != nil {
		t.Errorf("Lock() error = %s", err)
	}
}
