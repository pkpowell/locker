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

type Locker struct {
	pid      int
	lockfile *os.File
	file     string
}

// New creates a new Locker instance with the given lockfile name and path.
func New(lockfileName string, lockfilePath string) *Locker {
	return &Locker{
		file: path.Join(lockfilePath, lockfileName),
	}
}

// Init initializes the Locker by creating the lockfile if it doesn't exist,
// or checking if the existing lockfile contains a valid PID.
func (l *Locker) Init() error {
	var err error

	l.pid = os.Getpid()

	l.lockfile, err = os.OpenFile(l.file, os.O_RDWR, 0644)
	if err != nil {
		// if the lockfile doesn't exist, create it
		if errors.Is(err, os.ErrNotExist) {
			// defer l.lockfile.Close()

			// update pid in lockfile
			return l.updatePID()
		}
		return err
	}

	defer l.lockfile.Close()

	var line string
	scanner := bufio.NewReader(l.lockfile)
	line, err = scanner.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Println("EOF")
			fmt.Println(line) // Print last line if not empty
		} else {
			return fmt.Errorf("error reading from file: %w", err)
		}
	}

	if len(line) == 0 {
		return fmt.Errorf("error no data in lockfile")
	}

	num, err := strconv.Atoi(line)
	if err != nil {
		return err
	}

	fmt.Printf("lockfile pid %d, current pid %d\n", num, l.pid)

	// check if the process is currently running
	_, err = process.PidExists(int32(num))
	if err != nil {
		return err
	}

	proc, err := process.NewProcess(int32(num))
	// if the process does not exist check for other error. If ok create and update pid
	if err != nil {
		if err.Error() != "process does not exist" {
			return err
		}

		return l.updatePID()
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
	if isRunning && name == l.file {
		fmt.Printf("%s, pid %d is running. Exiting\n", name, num)
		return err
	}

	return l.updatePID()

	// } else {
	// 	fmt.Printf("no pid number found %s\n", sc.Text())
	// }
	// return nil
}

// updatePID writes the current process ID to the lockfile.
func (l *Locker) updatePID() (err error) {
	if l.lockfile == nil {
		err = l.create()
		if err != nil {
			fmt.Println("a.Create error: no pid number found")
			return err
		}
	}

	n, err := l.lockfile.Write([]byte(strconv.Itoa(l.pid)))
	if err != nil {
		fmt.Printf("lockFile.Write error %s\n", err)
		return err
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

func (l *Locker) Remove() error {
	return os.Remove(l.file)
}
