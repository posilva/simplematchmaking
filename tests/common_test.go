package tests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/google/uuid"
	"github.com/phayes/freeport"
	"github.com/posilva/simplematchmaking/cmd/simplematchmaking/app"
	"github.com/posilva/simplematchmaking/cmd/simplematchmaking/config"
	"github.com/stretchr/testify/require"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/redis/rueidis"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
	testcontainersredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	uniqueTestID             = uuid.NewString()
	defaultLbName            = "integration_lb_tests::" + uniqueTestID
	defaultLbNameSumMultiple = defaultLbName + "::Sum::multiple"
	defaultLbNameSum         = defaultLbName + "::Sum"
	defaultLbNameMax         = defaultLbName + "::Max"
	defaultLbNameMin         = defaultLbName + "::Min"
	defaultLbNameLast        = defaultLbName + "::Last"
	metadataDefault          = map[string]string{
		"country": "PT",
		"league":  "gold",
	}
)

type BaseTestSuite struct {
	suite.Suite
	Context            context.Context
	RedisContainer     *testcontainersredis.RedisContainer
	RedisClient        rueidis.Client
	DDBContainer       *localstack.LocalStackContainer
	AWSConfig          aws.Config
	DDBClient          *dynamodb.Client
	RedisEndpoint      string
	LocalstackEndpoint string
	ServiceEndpoint    string
}

func waitForService(suite *BaseTestSuite) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	c := &http.Client{Timeout: time.Second * 10}
	for {
		r, err := c.Get(fmt.Sprintf("http://%s/", suite.ServiceEndpoint))
		if err == nil && r.StatusCode == http.StatusOK {
			return
		}
		select {
		case <-time.After(time.Millisecond * 100):
		case <-ctx.Done():
			require.Fail(suite.T(), "timeout waiting for health check")
		}
	}
}

func setup(suite *BaseTestSuite) {
	log.Println("Running setup suite")

	suite.Context = context.Background()
	testcontainers.Logger = log.New(&ioutils.NopWriter{}, "", 0)
	setupRedisContainer(suite)

	port, err := freeport.GetFreePort()
	if err != nil {
		panic("failed to get free port: " + err.Error())
	}
	remoteAddress, remote := os.LookupEnv("TEST_REMOTE_ADDR")
	if remote {
		fmt.Println("Remote address: ", remoteAddress)
	}
	suite.ServiceEndpoint = fmt.Sprintf("127.0.0.1:%d", port)
	if remote {
		// this allows to run end 2 end tests against a remote endpoint
		suite.ServiceEndpoint = remoteAddress
	}

	log.Println("Service endpoint: ", suite.ServiceEndpoint)
	if !remote {
		go func() {
			config.SetAddr(suite.ServiceEndpoint)
			config.SetRedisAddr(suite.RedisEndpoint)
			config.SetLocal(true)
			app.Run()
		}()
		waitForService(suite)
		log.Printf("Service is running on %s", suite.ServiceEndpoint)
	}
}

func setupRedisContainer(suite *BaseTestSuite) {
	redisContainer, err := testcontainersredis.RunContainer(
		suite.Context,
		testcontainers.WithImage("redis:latest"),
		testcontainers.WithWaitStrategyAndDeadline(
			30*time.Second, wait.ForExposedPort()),
	)
	suite.NoError(err)

	ip, err := redisContainer.Host(suite.Context)
	suite.NoError(err)
	port, err := redisContainer.MappedPort(suite.Context, "6379")
	suite.NoError(err)

	endpoint := fmt.Sprintf("%s:%s", ip, port.Port())
	log.Printf("Redis endpoint: %s", endpoint)
	redisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{endpoint},
	})

	suite.NoError(err)
	suite.RedisEndpoint = endpoint
	pingCmd := redisClient.B().Ping().Build()
	err = redisClient.Do(suite.Context, pingCmd).Error()
	suite.NoError(err)

	suite.RedisContainer = redisContainer
	suite.RedisClient = redisClient
}

func teardown(suite *BaseTestSuite) {
	err := suite.RedisContainer.Terminate(suite.Context)
	suite.NoError(err)
	// err = suite.RedisClient.Do(suite.Context, suite.RedisClient.B().Flushall().Build()).Error()
	// suite.NoError(err)
}

/**
TEMPLATE OF INTERGRATION TEST
package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ExampleTestSuite struct {
	BaseTestSuite
}

func (suite *ExampleTestSuite) SetupSuite() {
	setup(&suite.BaseTestSuite)
	fmt.Println("Running setup suite")
}

func (suite *ExampleTestSuite) TearDownSuite() {
	fmt.Println("Running teardown suite")
	teardown(&suite.BaseTestSuite)
}

func (suite *ExampleTestSuite) SetupTest() {
	fmt.Println("Running setup test")
}

func (suite *ExampleTestSuite) TearDownTest() {
	fmt.Println("Running teardown test")
	err := suite.RedisClient.Do(suite.Context, suite.RedisClient.B().Flushall().Build()).Error()
	suite.NoError(err)
}

func (suite *ExampleTestSuite) BeforeTest(_ string, testName string) {
	fmt.Printf("Running before test: %s\n", testName)
}

func (suite *ExampleTestSuite) AfterTest(_ string, testName string) {
	fmt.Printf("Running after test: %s\n", testName)
}

func (suite *ExampleTestSuite) TestExample() {
}

func TestExample(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}

*/
