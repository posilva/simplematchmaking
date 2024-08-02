package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	configOK = `{
  "queues": {
    "queue1": {
      "name": "Queue 1",
      "maxPlayers": 10,
      "nrBrackets": 5,
      "maxRanking": 1000,
      "minRanking": 100,
      "makeIterations": 3
    },
    "queue2": {
      "name": "Queue 2",
      "maxPlayers": 20,
      "nrBrackets": 10,
      "maxRanking": 2000,
      "minRanking": 200,
      "makeIterations": 5
    }
  },
  "matchmakers": {
    "matchmaker1": {
      "name": "Matchmaker 1",
      "intervalSecs": 30,
      "makeTimeoutSecs": 60
    },
    "matchmaker2": {
      "name": "Matchmaker 2",
      "intervalSecs": 45,
      "makeTimeoutSecs": 90
    }
  }
}`
)

func TestLoad(t *testing.T) {
	data := base64.StdEncoding.EncodeToString([]byte(configOK))
	err := os.Setenv("MATCHMAKING_CFG", fmt.Sprintf("json.%s", data))
	require.Nil(t, err)

	e := NewEnvVar()
	err = e.Load()
	require.NoError(t, err)

	cfg := e.Get()
	require.NotNil(t, cfg)
	b, err := json.Marshal(cfg)
	require.Nil(t, err)
	require.JSONEq(t, configOK, string(b))
}
