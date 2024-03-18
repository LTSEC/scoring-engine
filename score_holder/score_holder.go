package score_holder

type Map map[string]interface{}
type ScoreMap = []map[string]map[string]interface{}

var AllTeams = make(ScoreMap, 0)

// Starts the score holder, takes an ORDERED LIST OF TEAMS as an input and then stores it in a map
func Startup(teams []string) error {
	TeamMapLength := len(teams)
	for i := 0; i < TeamMapLength; i++ {
		NewTeam := make(map[string]map[string]interface{})
		NewTeam["scores"] = Map{"ftp": 0, "http": 0, "ssh": 0}
		NewTeam["states"] = Map{"ftp": false, "http": false, "ssh": false}
		AllTeams = append(AllTeams, NewTeam)
	}
	return nil
}

// Returns the scoremap
func GetMap() ScoreMap {
	return AllTeams
}

// Creates a new score map and returns it
func NewScoreMap(ftp int, http int, ssh int) Map {
	return Map{"ftp": ftp, "http": http, "ssh": ssh}
}

// Creates a new state map and returns it
func NewStateMap(ftp bool, http bool, ssh bool) Map {
	return Map{"ftp": ftp, "http": http, "ssh": ssh}
}

// Takes a team, map of scores, and map of states as input to update a team's value
func UpdateTeam(TeamIndex int, scores Map, states Map) error {
	if scores != nil {
		for scoretype, score := range AllTeams[TeamIndex]["scores"] {
			add_score := scores[scoretype].(int)
			old_score := score.(int)
			AllTeams[TeamIndex]["scores"][scoretype] = add_score + old_score
		}
	}
	if states != nil {
		AllTeams[TeamIndex]["states"] = states
	}
	return nil
}

func SetTeam(TeamIndex int, scores Map, states Map) error {
	if scores != nil {
		AllTeams[TeamIndex]["scores"] = scores
	}
	if states != nil {
		AllTeams[TeamIndex]["states"] = states
	}
	return nil
}
