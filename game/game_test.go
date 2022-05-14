package game

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCanUpgradeFactory(t *testing.T) {
	game := &Game{
		IronFactory: NewFactory(Iron, time.Second),
	}
	t.Run("upgrade in progress", func(t *testing.T) {
		game.IronFactory.Status.IsInProgress = true
		err := game.CanUpgrade(Iron)
		assert.Equal(t, "upgrade allready in progress", err.Error())
	})

	t.Run("upgrade not in progress", func(t *testing.T) {
		game.IronFactory.Status.IsInProgress = false
		err := game.CanUpgrade(Iron)
		assert.NoError(t, err)
	})
}

func TestUpgradeFactory(t *testing.T) {
	game := &Game{
		IronFactory: NewFactory(Iron, time.Second),
		Resources:   &Resources{},
	}

	t.Run("without enough resources", func(t *testing.T) {
		err := game.UpgradeFactory(Iron)
		assert.Equal(t, "payment error: not enough iron for upgrade, not enough copper for upgrade, not enough gold for upgrade", err.Error())
	})

	t.Run("with enough resources", func(t *testing.T) {
		game.Resources = &Resources{
			Gold:   2,
			Iron:   500,
			Copper: 200,
		}
		err := game.UpgradeFactory(Iron)
		assert.NoError(t, err)
	})
}
