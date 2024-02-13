package main

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"time"
)

// StreamCmdOutput streams the output of a command to the logs
func StreamCmdOutput(c *cmd.Cmd) {
	// Print STDOUT and STDERR lines streaming from Cmd
	// N.B looks like this may not get used and may just be a vestige of copied coe
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		// Done when both channels have been closed
		// https://dave.cheney.net/2013/04/30/curious-channels
		for c.Stdout != nil || c.Stderr != nil {
			select {
			case line, open := <-c.Stdout:
				if !open {
					c.Stdout = nil
					continue
				}
				fmt.Printf("STDOUT: %s\n", line)
			case line, open := <-c.Stderr:
				if !open {
					c.Stderr = nil
					continue
				}
				fmt.Printf("STDERR: %s\n", line)
			}
		}
	}()
}

// startFailedCommand is intended to start a long running command
// but it will fail because the command doesn't exist
func startFailedCommand() *cmd.Cmd {
	opts := cmd.Options{
		//Buffered:  true,
		Streaming: true,
	}
	// This should fail with an error message indicating an unrecognized command
	cmd := cmd.NewCmdOptions(opts, "tail", "--notvalid")
	StreamCmdOutput(cmd)
	cmd.Start()

	// Wait for the process to actually have started
	for {
		s := cmd.Status()
		if s.StartTs > 0 {
			break
		}
		time.Sleep(1)
	}

	return cmd
}

func wrongWay() {
	// Calling stop and then waiting for Done appears to prevent the error message from being printed out
	// my suspicion is that startFailedCommand actually starts the command asynchronously. So its possible
	// we are calling stop before the command has actually started and as a result we never generate the output.
	cmd := startFailedCommand()
	cmd.Stop()
	<-cmd.Done()
}

func rightWay() {
	cmd := startFailedCommand()
	// We need to wait for the command to finish otherwise stdout/stderr won't get flushed and we might miss the error message
	<-cmd.Done()
}

func rightWayWithStop() {
	cmd := startFailedCommand()

	// Wait for cmd to start
	for {
		s := cmd.Status()
		if s.StartTs > 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	go func() {
		// Give the command time to start
		time.Sleep(1 * time.Second)
		s := cmd.Status()
		if s.Complete {
			fmt.Printf("Command completed before we could stop it\n")
			return
		}
		fmt.Printf("Stopping command\n")
		if err := cmd.Stop(); err != nil {
			fmt.Sprintf("Error stopping command: %s\n", err)
		}
	}()

	//if s.StopTs > 0 {
	//	fmt.Printf("Command already stopped\n")
	//} else {
	//	fmt.Printf("Command still running\n")
	//	cmd.Stop()
	//}

	<-cmd.Done()
}

func main() {
	rightWayWithStop()
}
