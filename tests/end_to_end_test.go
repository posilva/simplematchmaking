package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"
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
	err := os.Setenv("MATCHMAKING_CFG", "json.eyJxdWV1ZXMiOnsiZ2xvYmFsIjp7Im5hbWUiOiJnbG9iYWwiLCJtYXhQbGF5ZXJzIjoyLCJuckJyYWNrZXRzIjoxMDAsIm1heFJhbmtpbmciOjEwMDAsIm1pblJhbmtpbmciOjEsIm1ha2VJdGVyYXRpb25zIjozfX0sIm1hdGNobWFrZXJzIjp7Imdsb2JhbCI6eyJuYW1lIjoiZ2xvYmFsIiwiaW50ZXJ2YWxTZWNzIjo1LCJtYWtlVGltZW91dFNlY3MiOjN9fX0=")
	suite.Require().NoError(err)
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

type PlayerGenerator struct {
	mx         sync.Mutex
	MapPlayers map[string]string
	MapTickets map[string]string
	MapMatches map[string]string
	Matches    map[string]struct{}
}

func NewPlayerGenerator() *PlayerGenerator {
	return &PlayerGenerator{
		MapPlayers: make(map[string]string),
		MapTickets: make(map[string]string),
		MapMatches: make(map[string]string),
		Matches:    make(map[string]struct{}),
	}
}

func (pg *PlayerGenerator) AddPlayer(playerID, ticketID string) {
	pg.mx.Lock()
	defer pg.mx.Unlock()
	pg.MapPlayers[playerID] = ticketID
	pg.MapTickets[ticketID] = playerID
}

func (pg *PlayerGenerator) GeneratePlayers(maxPlayers int, min int, max int) {
	start := time.Now()
	defer func() {
		fmt.Println("Time to generate players: ", time.Since(start))
	}()
	var wg sync.WaitGroup
	wg.Add(maxPlayers)
	for i := 0; i < maxPlayers; i++ {
		go func() {
			playerID := ksuid.New().String()
			score := rand.Intn(max-min) + min
			response, err := findMatchRequestWithInput(handler.FindMatchInput{
				PlayerID: playerID,
				Score:    score,
			})
			if err != nil {
				fmt.Println("Error: ", err)
			}
			var out handler.FindMatchOutput
			err = json.Unmarshal([]byte(response), &out)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			pg.AddPlayer(playerID, out.TicketID)
			wg.Done()
		}()
	}
	wg.Wait()
}

func (pg *PlayerGenerator) CheckTickets(suite *E2ETestSuite) {
	for _, ticketID := range pg.MapPlayers {
		if _, ok := pg.MapMatches[ticketID]; ok {
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
		pg.Matches[match.MatchID] = struct{}{}
		// save the matchID for each ticket
		for _, ticketID := range match.Tickets {
			pg.MapMatches[ticketID] = match.MatchID
		}
	}
}

func (suite *E2ETestSuite) TestMatchmaking() {
	maxPlayers := 1000
	playersPerMatch := 2
	min := 0
	max := 1000

	pg := NewPlayerGenerator()
	pg.GeneratePlayers(maxPlayers, min, max)
	suite.Require().Equal(len(pg.MapPlayers), maxPlayers)
	fmt.Println("Players generated: ", len(pg.MapPlayers))

	for len(pg.Matches) < maxPlayers/playersPerMatch {
		pg.CheckTickets(suite)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Matches found: ", len(pg.Matches))

	// assert that all players are in a match
	i := len(pg.MapMatches)
	suite.Require().GreaterOrEqual(i, maxPlayers)

	for _, v := range pg.MapPlayers {
		if _, ok := pg.MapMatches[v]; ok {
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
