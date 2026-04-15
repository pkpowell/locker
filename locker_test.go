package locker

import (
	"errors"
	"testing"
)

func TestLock(t *testing.T) {
	t.Log("init 1")
	l, err := New("locker.test", "/tmp")
	if err != nil {
		if errors.Is(err, LOCKFILE_ACTIVE) {
			t.Logf("Init1 aborting %s", err)
		} else {
			t.Errorf("Init() error = %s", err)
		}
	}

	t.Log("init 2")
	l, err = New("locker.test", "/tmp")
	if err != nil {
		if errors.Is(err, LOCKFILE_ACTIVE) {
			t.Logf("Init2 aborting %s", err)
		} else {
			t.Errorf("Init2 aborting error = %s", err)
		}
	}

	err = l.Remove()
	if err != nil {
		t.Logf("Remove() error = %s", err)
	}

	t.Log("init 3 after remove()")
	l, err = New("locker.test", "/tmp")
	if err != nil {
		t.Errorf("Lock() error = %s", err)
	}
}
