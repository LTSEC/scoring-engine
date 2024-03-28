// These methods should be used ONLY in the main process

package logging

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	initialized    bool
	logFile        *os.File
	logger         *log.Logger
	logCh          chan string
	doneCh         chan bool
	semaphore      sync.WaitGroup
	activeMessages sync.WaitGroup
	/* ctx    context.Context
	cancel context.CancelFunc */
)

// USE: Make sure this function is called BEFORE any other functions are called in your program
// WHAT: Creates a txt file, channels to communicate with other routines, and the logger struct
func CreateLogFile() {
	if initialized {
		return
	}
	initialized = true
	filePath := "C:/Users/bobby/OneDrive/Desktop/testing conc/Logs" // <- Change this later
	var err error
	filename := filePath + "/Log " + time.Now().Format("2006-01-02 15-04-05") + ".txt"
	logFile, err = os.Create(filename)
	if err != nil {
		fmt.Println("Failed to open log file")
		return
	}
	logCh = make(chan string, 100) // Buffered channel for storing log messages
	doneCh = make(chan bool)       // Channel to signal done (to tell it to stop logging)
	logger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime)
}

// USE: Has to always follow CreateLogFile.
// HOW: Spawns another routine which listens to and logs anything on the message chanel
func StartLog() {
	if !initialized {
		return
	}
	semaphore.Add(1)
	go func() {
		defer semaphore.Done()
		for {
			select {
			case msg := <-logCh:
				if logFile != nil {
					logger.Println(msg) // Write directly to the file
					activeMessages.Done()
				}
			case <-doneCh:
				return
			}
		}
	}()
}

// USE: When you're done logging (should be called after all routines are finished)
// HOW: Closes doneCh which signals StartLog to finish and then closes the file
func StopLog() {
	if !initialized {
		return
	}
	activeMessages.Wait()
	close(doneCh)    // Signal to StartLog to stop taking any more information
	semaphore.Wait() // Wait for the logging goroutine to finish
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			fmt.Println("Failed to close log file")
		}
	}
}
