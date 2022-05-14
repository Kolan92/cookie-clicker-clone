package game

import (
	"errors"
	"sync"
	"time"

	"github.com/openlyinc/pointy"
)

type Resource int

const (
	Iron Resource = iota
	Copper
	Gold
)

var ResourceFromString = map[string]Resource{
	"iron":   Iron,
	"copper": Copper,
	"gold":   Gold,
}

var ironFactoryUpgrades = map[int]FactoryLevel{
	1: FactoryLevel{
		Level:               1,
		ProductionPerSecond: 10,
		NextUpgradeDuration: 15,
		UpgradeCost: &UpgradeCost{
			Iron:   pointy.Int(300),
			Copper: pointy.Int(100),
			Gold:   pointy.Int(1),
		},
	},
	2: FactoryLevel{
		Level:               2,
		ProductionPerSecond: 20,
		NextUpgradeDuration: 30,
		UpgradeCost: &UpgradeCost{
			Iron:   pointy.Int(800),
			Copper: pointy.Int(250),
			Gold:   pointy.Int(2),
		},
	},
	3: FactoryLevel{
		Level:               3,
		ProductionPerSecond: 40,
		NextUpgradeDuration: 60,
		UpgradeCost: &UpgradeCost{
			Iron:   pointy.Int(1600),
			Copper: pointy.Int(500),
			Gold:   pointy.Int(4),
		},
	},
	4: FactoryLevel{
		Level:               4,
		ProductionPerSecond: 80,
		NextUpgradeDuration: 90,
		UpgradeCost: &UpgradeCost{
			Iron:   pointy.Int(3000),
			Copper: pointy.Int(1000),
			Gold:   pointy.Int(8),
		},
	},
	5: FactoryLevel{
		Level:               5,
		ProductionPerSecond: 150,
		NextUpgradeDuration: 120,
		UpgradeCost:         &UpgradeCost{},
	},
}

type FactoryLevel struct {
	Level               int
	ProductionPerSecond int
	NextUpgradeDuration int
	UpgradeCost         *UpgradeCost
}
type UpgradeCost struct {
	Iron   *int
	Copper *int
	Gold   *int
}

type UpgradeStatus struct {
	IsInProgress           bool
	RemainingTimeInSeconds *int
	NextUpgradeCost        *int
}

type baseFactory struct {
	Level          int
	ProductionRate int
}

type Factory interface {
	Upgrade() error
}

type IronFactory struct {
	baseFactory
}

func (f *IronFactory) InitializeProduction(resources chan<- int) {
	for _ = range time.Tick(time.Second) {
		resources <- 10
	}
}

func (f *IronFactory) Upgrade() error {
	return errors.New("Not enough resources")
}

type ResourcesStock struct {
	Iron        int
	Copper      int
	Gold        int
	ironLock    sync.RWMutex
	ironFactory *IronFactory
}

func NewResourcesStock() *ResourcesStock {
	stock := &ResourcesStock{
		ironFactory: &IronFactory{},
	}

	ironProduction := make(chan int, 1)
	go stock.ironFactory.InitializeProduction(ironProduction)
	go func() {
		for {
			select {
			case iron := <-ironProduction:
				stock.ironLock.Lock()
				stock.Iron += iron
				stock.ironLock.Unlock()
			}
		}
	}()

	return stock
}

func (r *ResourcesStock) Upgrade(resource Resource) error {
	switch resource {
	case Iron:
		return r.ironFactory.Upgrade()
	case Copper, Gold:
		return errors.New("Not implemented")
	default:
		return errors.New("Unknown resource")
	}

}
