package main

import "time"

const (
	cacheTTLSeconds  = 10
	baseDBReadPeriod = 1
	gggbbbWorld      = "gggbbb"
)

type WorldDataRepository interface {
	GetWorldData() []WorldInformation
}

// Cache should only handle getting world data. Posting can happen in real time as it does not happen as often.
// TODO we may also want to cache ip address so we don't have to grab them from the db
// We should probably set up metrics to see how often posts even happen.
type Cache struct {
	cachedData []WorldInformation
	repository WorldDataRepository
}

func NewCache(repository WorldDataRepository) *Cache {
	// Init cache and goroutine that periodically grabs cached data.
	// If there are no gggbbb worlds, we should query the db

	c := &Cache{
		cachedData: repository.GetWorldData(),
		repository: repository,
	}

	go func() {
		for {
			c.refresh()
		}
	}()

	return c
}

func (c *Cache) Get() []WorldInformation {
	return c.cachedData
}

func (c *Cache) refresh() {
	// Read the db for latest data
	c.cachedData = c.repository.GetWorldData()

	// If we already have a gggbbb world, we don't have to serve the most up-to-date information
	if c.gggbbbWorldIsFound() {
		time.Sleep(cacheTTLSeconds * time.Second)
		return
	}

	// Wait 1s between reads at minimum to not overload the db
	time.Sleep(baseDBReadPeriod * time.Second)
}

func (c *Cache) gggbbbWorldIsFound() bool {
	for _, w := range c.cachedData {
		if w.StreamOrder == gggbbbWorld {
			return true
		}
	}

	return false
}
