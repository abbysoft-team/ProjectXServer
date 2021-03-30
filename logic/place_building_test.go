package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleLogic_PlaceBuilding(t *testing.T) {
	logic, db, session := NewLogicMock()
	request := &rpc.PlaceBuildingRequest{
		SessionID:  "sessionID",
		Location:   &rpc.Vector2D{X: 1.0, Y: 0.5},
		Rotation:   25.4,
		TownID:     1,
		BuildingID: rpc.BuildingType_HOUSE,
	}

	session.SelectedCharacter = &model.Character{
		ID:        1,
		AccountID: 1,
		Name:      "test",
	}
	session.SelectedCharacter.Towns = append(session.SelectedCharacter.Towns, model.Town{
		ID:         1,
		X:          0,
		Y:          0,
		OwnerName:  "",
		Population: 0,
		Name:       "",
		Buildings:  nil,
		Rotation:   0,
	})

	building := model.Buildings[request.BuildingID]
	building.Rotation = request.Rotation
	building.Location = model.ToModelVector(request.Location)

	session.SelectedCharacter.Resources = building.Cost
	session.SelectedCharacter.MaxPopulation = 10

	db.On("AddTownBuilding", int64(1), building).Return(nil)
	db.On("UpdateCharacter", mock.MatchedBy(func(char model.Character) bool {
		return char.MaxPopulation == 10+building.PopulationBonus &&
			!char.Resources.IsEnough(model.Resources{
				CharacterID: 0,
				Wood:        1,
				Food:        1,
				Stone:       1,
				Leather:     1,
			})
	})).Return(nil)

	resp, err := logic.PlaceBuilding(session, request)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
