package scoring

import (
	"testing"
)

func TestDBconnect(t *testing.T) {

	/*
		cfg_MySQL := Config{
			User:     "root",
			Password: "root",
			Host:     "localhost",
			Port:     3306,
			DBName:   "test",
			DBType:   "mysql",
			DBRef:    "test.yaml",
		}
	*/
	/*
		cfg_Postgres := Config{
			User:     "postgres",
			Password: "changeme",
			Host:     "localhost",
			Port:     5432,
			DBName:   "test",
			DBType:   "postgres",
			DBRef:    "test.yaml",
		}
	*/

	//CreateDatabase(cfg_Postgres)
	//ExecuteSQLFile(cfg_Postgres, "test.sql")
	//getTableNames_sqlite()
	//getTableNames_postgres()
	getTableNames_sqlite()
	/*	Test is divided into 3 different scenarios:
			a) User connection and Data being 1-1: points awarded
			b) Mismatch in schema (user connects): no points awarded
			c) Mismatch in data (user connects): no points awarded
		First tests for MySQL then tests for Postgres
	*/

	//fmt.Println("Testing MYSQL...")
	//RunTest(cfg_MySQL, "test.sql", "a")
	//RunTest(cfg_MySQL, "test2_schema.sql", "test_seed.sql", "b")
	//RunTest(cfg_MySQL, "test_schema.sql", "test3_seed.sql", "c")

	//fmt.Println("\nNow testing Postgres...")
	//RunTest(cfg_Postgres, "test.sql", "a")
	//RunTest(cfg_Postgres, "test2_schema.sql", "test_seed.sql", "b")
	//RunTest(cfg_Postgres, "test_schema.sql", "test3_seed.sql", "c")

	//conStr := createConStr_sqlite("test.sql")
	//getTableNames_sqlite()
	//fmt.Println(conStr)
	/*
		sliceExample := []string{"hello", "hi", "howdy"}
		for i := 0; i < len(sliceExample); i++ {
			fmt.Println(sliceExample[i])
		}
		fmt.Println("Went through")
	*/
}

/*
func RunTest(cfg Config, SQLFile string, testCase string) {
	var err error
	err = CreateDatabase(cfg)
	if err != nil {
		fmt.Println("issue with creating DB")
	}
	//defer DropDatabase(cfg)
	err = ExecuteSQLFile(cfg, SQLFile)
	if err != nil {
		fmt.Println("issue with file: %w", err)
	}
	//points, _, _ := ScoreDB(cfg.DBType, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.DBRef)
	//fmt.Printf("Test case %s: %d\n", testCase, points)
}
*/
