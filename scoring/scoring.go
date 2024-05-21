package scoring

import (
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"time"

	"github.com/LTSEC/scoring-engine/config"
	"github.com/LTSEC/scoring-engine/logging"
	"github.com/LTSEC/scoring-engine/score_holder"
)

var (
	ScoringOn bool
	mutex     sync.Mutex
	logger    *logging.Logger
	TeamNames []string
)

// Starts the scoring process
func ScoringStartup(yamlConfig *config.Yaml) error {
	// get list of teams from existing maps
	i := 0
	TeamNames = make([]string, len(yamlConfig.TeamScores))
	for k := range yamlConfig.TeamScores {
		TeamNames[i] = k
		i++
	}
	slices.Sort(TeamNames)

	score_holder.Startup(TeamNames)
	logger = new(logging.Logger)
	logger.StartLog()
	ScoringOn = true
	// todo: make scoring on pause scoring instead of just stopping it entirely, so that it may be resumed later
	activelyScore(yamlConfig)

	return nil
}

// Ends the scoring process
func ScoringToggle(state bool, yamlConfig *config.Yaml) error {
	logger.LogMessage(fmt.Sprint("Engine online: ", state), "Info")
	ScoringOn = state
	if state == true {
		activelyScore(yamlConfig)
	}
	return nil
}

// Restarts the scoring process
func activelyScore(yamlConfig *config.Yaml) {
	for ScoringOn {
		for index, teamName := range TeamNames {
			go scoreTeam(index, teamName, yamlConfig)
		}

		time.Sleep(time.Duration(yamlConfig.SleepTime) * time.Second)
	}
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
func scoreTeam(index int, teamName string, yamlConfig *config.Yaml) {
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
		logger.LogMessage(err.Error(), "FTP Error")
	}

	ssh, err := SSHConnect(yamlConfig.TeamIpsSSH[teamName], yamlConfig.SshPortNum, SSHUser, SSHPass)
	if err != nil {
		logger.LogMessage(err.Error(), "SSH Error")
	}

	http, err := CheckWeb(yamlConfig.WebDir, yamlConfig.TeamIpsWeb[teamName], yamlConfig.WebPortNum)
	if err != nil {
		logger.LogMessage(err.Error(), "HTTP Error")
	}

	score_holder.UpdateTeam(index, score_holder.NewScoreMap(yamlConfig.Ftpadd*b2i(ftp != ""),
		yamlConfig.Httpadd*b2i(http), yamlConfig.Sshadd*b2i(ssh)),
		score_holder.NewStateMap(ftp != "", http, ssh))
}
