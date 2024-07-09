package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/carlmjohnson/requests"
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

func (suite *E2ETestSuite) TestSimple() {
	resp, err := makeRequest()
	suite.NoError(err)
	suite.Equal("{\"id\":\"ticket1\"}", resp)
}

func TestE2E(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func makeRequest() (response string, err error) {
	path := fmt.Sprintf("/api/v1/mm")

	err = requests.
		URL(path).
		Host(baseURL).
		Scheme(defaultScheme).
		CheckStatus(http.StatusOK).
		ToString(&response).
		Fetch(context.Background())
	return response, err
}
