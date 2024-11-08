package oreillyapi

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	client := New()

	res, err := client.Search(SearchOption{
		// Query: "*",
		Publishers: []string{
			"O'Reilly Japan, Inc.",
		},
		Languages: []string{
			"ja",
		},
		Sort:  PublicationDate,
		Order: Desc,
		Limit: 1,
	})
	require.NoError(t, err)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(&res)
	require.NoError(t, err)
}
