package locker

import (
	"errors"
	"testing"
)

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
		if errors.Is(err, LOCKFILE_ACTIVE) {
			t.Logf("Init1 %s", err)
		} else {
			t.Errorf("Init() error = %s", err)
		}
	}
	t.Log("init 2")
	err = l.Init()
	if err != nil {
		if errors.Is(err, LOCKFILE_ACTIVE) {
			t.Logf("Init2 %s", err)
		} else {
			t.Errorf("Init2 error = %s", err)
		}
	}

	err = l.Remove()
	if err != nil {
		t.Logf("Remove() error = %s", err)
	}

	// t.Log("init 3 after remove()")
	// err = l.Init()
	// if err != nil {
	// 	t.Errorf("Lock() error = %s", err)
	// }
}
