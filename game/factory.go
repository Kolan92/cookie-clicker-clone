package game

import (
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

var factoryUpgrades = map[Resource]map[int]FactoryLevel{
	Iron: {
		1: FactoryLevel{
			Id:                  1,
			Production:          10,
			NextUpgradeDuration: 15,
			UpgradeCost: &UpgradeCost{
				Iron:   pointy.Int(300),
				Copper: pointy.Int(100),
				Gold:   pointy.Int(1),
			},
		},
		2: FactoryLevel{
			Id:                  2,
			Production:          20,
			NextUpgradeDuration: 30,
			UpgradeCost: &UpgradeCost{
				Iron:   pointy.Int(800),
				Copper: pointy.Int(250),
				Gold:   pointy.Int(2),
			},
		},
		3: FactoryLevel{
			Id:                  3,
			Production:          40,
			NextUpgradeDuration: 60,
			UpgradeCost: &UpgradeCost{
				Iron:   pointy.Int(1600),
				Copper: pointy.Int(500),
				Gold:   pointy.Int(4),
			},
		},
		4: FactoryLevel{
			Id:                  4,
			Production:          80,
			NextUpgradeDuration: 90,
			UpgradeCost: &UpgradeCost{
				Iron:   pointy.Int(3000),
				Copper: pointy.Int(1000),
				Gold:   pointy.Int(8),
			},
		},
		5: FactoryLevel{
			Id:                  5,
			Production:          150,
			NextUpgradeDuration: 120,
			UpgradeCost:         &UpgradeCost{},
		},
	},
	Copper: {
		1: FactoryLevel{
			Id:                  1,
			Production:          3,
			NextUpgradeDuration: 15,
			UpgradeCost: &UpgradeCost{
				Iron:   pointy.Int(200),
				Copper: pointy.Int(70),
			},
		},
		2: FactoryLevel{
			Id:                  2,
			Production:          7,
			NextUpgradeDuration: 30,
			UpgradeCost: &UpgradeCost{
				Iron:   pointy.Int(400),
				Copper: pointy.Int(150),
			},
		},
		3: FactoryLevel{
			Id:                  3,
			Production:          14,
			NextUpgradeDuration: 60,
			UpgradeCost: &UpgradeCost{
				Iron:   pointy.Int(800),
				Copper: pointy.Int(300),
			},
		},
		4: FactoryLevel{
			Id:                  4,
			Production:          30,
			NextUpgradeDuration: 90,
			UpgradeCost: &UpgradeCost{
				Iron:   pointy.Int(1600),
				Copper: pointy.Int(600),
			},
		},
		5: FactoryLevel{
			Id:                  5,
			Production:          60,
			NextUpgradeDuration: 120,
			UpgradeCost:         &UpgradeCost{},
		},
	},

	Gold: {
		1: FactoryLevel{
			Id:                  1,
			Production:          2,
			NextUpgradeDuration: 15,
			UpgradeCost: &UpgradeCost{
				Copper: pointy.Int(100),
				Gold:   pointy.Int(2),
			},
		},
		2: FactoryLevel{
			Id:                  2,
			Production:          3,
			NextUpgradeDuration: 30,
			UpgradeCost: &UpgradeCost{
				Copper: pointy.Int(200),
				Gold:   pointy.Int(4),
			},
		},
		3: FactoryLevel{
			Id:                  3,
			Production:          4,
			NextUpgradeDuration: 60,
			UpgradeCost: &UpgradeCost{
				Copper: pointy.Int(400),
				Gold:   pointy.Int(8),
			},
		},
		4: FactoryLevel{
			Id:                  4,
			Production:          6,
			NextUpgradeDuration: 90,
			UpgradeCost: &UpgradeCost{
				Copper: pointy.Int(800),
				Gold:   pointy.Int(16),
			},
		},
		5: FactoryLevel{
			Id:                  5,
			Production:          8,
			NextUpgradeDuration: 120,
			UpgradeCost:         &UpgradeCost{},
		},
	},
}

type FactoryLevel struct {
	Id                  int
	Production          int
	NextUpgradeDuration int
	UpgradeCost         *UpgradeCost
}

type UpgradeCost struct {
	Iron   *int `json:",omitempty"`
	Copper *int `json:",omitempty"`
	Gold   *int `json:",omitempty"`
}

type UpgradeStatus struct {
	IsInProgress           bool
	RemainingTimeInSeconds *int         `json:",omitempty"`
	NextUpgradeCost        *UpgradeCost `json:",omitempty"`
}

type Factory struct {
	Level              FactoryLevel
	FactoryType        Resource
	Status             *UpgradeStatus
	productionInterval time.Duration
}

func NewFactory(factoryType Resource, productionInterval time.Duration) *Factory {
	factory := &Factory{
		FactoryType:        factoryType,
		Level:              factoryUpgrades[factoryType][1],
		productionInterval: productionInterval,
	}
	status := &UpgradeStatus{
		IsInProgress:    false,
		NextUpgradeCost: factoryUpgrades[factoryType][factory.Level.Id+1].UpgradeCost,
	}

	factory.Status = status
	return factory
}

func (f *Factory) InitializeProduction(resources chan<- int) {
	for _ = range time.Tick(f.productionInterval) {
		resources <- factoryUpgrades[f.FactoryType][f.Level.Id].Production
	}
}

func (f *Factory) Upgrade() {
	nextUpgradeDuration := f.Level.NextUpgradeDuration
	f.Status.IsInProgress = true
	f.Status.RemainingTimeInSeconds = &nextUpgradeDuration
	f.Status.NextUpgradeCost = nil

	for _ = range time.Tick(time.Second) {
		f.Status.RemainingTimeInSeconds = pointy.Int(*f.Status.RemainingTimeInSeconds - 1)
		if *f.Status.RemainingTimeInSeconds <= 0 {
			f.Level = factoryUpgrades[f.FactoryType][f.Level.Id+1]
			f.Status = &UpgradeStatus{
				IsInProgress:    false,
				NextUpgradeCost: factoryUpgrades[f.FactoryType][f.Level.Id+1].UpgradeCost,
			}

			break
		}
	}
}
