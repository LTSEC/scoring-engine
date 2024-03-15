package scoring

import (
	"fmt"
	"time"

	"github.com/LTSEC/scoring-engine/score_holder"
)

func ScoringStartup(teams []string, ftpadd int, sshadd int, httpadd int) error {
	score_holder.Startup(teams)
	for i := 1; i < 3; i++ {
		for index, _ := range teams {

			var scoreaddftp int
			var scoreaddssh int
			var scoreaddhttp int

			ftp, err := FTPConnect("172.29.1.5", "ftpuser", "ftppass")
			if err != nil {
				fmt.Println(err)
			}

			ssh, err := SSHConnect("172.29.1.5", "22", "testuser", "testpass")
			if err != nil {
				fmt.Println(err)
			}

			http, err := CheckWeb("C:/Users/Aidan Feess/Documents/Projects/LTSEC/scoring-engine/tests/site_infos/site_info.html", "172.29.1.5")
			if err != nil {
				fmt.Println(err)
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
		time.Sleep(2 * time.Second)
	}
	return nil
}
