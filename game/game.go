package game

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Game struct {
	Resources *Resources

	ironLock    sync.RWMutex
	IronFactory *Factory

	copperLock    sync.RWMutex
	CopperFactory *Factory

	goldLock    sync.RWMutex
	GoldFactory *Factory
}

type Resources struct {
	Iron   int
	Copper int
	Gold   int
}

func NewGame() *Game {
	game := &Game{
		Resources:     &Resources{},
		IronFactory:   NewFactory(Iron, time.Second),
		CopperFactory: NewFactory(Copper, time.Second),
		GoldFactory:   NewFactory(Gold, time.Minute),
	}

	ironProduction := make(chan int, 1)
	copperProduction := make(chan int, 1)
	goldProduction := make(chan int, 1)

	go game.IronFactory.InitializeProduction(ironProduction)
	go game.CopperFactory.InitializeProduction(copperProduction)
	go game.GoldFactory.InitializeProduction(goldProduction)

	go func() {
		for {
			select {
			case iron := <-ironProduction:
				game.ironLock.Lock()
				game.Resources.Iron += iron
				game.ironLock.Unlock()
			case copper := <-copperProduction:
				game.copperLock.Lock()
				game.Resources.Copper += copper
				game.copperLock.Unlock()
			case gold := <-goldProduction:
				game.goldLock.Lock()
				game.Resources.Gold += gold
				game.goldLock.Unlock()
			}
		}
	}()

	return game
}

func (g *Game) UpgradeFactory(resource Resource) error {
	if err := g.CanUpgrade(resource); err != nil {
		return err
	}

	if err := g.PayForUpgrade(resource); err != nil {
		return err
	}

	switch resource {
	case Iron:
		go g.IronFactory.Upgrade()
		return nil
	case Copper:
		go g.CopperFactory.Upgrade()
		return nil
	case Gold:
		go g.GoldFactory.Upgrade()
		return nil
	default:
		return errors.New("unknown resource")
	}
}

func (g *Game) PayForUpgrade(resource Resource) error {

	var upgradeCost *UpgradeCost
	switch resource {
	case Iron:
		upgradeCost = g.IronFactory.Level.UpgradeCost
	case Copper:
		upgradeCost = g.CopperFactory.Level.UpgradeCost
	case Gold:
		upgradeCost = g.GoldFactory.Level.UpgradeCost
	}
	g.ironLock.Lock()
	g.copperLock.Lock()
	g.goldLock.Lock()

	defer g.ironLock.Unlock()
	defer g.copperLock.Unlock()
	defer g.goldLock.Unlock()

	var ironErr, copperErr, goldErr error
	if upgradeCost.Copper != nil && *upgradeCost.Copper > g.Resources.Copper {
		copperErr = errors.New("not enough copper for upgrade")
	}

	if upgradeCost.Iron != nil && *upgradeCost.Iron > g.Resources.Iron {
		ironErr = errors.New("not enough iron for upgrade")
	}

	if upgradeCost.Gold != nil && *upgradeCost.Gold > g.Resources.Gold {
		goldErr = errors.New("not enough gold for upgrade")
	}

	if ironErr != nil || copperErr != nil || goldErr != nil {

		return fmt.Errorf("payment error: %v, %v, %v", ironErr, copperErr, goldErr)
	}
	fmt.Println("Paying for upgrade")

	if upgradeCost.Copper != nil {
		g.Resources.Copper -= *upgradeCost.Copper
	}

	if upgradeCost.Iron != nil {
		g.Resources.Iron -= *upgradeCost.Iron
	}

	if upgradeCost.Gold != nil {
		g.Resources.Gold -= *upgradeCost.Gold
	}

	return nil
}

func (g *Game) CanUpgrade(resource Resource) error {
	switch resource {
	case Iron:
		if g.IronFactory.Status.IsInProgress {
			return errors.New("upgrade allready in progress")
		}
	case Copper:
		if g.CopperFactory.Status.IsInProgress {
			return errors.New("upgrade allready in progress")
		}
	case Gold:
		if g.GoldFactory.Status.IsInProgress {
			return errors.New("upgrade allready in progress")
		}
	default:
		return errors.New("Unknown resource")
	}
	return nil
}
