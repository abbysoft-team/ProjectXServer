package model

import rpc "abbysoft/gardarike-online/rpc/generated"

type Building struct {
	ID              rpc.BuildingType
	Name            string
	Cost            Resources
	Production      Resources
	Location        Vector2D
	Rotation        Vector2D
	PopulationBonus uint64
}

// CharacterBuildings - number of buildings of each type
type CharacterBuildings map[rpc.BuildingType]uint64

var (
	Buildings = map[rpc.BuildingType]Building{
		rpc.BuildingType_HOUSE: {
			ID:              rpc.BuildingType_HOUSE,
			Name:            "house",
			Cost:            Resources{Wood: 30, Food: 10, Stone: 15, Leather: 20},
			Production:      Resources{Food: 1},
			PopulationBonus: 5,
		},
		rpc.BuildingType_QUARRY: {
			ID:              rpc.BuildingType_QUARRY,
			Name:            "quarry",
			Cost:            Resources{Wood: 100, Food: 50, Stone: 0, Leather: 80},
			Production:      Resources{Stone: 1},
			PopulationBonus: 0,
		},
	}
)

func IsValidBuildingType(typeValue int32) bool {
	_, found := rpc.BuildingType_name[typeValue]
	return found
}
