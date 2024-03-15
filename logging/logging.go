package logging

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	file        *os.File
	mutex       sync.Mutex
	initialized bool

	// Loggers
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

// CreateLogFile prepares the logging environment by creating a log file based on the current date and time.
func CreateLogFile() {
	now := time.Now()
	fileName := now.Format("2006-01-02_15-04-05") + " Log.txt"
	var err error
	file, err = os.Create(fileName)
	if err != nil {
		log.Printf("Failed to create log file: %v", err)
		return
	}
	SetLogFile(file) // Pass the file descriptor directly.
}

// SetLogFile configures the logging mechanism to output to a specified file.
func SetLogFile(f *os.File) {
	mutex.Lock()
	defer mutex.Unlock()

	if file != nil {
		file.Close() // Ensure the previously opened file is closed before reassignment.
	}

	file = f // Assign the new file descriptor.

	// Initialize loggers with the new file.
	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	initialized = true // Mark the logging system as initialized.
}

// LogMessage allows logging a message with a specified severity.
func LogMessage(severity, message string) {
	if !initialized {
		fmt.Println("Logging system not initialized")
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	switch severity {
	case "info":
		infoLogger.Println(message)
	case "warning":
		warningLogger.Println(message)
	case "error":
		errorLogger.Println(message)
	default:
		fmt.Println("Invalid severity specified:", severity)
	}
}

// StopLog should be called to close the log file before the application exits.
func StopLog() {
	if file != nil {
		file.Close()
	}
}
