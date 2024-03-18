package scoring

import (
	"testing"
	"time"
)

func TestScoring(t *testing.T) {
	go func() {
		t.Log(ScoringStartup([]string{"Team1", "Team2", "Team3", "Team4"}, 5, 1, 1, 5, "172.29.1.", 5))
	}()
	time.Sleep(11 * time.Second)
	t.Log(ScoringToggle(false))
}
