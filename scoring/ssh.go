package scoring

import (
	"golang.org/x/crypto/ssh"
)

func SSHConnect(hostname string, port string, username string, password string) (bool, error) {

	host := hostname + ":" + port

	//var key ssh.Signer
	var err error

	//key, err = ssh.ParsePrivateKey(pKey)
	//if err != nil {
	//	return false, err
	//}

	//var hostkeyCallback ssh.HostKeyCallback
	//hostkeyCallback, err = knownhosts.New("C:/Users/Aidan Feess/.ssh/known_hosts")
	//if err != nil {
	//	return false, err
	//}

	conf := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			//ssh.PublicKeys(key),
			ssh.Password(password),
		},
	}

	var conn *ssh.Client

	conn, err = ssh.Dial("tcp", host, conf)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	var session *ssh.Session
	session, err = conn.NewSession()
	if err != nil {
		return false, err
	}
	defer session.Close()

	return true, nil
}
