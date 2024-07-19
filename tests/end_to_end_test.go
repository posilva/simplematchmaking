package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/carlmjohnson/requests"

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
	fmt.Println("Running teardown suite")
	teardown(&suite.BaseTestSuite)
}

func (suite *E2ETestSuite) TestFindMatch() {
	resp, err := findMatchRequest()
	suite.NoError(err)
	var out handler.FindMatchOutput
	err = json.Unmarshal([]byte(resp), &out)
	suite.NoError(err)
	suite.Regexp(regexp.MustCompile("^[a-zA-Z0-9]{27}$"), out.TicketID)
}

func (suite *E2ETestSuite) TestGetMatch() {
	resp, err := getMatchRequest("ticket1")
	suite.NoError(err)
	var out handler.GetMatchOutput
	err = json.Unmarshal([]byte(resp), &out)
	suite.NoError(err)
	suite.Equal("match1", out.MatchID)
}

func (suite *E2ETestSuite) TestCancelMatch() {
	err := cancelMatchRequest("ticket1")
	suite.NoError(err)
}

func TestE2E(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func findMatchRequest() (response string, err error) {
	path := fmt.Sprintf("/api/v1/queue")

	in, err := json.Marshal(&handler.FindMatchInput{
		PlayerID: "1",
		Score:    1,
	})

	if err != nil {
		return "", err
	}

	err = requests.
		URL(path).
		Put().
		BodyReader(bytes.NewReader(in)).
		Host(baseURL).
		Scheme(defaultScheme).
		CheckStatus(http.StatusOK).
		ToString(&response).
		Fetch(context.Background())
	return response, err
}

func getMatchRequest(ticketID string) (response string, err error) {
	path := fmt.Sprintf("/api/v1/queue/%s", ticketID)

	err = requests.
		URL(path).
		Host(baseURL).
		Scheme(defaultScheme).
		CheckStatus(http.StatusOK).
		ToString(&response).
		Fetch(context.Background())
	return response, err
}

func cancelMatchRequest(ticketID string) (err error) {
	path := fmt.Sprintf("/api/v1/queue/%s", ticketID)

	err = requests.
		URL(path).
		Delete().
		Host(baseURL).
		Scheme(defaultScheme).
		CheckStatus(http.StatusOK).
		Fetch(context.Background())
	return err
}
