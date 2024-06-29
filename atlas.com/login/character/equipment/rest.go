package equipment

import (
	"atlas-login/character/equipment/slot"
)

type RestModel struct {
	Hat      slot.RestModel `json:"hat"`
	Medal    slot.RestModel `json:"medal"`
	Forehead slot.RestModel `json:"forehead"`
	Ring1    slot.RestModel `json:"ring1"`
	Ring2    slot.RestModel `json:"ring2"`
	Eye      slot.RestModel `json:"eye"`
	Earring  slot.RestModel `json:"earring"`
	Shoulder slot.RestModel `json:"shoulder"`
	Cape     slot.RestModel `json:"cape"`
	Top      slot.RestModel `json:"top"`
	Pendant  slot.RestModel `json:"pendant"`
	Weapon   slot.RestModel `json:"weapon"`
	Shield   slot.RestModel `json:"shield"`
	Gloves   slot.RestModel `json:"gloves"`
	Bottom   slot.RestModel `json:"bottom"`
	Belt     slot.RestModel `json:"belt"`
	Ring3    slot.RestModel `json:"ring3"`
	Ring4    slot.RestModel `json:"ring4"`
	Shoes    slot.RestModel `json:"shoes"`
}

func Transform(model Model) RestModel {
	rm := RestModel{
		Hat:      slot.Transform(model.hat),
		Medal:    slot.Transform(model.medal),
		Forehead: slot.Transform(model.forehead),
		Ring1:    slot.Transform(model.ring1),
		Ring2:    slot.Transform(model.ring2),
		Eye:      slot.Transform(model.eye),
		Earring:  slot.Transform(model.earring),
		Shoulder: slot.Transform(model.shoulder),
		Cape:     slot.Transform(model.cape),
		Top:      slot.Transform(model.top),
		Pendant:  slot.Transform(model.pendant),
		Weapon:   slot.Transform(model.weapon),
		Shield:   slot.Transform(model.shield),
		Gloves:   slot.Transform(model.gloves),
		Bottom:   slot.Transform(model.bottom),
		Belt:     slot.Transform(model.belt),
		Ring3:    slot.Transform(model.ring3),
		Ring4:    slot.Transform(model.ring4),
		Shoes:    slot.Transform(model.shoes),
	}
	return rm
}

func Extract(model RestModel) Model {
	return Model{
		hat:      slot.Extract(model.Hat),
		medal:    slot.Extract(model.Medal),
		forehead: slot.Extract(model.Forehead),
		ring1:    slot.Extract(model.Ring1),
		ring2:    slot.Extract(model.Ring2),
		eye:      slot.Extract(model.Eye),
		earring:  slot.Extract(model.Earring),
		shoulder: slot.Extract(model.Shoulder),
		cape:     slot.Extract(model.Cape),
		top:      slot.Extract(model.Top),
		pendant:  slot.Extract(model.Pendant),
		weapon:   slot.Extract(model.Weapon),
		shield:   slot.Extract(model.Shield),
		gloves:   slot.Extract(model.Gloves),
		bottom:   slot.Extract(model.Bottom),
		belt:     slot.Extract(model.Belt),
		ring3:    slot.Extract(model.Ring3),
		ring4:    slot.Extract(model.Ring4),
		shoes:    slot.Extract(model.Shoes),
	}
}
