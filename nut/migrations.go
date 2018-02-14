package nut

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/migration"
	"github.com/astaxie/beego/orm"
)

// https://github.com/beego/bee/blob/master/cmd/commands/migrate/migrate.go

var (
	errBadDatabaseDriver = errors.New("not support database driver")
)

// Migration migration model
type Migration struct {
	migration.Migration
	up   string
	down string
}

// Up Run the migrations
func (p *Migration) Up() {
	p.SQL(p.up)
}

// Down Reverse the migrations
func (p *Migration) Down() {
	p.SQL(p.down)
}

func loadMigrations(driver string) error {
	return filepath.Walk(path.Join("db", driver, "migrations"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		name := info.Name()
		beego.Info("Find migration", name)
		up, err := ioutil.ReadFile(filepath.Join(path, "up.sql"))
		if err != nil {
			return nil
		}
		down, err := ioutil.ReadFile(filepath.Join(path, "down.sql"))
		if err != nil {
			return nil
		}
		mig := &Migration{up: string(up), down: string(down)}
		mig.Created = name[0:len(migration.DateFormat)]

		migration.Register(name, mig)
		return nil
	})
}

func showMigrationsTableSQL(driver orm.DriverType) (string, error) {
	switch driver {
	case orm.DRMySQL:
		return "SHOW TABLES LIKE 'migrations'", nil
	case orm.DRPostgres:
		return "SELECT * FROM pg_catalog.pg_tables WHERE tablename = 'migrations';", nil
	default:
		return "", errBadDatabaseDriver
	}
}

func selectMigrationsTableSQL(driver orm.DriverType) (string, error) {
	switch driver {
	case orm.DRMySQL:
		return "DESC migrations", nil
	case orm.DRPostgres:
		return "SELECT * FROM migrations WHERE false ORDER BY id_migration;", nil
	default:
		return "", errBadDatabaseDriver
	}
}

func createMigrationsTableSQL(driver orm.DriverType) (string, error) {
	switch driver {
	case orm.DRMySQL:
		return ``, nil
	case orm.DRPostgres:
		return `
		CREATE TYPE migrations_status AS ENUM('update', 'rollback');
		CREATE TABLE migrations (
			id_migration SERIAL PRIMARY KEY,
			name varchar(255) DEFAULT NULL,
			created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			statements text,
			rollback_statements text,
			status migrations_status
		);`, nil
	default:
		return "", errBadDatabaseDriver
	}
}

func getLatestMigration(db *sql.DB) (string, int64, error) {
	sql := "SELECT name FROM migrations where status = 'update' ORDER BY id_migration DESC LIMIT 1"
	rows, err := db.Query(sql)
	if err != nil {
		return "", 0, err
	}
	var name string
	var created int64
	if rows.Next() {
		if err = rows.Scan(&name); err != nil {
			return "", 0, err
		}

		t, err := time.Parse(migration.DateFormat, name[0:len(migration.DateFormat)])
		if err != nil {
			return "", 0, err
		}
		created = t.Unix()
	}

	return name, created, nil
}

func checkForSchemaUpdateTable(db *sql.DB, driver orm.DriverType) error {
	showTableSQL, err := showMigrationsTableSQL(driver)
	if err != nil {
		return err
	}

	rows, err := db.Query(showTableSQL)
	if err != nil {
		return err
	}

	if !rows.Next() {
		createTableSQL, err := createMigrationsTableSQL(driver)
		if err != nil {
			return err
		}
		beego.Notice("Creating 'migrations' table...")
		if _, err := db.Query(createTableSQL); err != nil {
			return err
		}
	}
	return nil
}
