//go:build windows
// +build windows

package winservice

import (
	"fmt"

	wsvc "golang.org/x/sys/windows/svc"
)

var chanGraceExit = make(chan int)

type serviceWindows struct{}

func init() {
	/*	interactive, err := wsvc.IsWindowsService()
		if err != nil {
			panic(err)
		}

		if interactive {
			return
		}*/
	go func() {
		err := wsvc.Run("CheckerReporter", serviceWindows{})
		if err != nil {
			fmt.Printf("Error Call wsvc.Run(CheckerReporter, serviceWindows{}) : %v \n", err)
		}
	}()
}

func (serviceWindows) Execute(args []string, r <-chan wsvc.ChangeRequest, s chan<- wsvc.Status) (svcSpecificEC bool, exitCode uint32) {
	const accCommands = wsvc.AcceptStop | wsvc.AcceptShutdown
	s <- wsvc.Status{State: wsvc.StartPending}
	s <- wsvc.Status{State: wsvc.Running, Accepts: accCommands}

	for {
		c := <-r
		switch c.Cmd {
		case wsvc.Interrogate:
			s <- c.CurrentStatus
		case wsvc.Stop, wsvc.Shutdown:
			chanGraceExit <- 1
			s <- wsvc.Status{State: wsvc.StopPending}
			return false, 0
		}
	}

}

func ShutdownChannel() <-chan int {
	return chanGraceExit
}
