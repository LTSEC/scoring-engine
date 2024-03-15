package scoring

import (
	"testing"
)

func TestScoring(t *testing.T) {
	t.Log(ScoringStartup([]string{"Team1", "Team2"}, 5, 1, 1))
}
