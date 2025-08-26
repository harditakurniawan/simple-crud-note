package databases

import (
	"log"
	entities "simple-crud-notes/databases/pgsql/entities"
)

func Migration() {
	err := DB.AutoMigrate(
		&entities.User{},
		&entities.AccessToken{},
		&entities.Note{},
	)

	if err != nil {
		log.Panicln("failed to migrate: " + err.Error())
	}
}
