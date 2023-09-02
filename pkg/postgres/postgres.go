package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"dyploma/config"
	"dyploma/structs"

	"github.com/golang-migrate/migrate/v4"
	mpostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Need for migrations
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

type Postgres struct {
	db *gorm.DB
}

func (p *Postgres) SaveSiteStruct(siteStruct structs.SiteStruct, onlyThisPage, force bool) error {
	baseReq := p.db.Model(&siteStruct).Session(&gorm.Session{FullSaveAssociations: true})

	if onlyThisPage && force {
		baseReq.Omit("Hierarchy.ParentLink")
	}

	err := baseReq.Where("base_url = ?", siteStruct.BaseURL).Save(&siteStruct).Error
	if err != nil {
		return fmt.Errorf("Error saving site structure: %w", err)
	}

	return nil
}

func (p *Postgres) GetFullData(link, url string, onlyThisPage bool) (structs.SiteStruct, error) {
	var result structs.SiteStruct

	baseReq := p.db.Model(&result)
	if !onlyThisPage {
		baseReq = baseReq.Preload("Hierarchy", preloadHierarchy)
	}

	err := baseReq.Where("base_url = ?", link).First(&result).Error
	if err != nil {
		return structs.SiteStruct{}, fmt.Errorf("Error getting site structure: %w", err)
	}

	if onlyThisPage {
		crawlerData, err := p.GetCrawlerData(url)
		if err != nil {
			return structs.SiteStruct{}, fmt.Errorf("Error getting crawler data: %w", err)
		}

		result.Hierarchy = &structs.Hierarchy{CrawlerData: crawlerData, ParentLink: link}
	}

	return result, nil
}

func (p *Postgres) GetSites() ([]structs.SiteStruct, error) {
	var result []structs.SiteStruct

	err := p.db.Model(&structs.SiteStruct{}).Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("Error getting sites: %w", err)
	}

	return result, nil
}

func (p *Postgres) GetCrawlerData(link string) (structs.CrawlerData, error) {
	var result structs.CrawlerData

	err := p.db.Model(&structs.Hierarchy{}).Where("link = ?", link).First(&result).Error
	if err != nil {
		return structs.CrawlerData{}, fmt.Errorf("Error getting crawler data: %w", err)
	}

	return result, nil
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
