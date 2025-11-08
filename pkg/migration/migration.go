package migration

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
)

const migrationsTableName = "migrations"

type DomainMigration interface {
	Migrations() []Migration
}

type Migration interface {
	TableName() string
	Description() string
	CreatedAt() string
	Name() string
	Up(*sql.DB) error
	Down(*sql.DB) error
}

type BaseMigration struct {
	TableNameValue   string
	CreatedAtValue   string
	DescriptionValue string
}

func (b BaseMigration) TableName() string {
	return b.TableNameValue
}

func (b BaseMigration) Description() string {
	return b.DescriptionValue
}

func (b BaseMigration) CreatedAt() string {
	return b.CreatedAtValue
}

func (b BaseMigration) Name() string {
	createdAt := strings.ReplaceAll(strings.ToLower(b.CreatedAt()), " ", "_")
	createdAt = strings.ReplaceAll(createdAt, "-", "")
	createdAt = strings.ReplaceAll(createdAt, ":", "")
	description := strings.ReplaceAll(strings.ToLower(b.Description()), " ", "_")
	return fmt.Sprintf("%s_%s", createdAt, description)
}

type Migrator struct {
	migrations map[string]Migration
	dbConn     *sql.DB
}

func NewMigrator(dbConn *sql.DB) *Migrator {
	return &Migrator{
		migrations: make(map[string]Migration, 0),
		dbConn:     dbConn,
	}
}

func (m *Migrator) Register(migrations ...DomainMigration) {
	for _, domainMigration := range migrations {
		for _, migration := range domainMigration.Migrations() {
			m.migrations[migration.Name()] = migration
		}
	}
}

func (m *Migrator) Migrate() error {
	m.createMigrationsTableIfNotCreated()
	batchNumber := m.getLastBatchNumber()
	batchNumber++
	migrations := m.getMigrationsSorted()
	var numOfMigrations int
	for _, migration := range migrations {
		if m.migrationExists(migration) {
			continue
		}
		if err := migration.Up(m.dbConn); err != nil {
			log.Fatalf("error migrating %s: %v\n", migration.Name(), err)
			return err
		}
		m.recordMigration(migration, batchNumber)
		log.Printf("migrated %s\n", migration.Name())
		numOfMigrations++
	}

	log.Printf("total migrations: %d\n", numOfMigrations)
	return nil
}

func (m *Migrator) Rollback() error {
	m.createMigrationsTableIfNotCreated()
	batchNumber := m.getLastBatchNumber()
	migrations := m.getBatchMigrations(batchNumber)
	for _, migration := range migrations {
		if err := migration.Down(m.dbConn); err != nil {
			log.Fatalf("error rolling back %s: %v\n", migration.Name(), err)
			return err
		}
		m.removeMigration(migration, batchNumber)
		log.Printf("rolled back %s\n", migration.Name())
	}
	fmt.Printf("rolled back %d migrations\n", len(migrations))
	return nil
}

func (m *Migrator) createMigrationsTableIfNotCreated() {
	query := `CREATE TABLE IF NOT EXISTS %s (
    id BIGSERIAL PRIMARY KEY,
    migration VARCHAR(255) NOT NULL,
    batch INT NOT NULL
)`
	query = fmt.Sprintf(query, migrationsTableName)
	_, err := m.dbConn.Exec(query)
	if err != nil {
		log.Fatalf("error creating migrations table: %v\n", err)
	}
}

func (m *Migrator) getMigrationsSorted() []Migration {
	migrations := make([]Migration, 0)
	for _, migration := range m.migrations {
		migrations = append(migrations, migration)
	}
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].CreatedAt() < migrations[j].CreatedAt()
	})

	return migrations
}

func (m *Migrator) getLastBatchNumber() int {
	query := `SELECT batch FROM %s ORDER BY batch DESC LIMIT 1`
	query = fmt.Sprintf(query, migrationsTableName)
	var batchNumber int
	err := m.dbConn.QueryRow(query).Scan(&batchNumber)
	if err != nil {
		return 0
	}

	return batchNumber
}

func (m *Migrator) migrationExists(migration Migration) bool {
	query := `SELECT id FROM %s WHERE migration = $1 LIMIT 1`
	query = fmt.Sprintf(query, migrationsTableName)
	var id int
	err := m.dbConn.QueryRow(query, migration.Name()).Scan(&id)
	if err != nil {
		return false
	}

	return true
}

func (m *Migrator) recordMigration(migration Migration, number int) {
	query := `INSERT INTO %s (migration, batch) VALUES ($1, $2)`
	query = fmt.Sprintf(query, migrationsTableName)
	_, err := m.dbConn.Exec(query, migration.Name(), number)
	if err != nil {
		log.Fatalf("error recording migration: %v\n", err)
	}
}

func (m *Migrator) getBatchMigrations(number int) []Migration {
	query := `SELECT migration FROM %s WHERE batch = $1`
	query = fmt.Sprintf(query, migrationsTableName)
	rows, err := m.dbConn.Query(query, number)
	if err != nil {
		log.Fatalf("error getting migrations: %v\n", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatalf("error closing migrations rows: %v\n", err)
		}
	}()

	migrations := make([]Migration, 0)
	for rows.Next() {
		var migrationName string
		if err := rows.Scan(&migrationName); err != nil {
			log.Fatalf("error getting migrations: %v\n", err)
		}
		migration, exists := m.migrations[migrationName]
		if !exists {
			continue
		}
		migrations = append(migrations, migration)
	}

	return migrations
}

func (m *Migrator) removeMigration(migration Migration, number int) {
	query := `DELETE FROM %s WHERE migration = $1 AND batch = $2`
	query = fmt.Sprintf(query, migrationsTableName)
	_, err := m.dbConn.Exec(query, migration.Name(), number)
	if err != nil {
		log.Fatalf("error removing migration: %v\n", err)
	}
}
