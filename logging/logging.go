package logging

import "fmt"

// for use in other programs

// for individual use in routines
func (l *Logger) LogMessage(msg string, status string) {
	// fmt.Println("LogMessage is NOT being called")
	if l != nil && l.status {
		// fmt.Println("LogMessage is being called")
		message := fmt.Sprintf(" %s: %s", status, msg)
		l.logger.Println(message)
	}
}

func (l *Logger) GetStatus() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.status
}

// takes a pointer to a pointer of a logger. creates the log file but doesn't immediately log
func CreateLogFile(l **Logger) {
	newLogger := new(Logger)
	newLogger.initialize()
	*l = newLogger
}

// called either to start logging or to resume logging (if PauseLog was called)
func StartLog(l *Logger) {
	if l != nil && !l.status {
		l.mutex.Lock()
		defer l.mutex.Unlock()
		l.setupLogger()
		l.status = true
		fmt.Println(fmt.Printf("%t", l.status))
	}
}

/* will come back to this later
func PauseLog(l *Logger) {
	if l != nil && l.status {
		l.mutex.Lock()
		defer l.mutex.Unlock()
		l.status = false
	}
}
*/

// should be called when you want to stop logging in the program as a whole- not when you're finished in one routine
func StopLog(l *Logger) {
	l.status = false
	if l != nil {
		l.cleanup()
	}
}
