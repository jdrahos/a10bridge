package util_test

import (
	"a10bridge/util"
	"testing"

	"github.com/stretchr/testify/suite"
)

type StringUtilsTestSuite struct {
	suite.Suite
	helper *util.TestHelper
}

func TestStringUtils(t *testing.T) {
	tests := new(StringUtilsTestSuite)
	tests.helper = new(util.TestHelper)
	suite.Run(t, tests)
}

func (suite *StringUtilsTestSuite) TestContains() {
	stringSlice := []string{
		"test 1",
		"test 2",
		"test 3",
	}
	contains := util.Contains(stringSlice, "test 1")
	suite.Assert().True(contains)

	contains = util.Contains(stringSlice, "test 4")
	suite.Assert().False(contains)
}

func (suite *StringUtilsTestSuite) TestApplyTemplate() {
	entity := struct {
		Name   string
		Number int
	}{
		Name:   "world",
		Number: 5,
	}
	result, err := util.ApplyTemplate(entity, "Hello {{.Name}}, it is {{.Number}} o'clock!")
	suite.Assert().Nil(err)
	suite.Assert().Equal("Hello world, it is 5 o'clock!", result)
}

func (suite *StringUtilsTestSuite) TestApplyTemplate_wrongTemplate() {
	entity := struct {
		Name   string
		Number int
	}{
		Name:   "world",
		Number: 5,
	}
	_, err := util.ApplyTemplate(entity, "Hello {{.Name}, it is {{.Number}} o'clock!")
	suite.Assert().NotNil(err)
}

func (suite *StringUtilsTestSuite) TestApplyTemplate_wrongEntity() {
	entity := struct {
		Name   string
		Number int
	}{
		Name:   "world",
		Number: 5,
	}
	_, err := util.ApplyTemplate(entity, "Hello {{.IDontExist}}, it is {{.Number}} o'clock!")
	suite.Assert().NotNil(err)
}

func (suite *StringUtilsTestSuite) TestToJSON() {
	entity := struct {
		Name   string `json:"name"`
		Number int    `json:"num"`
	}{
		Name:   "world",
		Number: 5,
	}
	result := util.ToJSON(entity)
	suite.Assert().Equal("{\n  \"name\": \"world\",\n  \"num\": 5\n}", result)
}

func (suite *StringUtilsTestSuite) TestToJSON_marshallingFails() {
	entity := struct {
		Name   chan string `json:"name"`
		Number int         `json:"num"`
	}{
		Name:   make(chan string),
		Number: 5,
	}
	result := util.ToJSON(entity)
	suite.Assert().Equal("json: unsupported type: chan string", result)
}
