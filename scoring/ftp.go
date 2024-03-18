package scoring

import (
	"io"
	"time"

	"github.com/jlaffaye/ftp"
)

// Attempts a connection via FTP and returns a boolean value representing success
func FTPConnect(address string, username string, password string) (string, error) {
	connection, err := ftp.Dial(address+":21", ftp.DialWithTimeout(250*time.Millisecond))
	if err != nil {
		return "", err
	}

	err = connection.Login(username, password)
	if err != nil {
		return "", err
	}

	result, err := connection.Retr("textfile.txt")
	if err != nil {
		return "", err
	}
	defer result.Close()

	buf, err := io.ReadAll(result)

	if err != nil {
		return "", err
	}

	err = connection.Quit()

	if err != nil {
		return "", err
	}

	return string(buf), nil
}
