package score_holder

type Map map[string]interface{}
type ScoreMap = []map[string]map[string]interface{}

var AllTeams = make(ScoreMap, 0)

// Starts the score holder, takes an ORDERED LIST OF TEAMS as an input and then stores it in a map
func Startup(teams []string) {
	TeamMapLength := len(teams)
	for i := 0; i < TeamMapLength; i++ {

		NewTeam := make(map[string]map[string]interface{})
		NewTeam["scores"] = Map{"ftp": 0, "http": 0}
		NewTeam["states"] = Map{"ftp": false, "http": false}
		AllTeams = append(AllTeams, NewTeam)
	}

}

// Returns the scoremap
func GetMap() ScoreMap {
	return AllTeams
}

func UpdateMap(TeamIndex int, scores Map, states Map) {
	if scores != nil {
		AllTeams[TeamIndex]["scores"] = scores
	}
	if states != nil {
		AllTeams[TeamIndex]["states"] = states
	}
}
