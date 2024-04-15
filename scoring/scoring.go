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
	logging.CreateLogFile()
	ScoringOn = true
	// todo: make scoring on pause scoring instead of just stopping it entirely, so that it may be resumed later
	for ScoringOn {
		for index, teamName := range TeamNames {
			go scoreTeam(index, teamName, yamlConfig)
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

// Scores an individual team
func scoreTeam(index int, teamName string, yamlConfig *config.Yaml) {
	mutex.Lock()
	defer mutex.Unlock()
	var scoreaddftp int
	var scoreaddssh int
	var scoreaddhttp int

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
		logging.LogMessage("error", err.Error())
	}

	ssh, err := SSHConnect(yamlConfig.TeamIpsSSH[teamName], yamlConfig.FtpPortNum, SSHUser, SSHPass)
	if err != nil {
		logging.LogMessage("error", err.Error())
	}

	http, err := CheckWeb(yamlConfig.WebDir, yamlConfig.TeamIpsWeb[teamName], yamlConfig.WebPortNum)
	if err != nil {
		logging.LogMessage("error", err.Error())
	}

	if ftp != "" {
		scoreaddftp += yamlConfig.Ftpadd
	}

	if ssh {
		scoreaddssh += yamlConfig.Sshadd
	}

	if http {
		scoreaddhttp += yamlConfig.Httpadd
	}

	score_holder.UpdateTeam(index, score_holder.NewScoreMap(scoreaddftp, scoreaddhttp, scoreaddssh),
		score_holder.NewStateMap(ftp != "", ssh, http))
}
