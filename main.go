package locker

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/shirou/gopsutil/process"
)

func main() {

}

type Locker struct {
	pid          int
	lockfile     *os.File
	LockfileName string
}

func (a *Locker) checkPID() error {
	var err error

	a.pid = os.Getpid()

	a.lockfile, err = os.Open(a.LockfileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			a.Create()
			a.Update()
			a.lockfile.Close()
		} else {
			return err
		}
		return nil
	}

	defer a.lockfile.Close()

	sc := bufio.NewScanner(a.lockfile)
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
			a.Create()
			a.Update()

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

		if isRunning && name == a.LockfileName {
			// Warnf("%s, pid %d is running. Exiting", name, num)
			return err
			// os.Exit(1)
		}

		a.Update()

	} else {
		fmt.Println("no pid number found")
	}
	return nil
}

func (a *Locker) Update() {
	if a.lockfile == nil {
		err := a.Create()
		if err != nil {

		}
	}

	n, err := a.lockfile.Write([]byte(strconv.Itoa(a.pid)))
	if err != nil {
		fmt.Printf("lockFile.Write error %s\n", err)
		return
	}

	fmt.Printf("%d bytes written to lockfile\n", n)
}

func (a *Locker) Create() error {
	var err error
	a.lockfile, err = os.Create(a.LockfileName)
	if err != nil {
		return err
		// os.Exit(1)
	}
	return nil
}

func (a *Locker) Remove() error {
	err := os.Remove(a.LockfileName)
	if err != nil {
		return err
		// os.Exit(1)
	}
	return nil
}
