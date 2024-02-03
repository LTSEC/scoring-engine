package scoring

import (
	"io"
	"time"

	"github.com/jlaffaye/ftp"
)

// Connects to an IP address and returns a connection,
// if it could connect
func Connect(address string, username string, password string) (string, error) {
	connection, err := ftp.Dial(address, ftp.DialWithTimeout(5*time.Second))
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
