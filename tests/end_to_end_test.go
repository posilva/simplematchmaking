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
	"github.com/segmentio/ksuid"

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

func TestE2E(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func findMatchRequest() (response string, err error) {
	path := fmt.Sprintf("/api/v1/queue")

	in, err := json.Marshal(&handler.FindMatchInput{
		PlayerID: ksuid.New().String(),
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
		ContentType("application/json").
		Scheme(defaultScheme).
		CheckStatus(http.StatusOK).
		ToString(&response).
		Fetch(context.Background())
	return response, err
}

func getMatchRequest(ticketID string, expectedStatus int) (response string, err error) {
	path := fmt.Sprintf("/api/v1/queue/%s", ticketID)

	err = requests.
		URL(path).
		Host(baseURL).
		Scheme(defaultScheme).
		CheckStatus(expectedStatus).
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
		CheckStatus(http.StatusNoContent).
		Fetch(context.Background())
	return err
}
