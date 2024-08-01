package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/segmentio/ksuid"
	"golang.org/x/exp/rand"

	"github.com/posilva/simplematchmaking/internal/adapters/input/handler"
	"github.com/stretchr/testify/suite"
)

var (
	baseURL       = "localhost:8808"
	defaultScheme = "http"
)

type E2ETestSuite struct {
	BaseTestSuite
}

func (suite *E2ETestSuite) SetupSuite() {
	setup(&suite.BaseTestSuite)
	baseURL = suite.ServiceEndpoint
}

func (suite *E2ETestSuite) TearDownSuite() {
	fmt.Println("Running teardown suite.Require()")
	teardown(&suite.BaseTestSuite)
}

func (suite *E2ETestSuite) TestCheckMatch() {
	fmResp, err := findMatchRequest()
	suite.Require().NoError(err)
	var fmOut handler.FindMatchOutput
	err = json.Unmarshal([]byte(fmResp), &fmOut)
	suite.Require().NoError(err)
	suite.Require().Regexp(regexp.MustCompile("^[a-zA-Z0-9]{27}$"), fmOut.TicketID)
	_, err = getMatchRequest(fmOut.TicketID, http.StatusNotFound)
	suite.Require().NoError(err)
}

func (suite *E2ETestSuite) TestCancelMatch() {
	resp, err := findMatchRequest()
	suite.Require().NoError(err)
	var out handler.FindMatchOutput
	err = json.Unmarshal([]byte(resp), &out)
	suite.Require().NoError(err)
	suite.Require().Regexp(regexp.MustCompile("^[a-zA-Z0-9]{27}$"), out.TicketID)
	err = cancelMatchRequest(out.TicketID)
	suite.Require().NoError(err)
}
func (suite *E2ETestSuite) SetupTest() {
	cmd := suite.RedisClient.B().Flushall().Build()
	err := suite.RedisClient.Do(suite.Context, cmd).Error()
	suite.Require().NoError(err)

	cmd = suite.RedisClient.B().Keys().Pattern("*").Build()
	keys, err := suite.RedisClient.Do(suite.Context, cmd).AsStrSlice()
	suite.Require().NoError(err)
	suite.Require().Empty(keys)
}
func (suite *E2ETestSuite) TestMatchmaking() {
	maxPlayers := 1000
	min := 0
	max := 1000

	var mapPlayers = make(map[string]string)
	var mapTickets = make(map[string]string)

	// generate a list of maxPlayers
	for i := 0; i < maxPlayers; i++ {
		playerID := ksuid.New().String()
		score := rand.Intn(max-min) + min
		response, err := findMatchRequestWithInput(handler.FindMatchInput{
			PlayerID: playerID,
			Score:    score,
		})
		suite.Require().NoError(err)
		// create a ticket for each player
		var out handler.FindMatchOutput
		err = json.Unmarshal([]byte(response), &out)
		suite.Require().NoError(err)
		suite.Require().Regexp(regexp.MustCompile("^[a-zA-Z0-9]{27}$"), out.TicketID)
		mapPlayers[playerID] = out.TicketID
		mapTickets[out.TicketID] = playerID
		suite.Require().NoError(err)
	}

	mapMatches := make(map[string]string)
	matches := make(map[string]struct{})

	for len(matches) < maxPlayers/2 {
		for _, ticketID := range mapPlayers {
			if _, ok := mapMatches[ticketID]; ok {
				continue
			}
			// check if a match is available
			response, err := getMatchRequestAny(ticketID)
			if requests.HasStatusErr(err, http.StatusNotFound) {
				//suite.T().Logf("Match not found for ticketID: %s", ticketID)
				continue
			}
			// match is available
			var match handler.GetMatchOutput
			err = json.Unmarshal([]byte(response), &match)
			suite.Require().NoError(err)
			matches[match.MatchID] = struct{}{}
			// save the matchID for each ticket
			for _, ticketID := range match.Tickets {
				mapMatches[ticketID] = match.MatchID
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	// assert that all players are in a match
	i := 0
	for range mapMatches {
		i++
	}
	suite.Require().GreaterOrEqual(i, maxPlayers)

	for _, v := range mapPlayers {
		if _, ok := mapMatches[v]; ok {
			suite.Require().True(ok)
		}
	}
	// TODO confirm the match to clean up tickets (or they will expire later)

}
func TestE2E(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func findMatchRequestWithScore(score int) (response string, err error) {
	return findMatchRequestWithInput(handler.FindMatchInput{
		PlayerID: ksuid.New().String(),
		Score:    score,
	})
}

func findMatchRequest() (response string, err error) {
	return findMatchRequestWithInput(handler.FindMatchInput{
		PlayerID: ksuid.New().String(),
		Score:    1,
	})
}
func findMatchRequestWithInput(input handler.FindMatchInput) (response string, err error) {
	path := fmt.Sprintf("/api/v1/queue/global")

	in, err := json.Marshal(&input)

	if err != nil {
		return "", err
	}

	err = requests.
		URL(path).
		Put().
		BodyReader(bytes.NewReader(in)).
		Host(baseURL).
		ContentType("application/json").
		Scheme(defaultScheme).
		CheckStatus(http.StatusOK).
		ToString(&response).
		Fetch(context.Background())
	return response, err
}

func getMatchRequest(ticketID string, expectedStatus int) (response string, err error) {
	path := fmt.Sprintf("/api/v1/ticket/%s", ticketID)
	err = requests.
		URL(path).
		Host(baseURL).
		Scheme(defaultScheme).
		CheckStatus(expectedStatus).
		ToString(&response).
		Fetch(context.Background())

	return response, err
}
func getMatchRequestAny(ticketID string) (response string, err error) {
	path := fmt.Sprintf("/api/v1/ticket/%s", ticketID)
	err = requests.
		URL(path).
		Host(baseURL).
		Scheme(defaultScheme).
		ToString(&response).
		Fetch(context.Background())
	return response, err
}

func cancelMatchRequest(ticketID string) (err error) {
	path := fmt.Sprintf("/api/v1/ticket/%s", ticketID)
	err = requests.
		URL(path).
		Delete().
		Host(baseURL).
		Scheme(defaultScheme).
		CheckStatus(http.StatusNoContent).
		Fetch(context.Background())
	return err
}
