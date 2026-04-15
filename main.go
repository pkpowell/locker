package locker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/shirou/gopsutil/process"
)

var (
	LOCKFILE_ACTIVE            = errors.New("lockfile-active")
	LOCKFILE_PERMISSION_DENIED = errors.New("lockfile-permission-denied")
	LOCKFILE_BAD_PID           = errors.New("lockfile-bad-pid")
)

type Locker struct {
	lockfile *os.File
	pid      int
	file     string
	procname string
}

// New creates a new Locker instance with the given lockfile name and path. Runs init and returns the Locker and any error.
func New(procname string, lockfilePath string) (l *Locker, err error) {
	l = new(procname, lockfilePath)
	err = l.init()

	return l, err
}

// new creates a new Locker instance with the given lockfile name and path.
func new(procname string, lockfilePath string) *Locker {
	return &Locker{
		file:     path.Join(lockfilePath, procname),
		procname: procname,
	}
}

// init initializes the Locker by creating the lockfile if it doesn't exist,
// or checks if an existing lockfile contains a valid PID.
func (locker *Locker) init() (err error) {
	locker.pid = os.Getpid()

	locker.lockfile, err = os.OpenFile(locker.file, os.O_RDWR, 0777)
	if err != nil {
		// if the lockfile doesn't exist, create it and update it
		if errors.Is(err, os.ErrNotExist) {
			return locker.updatePID()
		}

		// abort on permissions error. Nothing we can do here
		if errors.Is(err, os.ErrPermission) {
			return LOCKFILE_PERMISSION_DENIED
		}
		return err
	}

	defer locker.lockfile.Close()

	var line string

	// read first line. Abort if there is more than one line. nil == empty, EOF == one line read
	scanner := bufio.NewReader(locker.lockfile)
	line, err = scanner.ReadString('\n')

	switch err {
	case nil:
	case io.EOF:

	default:
		return err
	}

	// update if line == empty string
	if line == "" {
		fmt.Println("warn: no data in lockfile")
		locker.updatePID()
	}

	// convert string to int. Abort if NaN
	num, err := strconv.Atoi(line)
	if err != nil {
		fmt.Printf("strconv.Atoi error for <%s>, %s", line, err)
		return LOCKFILE_BAD_PID
	}

	// get process belonging to pid
	proc, err := process.NewProcess(int32(num))
	// if the process does not exist check for other error. If ok create and update pid
	if err != nil {
		if errors.Is(err, process.ErrorProcessNotRunning) {
			return locker.updatePID()
		}

		return err
	}

	name, err := proc.Name()
	if err != nil {
		return err
	}

	isRunning, err := proc.IsRunning()
	if err != nil {
		return err
	}

	// check if the process is running and matches the lockfile name
	if isRunning && name == locker.procname {
		return LOCKFILE_ACTIVE
	}

	return locker.updatePID()
}

// updatePID writes the current process ID to the lockfile.
func (locker *Locker) updatePID() (err error) {
	if locker.lockfile == nil {
		err = locker.create()
		if err != nil {
			return fmt.Errorf("locker.Create error: %w", err)
		}
	}

	err = locker.lockfile.Truncate(0)
	if err != nil {
		return fmt.Errorf("locker.lockfile.Truncate error %w", err)
	}

	_, err = locker.lockfile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("locker.lockfile.Seek error %w", err)
	}

	n, err := locker.lockfile.Write([]byte(strconv.Itoa(locker.pid)))
	if err != nil {
		return fmt.Errorf("locker.lockfile.Write error %w", err)
	}

	fmt.Printf("%d bytes written to lockfile\n", n)
	return nil
}

// create creates a new lockfile
func (locker *Locker) create() (err error) {
	locker.lockfile, err = os.Create(locker.file)

	return err
}

func (locker *Locker) LockfileName() string {
	return locker.file
}

func (locker *Locker) Remove() (err error) {
	err = os.Remove(locker.file)
	if errors.Is(err, os.ErrPermission) {
		return LOCKFILE_PERMISSION_DENIED
	}
	return
}
