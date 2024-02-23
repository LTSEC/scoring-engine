package score_holder

type Map map[string]interface{}

var AllTeams = make([]map[string]map[string]interface{}, 0)

func Startup(teams []string) {
	TeamMapLength := len(teams)
	for i := 0; i < TeamMapLength; i++ {

		NewTeam := make(map[string]map[string]interface{})
		NewTeam["scores"] = Map{"ftp": 0, "http": 0}
		NewTeam["states"] = Map{"ftp": false, "http": false}
		AllTeams = append(AllTeams, NewTeam)
	}

}

func GetMap() {

}

func UpdateMap() {

}
