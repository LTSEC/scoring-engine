// For logging capabilities in your own routines. Use it by importing this package into your program
package logging

// include this each time you need to log something. Adds a message to the channel so it can be listened
func LogMessage(message string) {
	activeMessages.Add(1)
	logCh <- message
}

// maybe later add a pause and etc
