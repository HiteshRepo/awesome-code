package migrations

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	dbCnf "github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/pkg/configs/db"
	"gorm.io/gorm"
	"log"
	"path/filepath"
	"strings"
)

const (
	cutSet       = "file://"
	databaseName = "postgres"
)

type Migrator struct {
	pgDBMigrate        *migrate.Migrate
}

func ProvideMigrator(config dbCnf.DatabaseConfig, pgDB *gorm.DB) (*Migrator, error) {
	dbConn, err := pgDB.DB()
	if err != nil {
		return nil, err
	}

	pgDBMigrate, err := initMigrate(dbConn, config.MigrationPath)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		pgDBMigrate: pgDBMigrate,
	}, nil
}

func (m Migrator) RunMigrations() {
	m.RunMigrationsWith(m.pgDBMigrate, "Postgres Database")
}

func (m Migrator) RunMigrationsWith(migrateInstance *migrate.Migrate, dBName string) {
	if err := migrateInstance.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println(fmt.Sprintf("No change detected after running the migrations for %s", dBName))
			return
		}
		log.Println(fmt.Sprintf("Migration Failed for %s", dBName), err)
	}
	log.Println(fmt.Sprintf("Migrations applied successfully to %s", dBName))
}

func initMigrate(dbConn *sql.DB, directory string) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	sourcePath, err := getSourcePath(directory)
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(sourcePath, databaseName, driver)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func getSourcePath(directory string) (string, error) {
	directory = strings.TrimPrefix(directory, cutSet)

	absPath, err := filepath.Abs(directory)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", cutSet, absPath), nil
}
