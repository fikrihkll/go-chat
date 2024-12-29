package common

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fikrihkll/chat-app/config"
	"github.com/fikrihkll/chat-app/infrastructure"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// migration migrate to database
func Migration(wd string) error {
	cfg := config.Load(fmt.Sprintf("%s/.env", wd))
	pgConn, err := infrastructure.NewPgConnection(cfg)
	driver, err := postgres.WithInstance(pgConn, &postgres.Config{})
	if err != nil {
		return errors.New("CANNOT CONNECT TO DATABASE")
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/migrations", wd),
		"postgres",
		driver,
	)

	if len(os.Args) > 2 {
		if os.Args[2] == "down" {
			if err := m.Down(); err != nil {
				log.Println("migration down error")
				return err
			}
			log.Println("migration down success")

			return nil
		} else if os.Args[2] == "step-down" {
			if err := m.Steps(-1); err != nil {
				log.Println("migration step-down to the 1 last migration error")
				return err
			}
			log.Println("migration step-down to the 1 last migration success")

			return nil
		} else if os.Args[2] == "info" {
			version, dirty, err := m.Version()
			if err != nil {
				log.Println("migration info error")
				return err
			}

			log.Printf("Version: %d\nDirty: %t", version, dirty)

			return nil
		}  else if os.Args[2] == "force" {
			version, err := strconv.Atoi(os.Args[3])
			if err != nil {
				log.Println("error force")
				return err
			}

			err = m.Force(version)
			if err != nil {
				log.Println("migration info error")
				return err
			}

			log.Printf("Force to version %d successfully executed", version)

			return nil
		}
	}

	if err := m.Up(); err != nil {
		log.Println("migration up error")
		return err
	}
	log.Println("migration up success")

	return nil
}
