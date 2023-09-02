package cron

import (
	"context"
	"fmt"
	"time"

	"dyploma/structs"

	"github.com/go-co-op/gocron"
)

type Crawler interface {
	PageWalker(page string, exclude []string, onlyThisPage, forceCollect bool, headers map[string]string) (siteStruct structs.SiteStruct, err error)
}

type Postgres interface {
	GetSites() ([]structs.SiteStruct, error)
}

type Cron struct {
	s        *gocron.Scheduler
	crawler  Crawler
	postgres Postgres
}

func (c *Cron) Start(_ context.Context) error {
	c.s = gocron.NewScheduler(time.UTC)
	_, err := c.s.Every(1).Day().Do(c.UpdateSites)
	if err != nil {
		return fmt.Errorf("Cron start error: %w", err)
	}

	c.s.StartAsync()
	return nil
}

func (c *Cron) Stop(_ context.Context) error {
	if c.s != nil {
		c.s.Stop()
		c.s.Clear()
	}

	return nil
}

func NewCron(crawler Crawler, postgres Postgres) *Cron {
	return &Cron{crawler: crawler, postgres: postgres}
}
