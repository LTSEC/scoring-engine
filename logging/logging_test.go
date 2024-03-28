// run by using go test logging.go shared.go logging_test.go

package logging

import (
	"sync"
	"testing"
)

// Important -> Read Below

func TestLogging(t *testing.T) {
	var sema sync.WaitGroup
	sema.Add(3)
	CreateLogFile()
	StartLog()
	go func() {
		defer sema.Done()
		simulate()
	}()
	go func() {
		defer sema.Done()
		simulate()
	}()
	go func() {
		defer sema.Done()
		simulate()
	}()
	sema.Wait()
	StopLog()
	t.Log("Finished")
}

func simulate() {
	// Send a message into the channel
	LogMessage("This is a test message")
}

/* This Setup will NOT work on routines that spawn child routines. To fix that make sure you
   have synchronization in place (either a waitgroup or channel) that is able to know when
   a child routine finished executing. Otherwise, your main function will finish executing and
   the child routines will terminate prematurely and some info won't get logged.
*/
