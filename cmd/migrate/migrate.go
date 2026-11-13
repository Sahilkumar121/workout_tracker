package main

import (
	"log"
	"os"

	"github.com/Sahilkumar121/workout_tracker/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	// config
	cfg := config.MustLoad()

	// check the argument
	if len(os.Args) < 2 {
		log.Fatalln("Invalid argument required <up | down>")
	}

	m, err := migrate.New("file://migrations", cfg.DbUrl)
	if err != nil {
		log.Fatalf("Migration error:  %v\n", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatalf("Migration up failed: %v\n", err)
		}
		log.Println("Migration up successfull...")
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("Migration down failed: %v\n", err)
		}
		log.Println("Migration down successfull...")
	default:
		log.Fatalln("Unknown keyword...")
	}

}
