package inventory

import (
	"atlas-login/character/inventory/equipable"
	"atlas-login/character/inventory/item"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/jtumidanski/api2go/jsonapi"
)

type RestModel struct {
	Equipable EquipableRestModel `json:"equipable"`
	Useable   ItemRestModel      `json:"useable"`
	Setup     ItemRestModel      `json:"setup"`
	Etc       ItemRestModel      `json:"etc"`
	Cash      ItemRestModel      `json:"cash"`
}

type EquipableRestModel struct {
	Type     string                `json:"-"`
	Capacity uint32                `json:"capacity"`
	Items    []equipable.RestModel `json:"items"`
}

func (r EquipableRestModel) GetName() string {
	return "inventories"
}

func (r EquipableRestModel) GetID() string {
	return r.Type
}

func (r EquipableRestModel) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "equipables",
			Name: "equipables",
		},
	}
}

func (r EquipableRestModel) GetReferencedIDs() []jsonapi.ReferenceID {
	var result []jsonapi.ReferenceID
	for _, v := range r.Items {
		result = append(result, jsonapi.ReferenceID{
			ID:   v.GetID(),
			Type: "equipables",
			Name: "equipables",
		})
	}
	return result
}

func (r EquipableRestModel) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	var result []jsonapi.MarshalIdentifier
	for key := range r.Items {
		result = append(result, r.Items[key])
	}

	return result
}

type ItemRestModel struct {
	Type     string           `json:"-"`
	Capacity uint32           `json:"capacity"`
	Items    []item.RestModel `json:"items"`
}

func (r ItemRestModel) GetName() string {
	return "inventories"
}

func (r ItemRestModel) GetID() string {
	return r.Type
}

func (r ItemRestModel) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "items",
			Name: "items",
		},
	}
}

func (r ItemRestModel) GetReferencedIDs() []jsonapi.ReferenceID {
	var result []jsonapi.ReferenceID
	for _, v := range r.Items {
		result = append(result, jsonapi.ReferenceID{
			ID:   v.GetID(),
			Type: "items",
			Name: "items",
		})
	}
	return result
}

func (r ItemRestModel) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	var result []jsonapi.MarshalIdentifier
	for key := range r.Items {
		result = append(result, r.Items[key])
	}

	return result
}

func Extract(model RestModel) (Model, error) {
	e, err := ExtractEquipable(model.Equipable)
	if err != nil {
		return Model{}, err
	}
	use, err := ExtractItem(model.Useable)
	if err != nil {
		return Model{}, err
	}
	setup, err := ExtractItem(model.Setup)
	if err != nil {
		return Model{}, err
	}
	etc, err := ExtractItem(model.Etc)
	if err != nil {
		return Model{}, err
	}
	cash, err := ExtractItem(model.Cash)
	if err != nil {
		return Model{}, err
	}

	return Model{
		equipable: e,
		useable:   use,
		setup:     setup,
		etc:       etc,
		cash:      cash,
	}, nil
}

func ExtractItem(rm ItemRestModel) (ItemModel, error) {
	is, err := model.SliceMap[item.RestModel, item.Model](item.Extract)(model.FixedProvider(rm.Items))(model.ParallelMap())()
	if err != nil {
		return ItemModel{}, err
	}

	return ItemModel{
		capacity: rm.Capacity,
		items:    is,
	}, nil
}

func ExtractEquipable(rm EquipableRestModel) (EquipableModel, error) {
	es, err := model.SliceMap(equipable.Extract)(model.FixedProvider(rm.Items))(model.ParallelMap())()
	if err != nil {
		return EquipableModel{}, err
	}

	return EquipableModel{
		capacity: rm.Capacity,
		items:    es,
	}, nil
}
