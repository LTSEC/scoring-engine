package scoring

import (
	"testing"
)

func TestSSH(t *testing.T) {
	//pKeyFile := "C:/Users/Aidan Feess/.ssh/id_rsa"

	//pKey, err := os.ReadFile(pKeyFile)
	//if err != nil {
	//	log.Fatal(err)
	//}

	t.Log(SSHConnect("172.29.1.5", "22", "testuser", "testpass"))
}
