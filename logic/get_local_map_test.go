package logic

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/assert"
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
