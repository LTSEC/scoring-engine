package scoring

import (
	"math/rand"
	"sync"
	"time"

	"github.com/LTSEC/scoring-engine/config"
	"github.com/LTSEC/scoring-engine/logging"
	"github.com/LTSEC/scoring-engine/score_holder"
)

var ScoringOn bool
var mutex sync.Mutex

// Starts the scoring process
func ScoringStartup(yamlConfig *config.Yaml) error {
	// get list of teams from existing maps
	i := 0
	TeamNames := make([]string, len(yamlConfig.TeamScores))
	for k := range yamlConfig.TeamScores {
		TeamNames[i] = k
		i++
	}

	score_holder.Startup(TeamNames)
	logger := new(logging.Logger)
	logger.StartLog()
	ScoringOn = true
	// todo: make scoring on pause scoring instead of just stopping it entirely, so that it may be resumed later
	for ScoringOn {
		for index, teamName := range TeamNames {
			go scoreTeam(index, teamName, yamlConfig, logger)
		}

		time.Sleep(time.Duration(yamlConfig.SleepTime) * time.Second)
	}

	return nil
}

// Ends the scoring process
func ScoringToggle(state bool) error {
	//logging.LogMessage("info", "Stopping scoring")
	ScoringOn = state
	return nil
}

func b2i(b bool) int {
	// The compiler currently only optimizes this form.
	// See issue 6011.
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}

// Scores an individual team
func scoreTeam(index int, teamName string, yamlConfig *config.Yaml, logger *logging.Logger) {
	mutex.Lock()
	defer mutex.Unlock()
	var FTPUser string
	var FTPPass string

	var SSHUser string
	var SSHPass string

	// Choose a random FTP user
	k := rand.Intn(len(yamlConfig.FtpCreds))

	for User, Pass := range yamlConfig.FtpCreds {
		if k == 0 {
			FTPUser = User
			FTPPass = Pass
			break
		}
		k--
	}
	// Choose a random SSH user
	k = rand.Intn(len(yamlConfig.FtpCreds))

	for User, Pass := range yamlConfig.SshCreds {
		if k == 0 {
			SSHUser = User
			SSHPass = Pass
			break
		}
		k--
	}

	ftp, err := FTPConnect(yamlConfig.TeamIpsFTP[teamName], yamlConfig.FtpPortNum, FTPUser, FTPPass)
	if err != nil {
		logger.LogMessage(err.Error(), "error")
	}

	ssh, err := SSHConnect(yamlConfig.TeamIpsSSH[teamName], yamlConfig.SshPortNum, SSHUser, SSHPass)
	if err != nil {
		logger.LogMessage(err.Error(), "error")
	}

	http, err := CheckWeb(yamlConfig.WebDir, yamlConfig.TeamIpsWeb[teamName], yamlConfig.WebPortNum)
	if err != nil {
		logger.LogMessage(err.Error(), "error")
	}

	score_holder.UpdateTeam(index, score_holder.NewScoreMap(yamlConfig.Ftpadd*b2i(ftp != ""),
		yamlConfig.Sshadd*b2i(ssh), yamlConfig.Httpadd*b2i(ssh)),
		score_holder.NewStateMap(ftp != "", ssh, http))
}
