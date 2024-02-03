package scoring

func Scoring() []bool {
	dir := "a"
	site_ips := make([]string, 0)
	status, err := CheckWeb(dir, site_ips)
	if err != nil {
		// do stuff
	}
	return status
}
