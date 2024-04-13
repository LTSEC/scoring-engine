package logging

// handles setting up the struct and helper functions

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	fileName string
	logFile  *os.File
	logger   *log.Logger
	mutex    sync.Mutex
	status   bool
}

func (l *Logger) initialize() error {
	var err error
	filePath := getLogPath('\\')
	timeAndDate := time.Now().Format("2006-01-02 15-04-05")
	l.fileName = fmt.Sprintf("%s%s.txt", filePath, timeAndDate)
	_, err = os.Create(l.fileName)
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) setupLogger() error {
	var err error
	l.logFile, err = os.OpenFile(l.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// l.logFile, err = os.OpenFile("2024-04-12 22-03-32.txt", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error trying to open file")
		return err
	}
	l.logger = log.New(l.logFile, "", log.Ltime)
	return nil
}

// called whenever the Main code is finished as a cleanup
func (l *Logger) cleanup() error {
	var err error
	if l != nil && l.logFile != nil {
		if l.logFile != nil {
			err := l.logFile.Close()
			if err != nil {
				fmt.Println("Failed to close log file")
			}
		}
	}
	return err
}

// for creating a "Logs" folder in the directory. arguments are '\\' for windows and '/' for linux
func getLogPath(fileSeparator byte) string {
	var lastIndex int
	// gets current directory
	dirPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get working directory")
		return ""
	}
	// finds the index starting from right to left to find the delimiter's index
	for i := len(dirPath) - 1; i >= 0; i-- {
		if dirPath[i] == fileSeparator {
			lastIndex = i
			break
		}
	}
	// truncates everything until that index. this gives us the base path
	dirPath = dirPath[:lastIndex+1]
	// joins "Logs" with the base path which is where we will store our Log files
	newPath := fmt.Sprintf("%sLogs", dirPath)
	fmt.Println(newPath)
	// if the path doesn't exist, it creates one. may need to change the permissions later
	if _, err := os.Stat(newPath); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(newPath, 0755)
		}
	}
	// adds the delimiter ahead of the path name. this so that all we need to do is append a name to it later on
	// to create a file underneath this directory
	newPath = fmt.Sprintf("%s%c", newPath, fileSeparator)
	return newPath
}
