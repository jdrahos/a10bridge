package util_test

import (
	"a10bridge/util"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NetUtilsTestSuite struct {
	suite.Suite
}

func TestNetUtils(t *testing.T) {
	tests := new(NetUtilsTestSuite)
	suite.Run(t, tests)
}

func (suite *NetUtilsTestSuite) TestLookupLocalhost() {
	localIp, err := util.LookupIP("localhost")
	suite.Assert().Nil(err)

	if strings.Contains(localIp, ":") {
		suite.Assert().Equal("::1", localIp)
	} else {
		suite.Assert().Equal("127.0.0.1", localIp)
	}
}

func (suite *NetUtilsTestSuite) TestLookupFailure() {
	_, err := util.LookupIP("i.m.definitelly.not.a.hosname.of.any.kind")
	suite.Assert().NotNil(err)
}

func (suite *NetUtilsTestSuite) TestResolverInjection() {
	_, err := util.LookupIP("i.m.definitelly.not.a.hosname.of.any.kind")
	suite.Assert().NotNil(err)

	original := util.InjectIPResolver(testResolver{})
	defer util.InjectIPResolver(original)

	ip, err := util.LookupIP("i.m.definitelly.not.a.hosname.of.any.kind")
	suite.Assert().Nil(err)
	suite.Assert().Equal("10.10.10.10", ip)
}

type testResolver struct{}

func (resolver testResolver) LookupIP(hostname string) (string, error) {
	return "10.10.10.10", nil
}
