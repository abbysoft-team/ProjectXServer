package tests

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestGetEmpiresRating(t *testing.T) {
	TestSelectCharacter(t)

	rand.Seed(time.Now().UnixNano())

	var request rpc.Request
	request.Data = &rpc.Request_GetEmpiresRatingRequest{
		GetEmpiresRatingRequest: &rpc.GetEmpiresRatingRequest{
			SessionID: sessionID,
			Criteria:  rpc.EmpiresRatingCriteria_POPULATION,
			Offset:    0,
			Limit:     4,
		},
	}

	resp, err := client.SendRequest(request)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetGetEmpiresRatingResponse())

	empireResp := resp.GetGetEmpiresRatingResponse()

	if empireResp.PlayerRating == nil {
		assert.True(t, len(empireResp.Entries) <= 4)
	} else {
		assert.Equal(t, characterName, empireResp.PlayerRating.EmpireName)
	}

	require.NotEmpty(t, empireResp.Entries)
	assert.NotNil(t, empireResp.Entries[0])
}
