package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Sahilkumar121/workout_tracker/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	cfg := config.MustLoad()

	// check if argv is valid
	if len(os.Args) < 2 {
		log.Fatalln("migrate usage: <up | down>")
	}

	m, err := migrate.New("file://migrations", cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("failed to migrate %v \n", err)
	}
	switch os.Args[1] {
	case "up":
		m.Up()
	case "down":
		m.Steps(-1)
	default:
		log.Fatalln("unknown keyword")
	}

	fmt.Println("miragtion successful ...")
}
