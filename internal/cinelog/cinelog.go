package cinelog

import "github.com/Cinelog-tv/cinelog/internal/database"

func Run() {

	databaseConfig := database.DatabaseConfig{
		Driver: database.DriverSQLite,
		DSN:    "cinelog.db",
	}

	_, err := database.NewDatabase(databaseConfig)

	if err != nil {
		panic(err)
	}
}
