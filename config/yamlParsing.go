package config

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

// The Yaml type is meant to store the data from a configuration yaml.
// This struct needs to be updated any time there is new information meant to be brought in through the yaml.
type Yaml struct {
	WebScore int    `yaml:"webscore"`
	WebIP    string `yaml:"webIP"`
	WebDir   string `yaml:"webDir"`

	FtpScore int               `yaml:"ftpScore"`
	FtpIP    string            `yaml:"ftpIP"`
	FtpCreds map[string]string `yaml:"ftpCreds"`

	TeamScores map[string]int `yaml:"teamScores"`

	SshIP    string            `yaml:"sshIP"`
	PortNum  int               `yaml:"portNum"`
	SshCreds map[string]string `yaml:"sshCreds"`
}

// Parse uses the go-yaml library in order to take information out of a .yaml config file and place into a Yaml struct.
// This is accomplished by opening the .yaml file and then using yaml.Unmarshal in order to import the information from the yaml.
// Parse then returns the struct.
func Parse(path string) Yaml {

	var yamlPath = path

	file, err := os.ReadFile(yamlPath)
	if err != nil {
		fmt.Println("Failed to open the file: ", err)
	}

	var config Yaml

	if err := yaml.Unmarshal(file, &config); err != nil {
		fmt.Println("Failed to unmarshal the .yaml: ", err)
	}

	return config
}
