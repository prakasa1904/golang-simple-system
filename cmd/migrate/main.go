package main

import (
	"os"

	"github.com/devetek/golang-webapp-boilerplate/internal/config"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/group"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/member"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/order"
)

func main() {
	cfg := config.NewConfig()
	log := config.NewLogger(cfg)
	db := config.NewDatabase(config.DatabaseOption{
		Driver:   cfg.GetString("database.driver"),
		DBName:   cfg.GetString("database.name"),
		Username: cfg.GetString("database.username"),
		Password: cfg.GetString("database.password"),
	})

	// runtime env
	env := os.Getenv("ENV")

	// Define models
	var groupModel = &group.Entity{}
	var memberModel = &member.Entity{}
	var orderModel = &order.Entity{}

	if err := db.Migrator().AutoMigrate(groupModel, memberModel, orderModel); err != nil {
		log.Errorf("Migration error : %+v", err)
	}

	// Seeder for development
	if env == "development" {
		var email string = "admin@devetek.com"
		result := db.First(orderModel, "email = ?", email)
		if result.Error != nil {
			log.Errorf("Create seed development error : %+v", result.Error)
		}

		if result.RowsAffected != 0 {
			log.Warnf("Development seed data already exist")
		}

		if result.RowsAffected < 1 {
			var newMember = member.Entity{
				Username: "administrator",
				Email:    email,
			}

			result := db.Create(&newMember)
			if result.Error != nil {
				log.Errorf("Create database error : %+v", result.Error)
				panic(result.Error)
			}
		}
	}
}
