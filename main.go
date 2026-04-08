package locker

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/shirou/gopsutil/process"
)

func main() {

}

type Locker struct {
	pid      int
	lockfile *os.File
	file     string
}

func New(lockfileName string, lockfilePath string) *Locker {
	return &Locker{
		file: path.Join(lockfilePath, lockfileName),
	}
}

func (l *Locker) Init() error {
	var err error

	l.pid = os.Getpid()

	l.lockfile, err = os.Open(l.file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			l.create()
			l.update()
			l.lockfile.Close()
		} else {
			return err
		}
		return nil
	}

	defer l.lockfile.Close()

	sc := bufio.NewScanner(l.lockfile)
	if sc.Scan() {
		line := sc.Text()

		num, err := strconv.Atoi(line)
		if err != nil {
			return err
			// os.Exit(1)
		}

		_, err = process.PidExists(int32(num))
		if err != nil {
			return err
			// os.Exit(1)
		}

		proc, err := process.NewProcess(int32(num))

		if err != nil {
			if err.Error() != "process does not exist" {
				return err
				// os.Exit(1)
			}
			l.create()
			l.update()

			return nil
		}

		name, err := proc.Name()
		if err != nil {
			return err
			// os.Exit(1)
		}

		isRunning, err := proc.IsRunning()
		if err != nil {
			return err
			// os.Exit(1)
		}

		if isRunning && name == l.file {
			fmt.Printf("%s, pid %d is running. Exiting\n", name, num)
			return err
			// os.Exit(1)
		}

		l.update()

	} else {
		fmt.Println("no pid number found")
	}
	return nil
}

// update writes the current process ID to the lockfile.
func (l *Locker) update() {
	if l.lockfile == nil {
		err := l.create()
		if err != nil {
			fmt.Println("a.Create error: no pid number found")
			return
		}
	}

	n, err := l.lockfile.Write([]byte(strconv.Itoa(l.pid)))
	if err != nil {
		fmt.Printf("lockFile.Write error %s\n", err)
		return
	}

	fmt.Printf("%d bytes written to lockfile\n", n)
}

// create creates a new lockfile
func (l *Locker) create() error {
	var err error
	l.lockfile, err = os.Create(l.file)
	if err != nil {
		return err
		// os.Exit(1)
	}
	return nil
}

func (l *Locker) Remove() error {
	return os.Remove(l.file)
}
