package tests

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestGetWorldMap(t *testing.T) {
	TestLoginSuccessful(t)

	rand.Seed(time.Now().UnixNano())

	var request rpc.Request
	request.Data = &rpc.Request_GetWorldMapRequest{
		GetWorldMapRequest: &rpc.GetWorldMapRequest{
			Location: &rpc.IntVector2D{
				X: int32(rand.Intn(100) - 50),
				Y: int32(rand.Intn(100) - 50),
			},
			SessionID: sessionID,
		},
	}

	resp, err := client.SendRequest(request)

	if !assert.NoError(t, err, "request error is not nil") {
		return
	}
	if !assert.NotNil(t, resp, "response is nil") {
		return
	}
	if !assert.NotNil(t, resp.GetGetWorldMapResponse(), "response isn't a get world map response") {
		return
	}

	assert.NotEmpty(t, resp.GetGetWorldMapResponse().Map.Data)
}
