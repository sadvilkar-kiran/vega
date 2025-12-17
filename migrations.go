package vega

import (
	"log"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateUp runs all pending migrations
func (v *Vega) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+v.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		log.Println("Error running migration:", err)
		return err
	}
	return nil
}

// MigrateDownAll rolls back all migrations
func (v *Vega) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+v.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil {
		return err
	}

	return nil
}

// Steps runs n migration steps (positive for up, negative for down)
func (v *Vega) Steps(n int, dsn string) error {
	m, err := migrate.New("file://"+v.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(n); err != nil {
		return err
	}

	return nil
}

// MigrateForce forces the migration version
func (v *Vega) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+v.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Force(-1); err != nil {
		return err
	}

	return nil
}

