package scoring

import (
	"fmt"
	"sync"
	"time"

	"github.com/LTSEC/scoring-engine/logging"
	"github.com/LTSEC/scoring-engine/score_holder"
)

var ScoringOn bool
var mutex sync.Mutex

// Starts the scoring process
func ScoringStartup(teams []string, ftpadd int, sshadd int, httpadd int, sleeptime int, ip string, ipstartnum int) error {
	score_holder.Startup(teams)
	logging.CreateLogFile()
	ScoringOn = true
	// todo: make scoring on pause scoring instead of just stopping it entirely, so that it may be resumed later
	for ScoringOn {
		for index, _ := range teams {
			go scoreTeam(index, ip, ipstartnum, ftpadd, sshadd, httpadd)
			fmt.Println(score_holder.GetMap()[index])
		}

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

// Scores an individual team
func scoreTeam(index int, ip string, ipstartnum int, ftpadd int, sshadd int, httpadd int) {
	mutex.Lock()
	defer mutex.Unlock()
	var scoreaddftp int
	var scoreaddssh int
	var scoreaddhttp int
	conn := ip + fmt.Sprint(ipstartnum+index)

	ftp, err := FTPConnect(conn, "ftpuser", "ftppass")
	if err != nil {
		//logging.LogMessage("info", err.Error())
	}

	ssh, err := SSHConnect(conn, "22", "testuser", "testpass")
	if err != nil {
		//logging.LogMessage("info", err.Error())
	}

	http, err := CheckWeb("C:/Users/Aidan Feess/Documents/Projects/LTSEC/scoring-engine/tests/site_infos/site_info.html", conn)
	if err != nil {
		//logging.LogMessage("info", err.Error())
	}

	if ftp != "" {
		scoreaddftp += ftpadd
	}

	if ssh {
		scoreaddssh += sshadd
	}

	if http {
		scoreaddhttp += httpadd
	}

	score_holder.UpdateTeam(index, score_holder.NewScoreMap(scoreaddftp, scoreaddhttp, scoreaddssh),
		score_holder.NewStateMap(ftp != "", ssh, http))
}
