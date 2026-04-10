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

// New creates a new Locker instance with the given lockfile name and path.
func New(procname string, lockfilePath string) *Locker {
	return &Locker{
		file:     path.Join(lockfilePath, procname),
		procname: procname,
	}
}

// Init initializes the Locker by creating the lockfile if it doesn't exist,
// or checking if the existing lockfile contains a valid PID.
func (l *Locker) Init() error {
	var err error

	l.pid = os.Getpid()

	l.lockfile, err = os.OpenFile(l.file, os.O_RDWR, 0777)
	if err != nil {
		// if the lockfile doesn't exist, create it
		if errors.Is(err, os.ErrNotExist) {
			// update pid in lockfile
			return l.updatePID()
		}

		// abort on permissions error
		if errors.Is(err, os.ErrPermission) {
			return LOCKFILE_PERMISSION_DENIED
		}
		return err
	}

	defer l.lockfile.Close()

	var line string

	scanner := bufio.NewReader(l.lockfile)
	line, err = scanner.ReadString('\n')
	switch err {
	case nil:
	case io.EOF:

	default:
		return err
	}

	if line == "" {
		fmt.Println("warn: no data in lockfile")
		l.updatePID()
	}

	num, err := strconv.Atoi(line)
	if err != nil {
		fmt.Printf("strconv.Atoi error for <%s>, %s", line, err)
		return LOCKFILE_BAD_PID
	}

	proc, err := process.NewProcess(int32(num))
	// if the process does not exist check for other error. If ok create and update pid
	if err != nil {
		if errors.Is(err, process.ErrorProcessNotRunning) {
			// fmt.Printf("process.NewProcess error %s\n", err)
			return l.updatePID()
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
	if isRunning && name == l.procname {
		return LOCKFILE_ACTIVE
	}

	return l.updatePID()
}

// updatePID writes the current process ID to the lockfile.
func (l *Locker) updatePID() (err error) {
	if l.lockfile == nil {
		err = l.create()
		if err != nil {
			return fmt.Errorf("l.Create error: %w", err)
		}
	}

	err = l.lockfile.Truncate(0)
	if err != nil {
		return fmt.Errorf("Truncate error %w", err)
	}

	_, err = l.lockfile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("Seek error %w", err)
	}

	n, err := l.lockfile.Write([]byte(strconv.Itoa(l.pid)))
	if err != nil {
		return fmt.Errorf("Write error %w", err)
	}

	fmt.Printf("%d bytes written to lockfile\n", n)
	return nil
}

// create creates a new lockfile
func (l *Locker) create() error {
	var err error
	l.lockfile, err = os.Create(l.file)

	return err
}

func (l *Locker) Remove() (err error) {
	err = os.Remove(l.file)
	if errors.Is(err, os.ErrPermission) {
		return LOCKFILE_PERMISSION_DENIED
	}
	return
}
