// go:build integration

package integration_test

import (
	"context"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/api/eventservice"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SystemStatsSuite struct {
	suite.Suite
	statConn   *grpc.ClientConn
	statClient eventservice.CalendarClient
	ctx        context.Context
	respEmpty  *eventservice.GetEventsResponse
}

func (s *SystemStatsSuite) SetupSuite() {
	config.InitCalendarSettings()
	cfg := config.CalendarSettings.Server

	statConn, err := grpc.Dial(net.JoinHostPort(cfg.Host, cfg.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.statConn = statConn
	s.statClient = eventservice.NewCalendarClient(s.statConn)

	// s.ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	ctx := context.Background()
	s.ctx = ctx

	respEmpty := eventservice.GetEventsResponse{}
	s.respEmpty = &respEmpty
}

func (s *SystemStatsSuite) SetupTest() {
	seed := time.Now().UnixNano()
	rand.NewSource(seed)
	s.T().Log("seed:", seed)
}

func (s *SystemStatsSuite) TearDownSuite() {
	err := s.statConn.Close()
	if err != nil {
		s.T().Log(err)
	}
}

func (s *SystemStatsSuite) TearDownTest() {
	s.T().Logf("%s - done", s.T().Name())
}

func (s *SystemStatsSuite) TestStandard() {
}

func TestSystemStatusSuite(t *testing.T) {
	suite.Run(t, new(SystemStatsSuite))
}
