package factory

import (
	"atlas-login/character"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

func SeedCharacter(l logrus.FieldLogger, ctx context.Context) func(accountId uint32, worldId byte, name string, jobIndex uint32, subJobIndex uint16, face uint32, hair uint32, color uint32, skinColor uint32, gender byte, top uint32, bottom uint32, shoes uint32, weapon uint32, strength byte, dexterity byte, intelligence byte, luck byte) (character.Model, error) {
	return func(accountId uint32, worldId byte, name string, jobIndex uint32, subJobIndex uint16,
		face uint32, hair uint32, color uint32, skinColor uint32, gender byte,
		top uint32, bottom uint32, shoes uint32, weapon uint32,
		strength byte, dexterity byte, intelligence byte, luck byte) (character.Model, error) {
		l.Debugf("Create character [%s] with job [%d:%d] and gender [%d].", name, jobIndex, subJobIndex, gender)
		l.Debugf("Face [%d], Hair [%d], HairColor [%d] SkinColor [%d].", face, hair, color, skinColor)
		l.Debugf("Top [%d], Bottom [%d], Shoes [%d], Weapon [%d].", top, bottom, shoes, weapon)
		l.Debugf("Strength [%d], Dexterity [%d], Intelligence [%d], Luck [%d].", strength, dexterity, intelligence, luck)
		c, err := requestCreate(accountId, worldId, name, jobIndex, subJobIndex, face, hair, color, skinColor, gender, top, bottom, shoes, weapon, strength, dexterity, intelligence, luck)(l, ctx)
		if err != nil {
			return character.Model{}, err
		}
		m, err := character.Extract(c)
		if err != nil {
			return character.Model{}, err
		}
		if m.Id() == 0 {
			return character.Model{}, errors.New("unable to create character")
		}
		return m, nil
	}
}
