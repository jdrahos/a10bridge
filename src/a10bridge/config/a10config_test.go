package config_test

import (
	"a10bridge/config"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
)

type A10ConfigTestSuite struct {
	suite.Suite
}

func TestA10Config(t *testing.T) {
	tests := new(A10ConfigTestSuite)
	suite.Run(t, tests)
}

func (suite *A10ConfigTestSuite) TestSort() {
	a10Instances := config.A10Instances{
		config.A10Instance{
			Name: "Z",
		},
		config.A10Instance{
			Name: "a",
		},
		config.A10Instance{
			Name: "z",
		},
		config.A10Instance{
			Name: "A",
		},
	}
	sort.Sort(a10Instances)
	suite.Assert().Equal("A", a10Instances[0].Name)
	suite.Assert().Equal("Z", a10Instances[1].Name)
	suite.Assert().Equal("a", a10Instances[2].Name)
	suite.Assert().Equal("z", a10Instances[3].Name)
}
