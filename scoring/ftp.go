package scoring

import (
	"io"
	"time"

	"github.com/jlaffaye/ftp"
)

// Connects to an IP address and returns a connection,
// if it could connect
func Connect(address string) error {
	connection, err := ftp.Dial(address, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}

	err = connection.Login("anonymous", "anonymous")
	if err != nil {
		return err
	}

	result, err := connection.Retr("test-file.txt")
	if err != nil {
		return err
	}
	defer result.Close()

	buf, err := io.ReadAll(result)
	println(string(buf))

	err = connection.Quit()

	if err != nil {
		return err
	}

	return nil
}
