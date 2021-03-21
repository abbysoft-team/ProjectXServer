package tests

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRenameTown(t *testing.T) {
	TestSelectCharacter(t)

	var request rpc.Request
	request.Data = &rpc.Request_RenameTownRequest{
		RenameTownRequest: &rpc.RenameTownRequest{
			SessionID: sessionID,
			TownID:    4,
			NewName:   "New town name",
		},
	}

	resp, err := client.SendRequest(request)

	require.NoError(t, err, "request error is not nil")
	require.NotNil(t, resp, "response is nil")
	require.NotNil(t, resp.GetRenameTownResponse(), "response isn't a rename town response")
}
