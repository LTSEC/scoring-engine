package logging

// This is for getting an idea of how to use the package and also tests simple concurrency

import (
	"fmt"
	"sync"
	"testing"
)

// makes sure you create a pointer to logger. this is because we need to pass it around.
var (
	semaphoreTest sync.WaitGroup
	logger        *Logger
)

func TestSimple(t *testing.T) {
	CreateLogFile(&logger) // serves as an initializer
	StartLog(logger)
	for i := 0; i <= 10; i++ {
		semaphoreTest.Add(1) // tells us that a routine is getting started
		go Hello(logger, i)
	}
	semaphoreTest.Wait() // makes sure all routines are finished before we close out
	StopLog(logger)
}

func Hello(l *Logger, num int) {
	defer semaphoreTest.Done() // signals that a routine is finished
	msg := fmt.Sprintf("Hello, I am in Routine %d", num)
	l.LogMessage(msg, "INFO")
}
