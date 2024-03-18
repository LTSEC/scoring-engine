package scoring

import (
	"fmt"
	"time"

	"github.com/LTSEC/scoring-engine/logging"
	"github.com/LTSEC/scoring-engine/score_holder"
)

var ScoringOn bool

// Starts the scoring process
func ScoringStartup(teams []string, ftpadd int, sshadd int, httpadd int, sleeptime int) error {
	score_holder.Startup(teams)
	logging.CreateLogFile()
	ScoringOn = true
	// todo: make scoring on pause scoring instead of just stopping it entirely, so that it may be resumed later
	for ScoringOn {
		for index, _ := range teams {

			var scoreaddftp int
			var scoreaddssh int
			var scoreaddhttp int

			ftp, err := FTPConnect("172.29.1.5", "ftpuser", "ftppass")
			if err != nil {
				//logging.LogMessage("info", err.Error())
			}

			ssh, err := SSHConnect("172.29.1.5", "22", "testuser", "testpass")
			if err != nil {
				//logging.LogMessage("info", err.Error())
			}

			http, err := CheckWeb("C:/Users/Aidan Feess/Documents/Projects/LTSEC/scoring-engine/tests/site_infos/site_info.html", "172.29.1.5")
			if err != nil {
				//logging.LogMessage("info", err.Error())
			}

			if ftp != "" {
				scoreaddftp += ftpadd
			}

			if ssh {
				scoreaddssh += ftpadd
			}

			if http {
				scoreaddhttp += httpadd
			}

			score_holder.UpdateTeam(index, score_holder.NewScoreMap(scoreaddftp, scoreaddhttp, scoreaddssh),
				score_holder.NewStateMap(ftp != "", ssh, http))
		}
		fmt.Println(score_holder.GetMap()[1]["scores"])
		time.Sleep(time.Duration(sleeptime) * time.Second)
	}

	return nil
}

// Ends the scoring process
func ScoringToggle(state bool) error {
	//logging.LogMessage("info", "Stopping scoring")
	ScoringOn = state
	return nil
}
