//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a process
	proc := MockProcess{}

	sigsChan := make(chan os.Signal, 1)

	// Listen for SIGINT
	signal.Notify(sigsChan, syscall.SIGINT)

	go func() {
		<-sigsChan
		fmt.Println("Gracefully stopping the program.")
		go proc.Stop()
		fmt.Println("Waiting to forcefully terminate the program.")
		<-sigsChan
		fmt.Println("Forcefully terminating the program.")
		os.Exit(1)
	}()

	// Run the process (blocking)
	proc.Run()
}
