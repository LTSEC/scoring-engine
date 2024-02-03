package scoring

import "testing"

func TestFTP(t *testing.T) {
	t.Log(Connect("172.31.255.5:21", "ftpuser", "password"))
}
