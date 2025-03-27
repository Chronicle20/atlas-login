package equipable

import (
	"strconv"
	"time"
)

type RestModel struct {
	Id             uint32    `json:"-"`
	ItemId         uint32    `json:"itemId"`
	Slot           int16     `json:"slot"`
	Strength       uint16    `json:"strength"`
	Dexterity      uint16    `json:"dexterity"`
	Intelligence   uint16    `json:"intelligence"`
	Luck           uint16    `json:"luck"`
	HP             uint16    `json:"hp"`
	MP             uint16    `json:"mp"`
	WeaponAttack   uint16    `json:"weaponAttack"`
	MagicAttack    uint16    `json:"magicAttack"`
	WeaponDefense  uint16    `json:"weaponDefense"`
	MagicDefense   uint16    `json:"magicDefense"`
	Accuracy       uint16    `json:"accuracy"`
	Avoidability   uint16    `json:"avoidability"`
	Hands          uint16    `json:"hands"`
	Speed          uint16    `json:"speed"`
	Jump           uint16    `json:"jump"`
	Slots          uint16    `json:"slots"`
	OwnerName      string    `json:"ownerName"`
	Locked         bool      `json:"locked"`
	Spikes         bool      `json:"spikes"`
	KarmaUsed      bool      `json:"karmaUsed"`
	Cold           bool      `json:"cold"`
	CanBeTraded    bool      `json:"canBeTraded"`
	LevelType      byte      `json:"levelType"`
	Level          byte      `json:"level"`
	Experience     uint32    `json:"experience"`
	HammersApplied uint32    `json:"hammersApplied"`
	Expiration     time.Time `json:"expiration"`
}

func (r RestModel) GetName() string {
	return "equipables"
}

func (r RestModel) GetID() string {
	return strconv.Itoa(int(r.Id))
}

func Transform(m Model) (RestModel, error) {
	rm := RestModel{
		Id:             m.id,
		ItemId:         m.itemId,
		Slot:           m.slot,
		Strength:       m.strength,
		Dexterity:      m.dexterity,
		Intelligence:   m.intelligence,
		Luck:           m.luck,
		HP:             m.hp,
		MP:             m.mp,
		WeaponAttack:   m.weaponAttack,
		MagicAttack:    m.magicAttack,
		WeaponDefense:  m.weaponDefense,
		MagicDefense:   m.magicDefense,
		Accuracy:       m.accuracy,
		Avoidability:   m.avoidability,
		Hands:          m.hands,
		Speed:          m.speed,
		Jump:           m.jump,
		Slots:          m.slots,
		OwnerName:      m.ownerName,
		Locked:         m.locked,
		Spikes:         m.spikes,
		KarmaUsed:      m.karmaUsed,
		Cold:           m.cold,
		CanBeTraded:    m.canBeTraded,
		LevelType:      m.levelType,
		Level:          m.level,
		Experience:     m.experience,
		HammersApplied: m.hammersApplied,
		Expiration:     m.expiration,
	}
	return rm, nil
}

func Extract(rm RestModel) (Model, error) {
	return Model{
		id:             rm.Id,
		itemId:         rm.ItemId,
		slot:           rm.Slot,
		strength:       rm.Strength,
		dexterity:      rm.Dexterity,
		intelligence:   rm.Intelligence,
		luck:           rm.Luck,
		hp:             rm.HP,
		mp:             rm.MP,
		weaponAttack:   rm.WeaponAttack,
		magicAttack:    rm.MagicAttack,
		weaponDefense:  rm.WeaponDefense,
		magicDefense:   rm.MagicDefense,
		accuracy:       rm.Accuracy,
		avoidability:   rm.Avoidability,
		hands:          rm.Hands,
		speed:          rm.Speed,
		jump:           rm.Jump,
		slots:          rm.Slots,
		ownerName:      rm.OwnerName,
		locked:         rm.Locked,
		spikes:         rm.Spikes,
		karmaUsed:      rm.KarmaUsed,
		cold:           rm.Cold,
		canBeTraded:    rm.CanBeTraded,
		levelType:      rm.LevelType,
		level:          rm.Level,
		experience:     rm.Experience,
		hammersApplied: rm.HammersApplied,
		expiration:     rm.Expiration,
	}, nil
}
