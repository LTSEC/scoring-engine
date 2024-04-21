package config

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

// The Yaml type is meant to store the data from a configuration yaml.
// This struct needs to be updated any time there is new information meant to be brought in through the yaml.
type Yaml struct {
	WebPortNum string `yaml:"webPortNum"`
	WebDir     string `yaml:"webDir"`
	Httpadd    int    `yaml:"httpAdd"`

	FtpPortNum string            `yaml:"ftpPortNum"`
	FtpCreds   map[string]string `yaml:"ftpCreds"`
	Ftpadd     int               `yaml:"ftpAdd"`

	SshPortNum string            `yaml:"sshPortNum"`
	SshCreds   map[string]string `yaml:"sshCreds"`
	Sshadd     int               `yaml:"sshAdd"`

	TeamScores map[string]int    `yaml:"teamScores"`
	TeamIpsFTP map[string]string `yaml:"teamIpsftp"`
	TeamIpsSSH map[string]string `yaml:"teamIpsssh"`
	TeamIpsWeb map[string]string `yaml:"teamIpsweb"`

	SleepTime int `yaml:"sleepTime"`
}

// Parse uses the go-yaml library in order to take information out of a .yaml config file and place into a Yaml struct.
// This is accomplished by opening the .yaml file and then using yaml.Unmarshal in order to import the information from the yaml.
// Parse then returns the struct.
func Parse(path string) *Yaml {

	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Failed to open the file: ", err)
	}

	var config Yaml

	if err := yaml.Unmarshal(file, &config); err != nil {
		fmt.Println("Failed to unmarshal the .yaml: ", err)
	}

	return &config
}
