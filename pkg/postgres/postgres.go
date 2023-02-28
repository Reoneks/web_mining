package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"test/config"
	"test/structs"

	"github.com/golang-migrate/migrate/v4"
	mpostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

type Postgres struct {
	db *gorm.DB
}

func (p *Postgres) SaveSiteStruct(siteStruct structs.SiteStruct) error {
	return nil
}

func (p *Postgres) GetFullData(link string) (structs.SiteStruct, error) {
	return structs.SiteStruct{}, nil
}

func (p *Postgres) GetCrawlerData(link string) (structs.CrawlerData, error) {
	return structs.CrawlerData{}, nil
}

func NewPostgres(cfg *config.Config) (postgres *Postgres, err error) {
	postgres = new(Postgres)
	postgres.db, err = newDB(cfg.DSN, cfg.MigrationURL)
	if err != nil {
		return
	}

	return
}

func (p *Postgres) Stop(ctx context.Context) error {
	sql, err := p.db.WithContext(ctx).DB()
	if err != nil {
		return fmt.Errorf("Postgres.Stop get database error: %w", err)
	}

	err = sql.Close()
	if err != nil {
		return fmt.Errorf("Postgres.Stop connection close error: %w", err)
	}

	return nil
}

func newDB(dsn, migrationsURL string) (*gorm.DB, error) {
	if err := migrations(dsn, migrationsURL); err != nil {
		return nil, fmt.Errorf("Postgres.newDB migrations error: %w", err)
	}

	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: glog.Default.LogMode(glog.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("Postgres.newDB gorm connection open error: %w", err)
	}

	return client, nil
}

func migrations(dsn, migrationsURL string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("Postgres.migrations sql connection open error: %w", err)
	}

	driver, err := mpostgres.WithInstance(db, &mpostgres.Config{})
	if err != nil {
		return fmt.Errorf("Postgres.migrations create migration instance error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("Postgres.migrations migration init error: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("Postgres.migrations migration error: %w", err)
	}

	return nil
}
