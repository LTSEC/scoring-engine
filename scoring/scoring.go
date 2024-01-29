package scoring

func Scoring() {
	dir := "a"
	site_ips := make([]string, 0)
	status, err := CheckWeb(dir, site_ips)
	if err != nil {
		// do stuff
	}
	println(status)
}
