package main

import (
	"os"

	"github.com/devetek/golang-webapp-boilerplate/internal/config"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/member"
)

func main() {
	cfg := config.NewConfig()
	log := config.NewLogger(cfg)
	db := config.New(config.DatabaseOption{
		Driver:   cfg.GetString("database.driver"),
		DBName:   cfg.GetString("database.name"),
		Username: cfg.GetString("database.username"),
		Password: cfg.GetString("database.password"),
	})

	// runtime env
	env := os.Getenv("ENV")

	// Define models
	var user = &member.Entity{}

	if err := db.Migrator().AutoMigrate(user); err != nil {
		log.Errorf("Migration error : %+v", err)
	}

	// Seeder for development
	if env == "development" {
		var email string = "admin@devetek.com"
		result := db.First(user, "email = ?", email)
		if result.Error != nil {
			log.Errorf("Create seed development error : %+v", result.Error)
		}

		if result.RowsAffected != 0 {
			log.Warnf("Development seed data already exist")
		}

		if result.RowsAffected < 1 {
			var newUser = member.Entity{
				Username: "administrator",
				Email:    email,
			}

			result := db.Create(&newUser)
			if result.Error != nil {
				log.Errorf("Create database error : %+v", result.Error)
				panic(result.Error)
			}
		}
	}
}
