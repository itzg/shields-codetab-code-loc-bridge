package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestCodeLoc(t *testing.T) {
	resp := httptest.NewRecorder()

	CodeLoc(resp, httptest.NewRequest("GET", "/api/codeloc?github=itzg/mc-image-helper&language=java", nil))

	assert.Equal(t, resp.Code, 200)

	var shieldsResp ShieldsEndpointResponse
	err := json.NewDecoder(resp.Result().Body).Decode(&shieldsResp)
	assert.NoError(t, err)

	assert.Equal(t, 1, shieldsResp.SchemaVersion)

	assert.Equal(t, "LOC (java)", shieldsResp.Label)

	assert.NotEmpty(t, shieldsResp.Message)
	assert.NotEqual(t, "", shieldsResp.Message)
}
