package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func Test_getGlobalChunkCoordsForPosition(t *testing.T) {
	const mapChunkSize = 10

	testFunc := func(l rpc.Vector2D, x, y, number int64) {
		xReal, yReal, numberReal := getGlobalChunkCoordsForPosition(l, mapChunkSize)
		assert.Equal(t, x, xReal, "global chunk x is wrong")
		assert.Equal(t, y, yReal, "global chunk y is wrong")
		assert.Equal(t, number, numberReal, "local chunk number is wrong")
	}

	testFunc(rpc.Vector2D{}, 0, 0, 3)
	testFunc(rpc.Vector2D{X: -1, Y: 0}, -1, 0, 4)
	testFunc(rpc.Vector2D{X: -1, Y: 0}, -1, 0, 4)
	testFunc(rpc.Vector2D{X: -1, Y: -1}, -1, -1, 2)
	testFunc(rpc.Vector2D{X: -1, Y: 1}, -1, 0, 4)
	testFunc(rpc.Vector2D{X: 0, Y: -1}, 0, -1, 1)
	testFunc(rpc.Vector2D{X: mapChunkSize, Y: 0}, 1, 0, 3)
	testFunc(rpc.Vector2D{X: mapChunkSize, Y: -mapChunkSize}, 1, -1, 3)
}

func TestSimpleLogic_GetLocalMap(t *testing.T) {
	logic, db, session, generator := NewLogicMockWithTerrainGenerator()
	request := &rpc.GetLocalMapRequest{
		SessionID: "sessionID",
		Location: &rpc.Vector2D{
			X: float32(rand.Intn(10000)),
			Y: float32(-rand.Intn(100000)),
		},
	}

	logic.config.AlwaysRegenerateMap = false
	logic.config.ChunkSize = 10

	db.On("GetMapChunk", mock.Anything, mock.Anything, mock.MatchedBy(func(number int64) bool {
		return number >= 0 && number <= 4
	})).Return(model.WorldMapChunk{}, sql.ErrNoRows)
	db.On("SaveMapChunkOrUpdate", mock.Anything, mock.Anything).Return(nil)

	generator.On("GenerateTerrain", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]float32{10.0, 10.0})

	response, err := logic.GetLocalMap(session, request)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Map)
	require.Equal(t, response.Map.Data, []float32{10.0, 10.0})
}
