package channel

import (
	"context"
	"errors"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/Chronicle20/atlas-rest/requests"
	"github.com/sirupsen/logrus"
	"math/rand"
)

func ByIdModelProvider(l logrus.FieldLogger, ctx context.Context) func(worldId byte, channelId byte) model.Provider[Model] {
	return func(worldId byte, channelId byte) model.Provider[Model] {
		return requests.Provider[RestModel, Model](l, ctx)(requestChannel(worldId, channelId), Extract)
	}
}

func GetById(l logrus.FieldLogger, ctx context.Context) func(worldId byte, channelId byte) (Model, error) {
	return func(worldId byte, channelId byte) (Model, error) {
		return ByIdModelProvider(l, ctx)(worldId, channelId)()
	}
}

func ByWorldModelProvider(l logrus.FieldLogger, ctx context.Context) func(worldId byte) model.Provider[[]Model] {
	return func(worldId byte) model.Provider[[]Model] {
		return requests.SliceProvider[RestModel, Model](l, ctx)(requestChannelsForWorld(worldId), Extract, model.Filters[Model]())
	}
}

func GetForWorld(l logrus.FieldLogger, ctx context.Context) func(worldId byte) ([]Model, error) {
	return func(worldId byte) ([]Model, error) {
		return ByWorldModelProvider(l, ctx)(worldId)()
	}
}

func GetRandomInWorld(l logrus.FieldLogger, ctx context.Context) func(worldId byte) (Model, error) {
	return func(worldId byte) (Model, error) {
		cs, err := GetForWorld(l, ctx)(worldId)
		if err != nil {
			return Model{}, err
		}
		if len(cs) == 0 {
			return Model{}, errors.New("no channel found")
		}

		ri := rand.Intn(len(cs))
		return cs[ri], nil
	}
}
