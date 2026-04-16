package locker

import (
	"errors"
	"testing"
)

func TestLock(t *testing.T) {
	var locker *Locker
	var err error

	filename := "locker.test"
	path := "/tmp"

	t.Log("New 1")

	locker, err = New(filename, path, true)
	if err != nil {
		if errors.Is(err, LOCKFILE_ACTIVE) {
			t.Logf("New 1 aborting %s", err)
		} else {
			t.Errorf("New 1 error = %s", err)
		}
	}

	t.Logf("New 2 %v", locker)

	locker, err = New(filename, path, true)
	if err != nil {
		if errors.Is(err, LOCKFILE_ACTIVE) {
			t.Logf("New 2 aborting %s", err)
		} else {
			t.Errorf("New 2 aborting error = %s", err)
		}
	}

	t.Logf("New 3 before Remove() %v", locker)

	err = locker.Remove()
	if err != nil {
		t.Errorf("Remove() error = %s", err)
	}

	t.Logf("New 3 after Remove() %v", locker)
	locker, err = New(filename, path, true)
	if err != nil {
		if errors.Is(err, LOCKFILE_ACTIVE) {
			t.Logf("New 3 aborting %s", err)
		} else {
			t.Errorf("New 3 aborting error = %s", err)
		}
	}
}
