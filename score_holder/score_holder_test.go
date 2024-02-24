package score_holder

import (
	"testing"
)

func TestSaveData(t *testing.T) {
	teams := []string{"TeamA", "TeamB", "TeamC"}
	t.Log(Startup(teams))
	t.Log(GetMap())
	t.Log(UpdateTeam(1, NewScoreMap(0, 0), NewStateMap(false, false)))
}
