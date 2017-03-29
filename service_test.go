package service

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	s   Service
	url string
}

func (s *ServiceSuite) SetupSuite() {
	httpmock.Activate()
	s.url = "http://localhost:8080"
	s.s = New()
}

func (s *ServiceSuite) BeforeTest() {
	httpmock.Reset()
}

func (s *ServiceSuite) AfterTest() {

}

func (s *ServiceSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func TestRunServiceSuite(t *testing.T) {
	s := new(ServiceSuite)
	suite.Run(t, s)
}

func (s *ServiceSuite) TestGetSuccess() {
	assert := assert.New(s.T())

	data := s.SuccessResponse()
	temp := make([]string, 1)
	temp[0] = "one"
	data.Data = temp
	res, _ := httpmock.NewJsonResponder(200, data)
	req := &Request{
		Method: "GET",
		Url:    s.url,
	}

	httpmock.RegisterResponder("GET", s.url, res)

	var body []string
	err := s.s.Send(req, &body)

	assert.NoError(err)
	assert.NotNil(body)
	assert.Equal(1, len(body))
	assert.EqualValues("one", body[0])
}

func (s *ServiceSuite) TestPostSuccess() {
	assert := assert.New(s.T())

	data := s.SuccessResponse()
	res, _ := httpmock.NewJsonResponder(200, data)
	req := &Request{
		Method: "GET",
		Url:    s.url,
		Data:   s.SuccessResponse(),
	}

	httpmock.RegisterResponder("POST", s.url, res)

	err := s.s.Send(req, nil)

	assert.NoError(err)
}

func (s *ServiceSuite) TestGetError() {
	assert := assert.New(s.T())

	data := s.ErrorResponse("ERROR")
	res, _ := httpmock.NewJsonResponder(500, data)
	req := &Request{
		Method: "GET",
		Url:    s.url,
	}

	httpmock.RegisterResponder("GET", s.url, res)

	var i interface{}
	err := s.s.Send(req, i)

	assert.Error(err)
}

func (s *ServiceSuite) TestGetJsonParseError() {
	assert := assert.New(s.T())

	data := "{{}"
	res, _ := httpmock.NewJsonResponder(500, data)
	req := &Request{
		Method: "GET",
		Url:    s.url,
	}

	httpmock.RegisterResponder("GET", s.url, res)

	var i interface{}
	err := s.s.Send(req, i)

	assert.Error(err)
}

func (s *ServiceSuite) TestGetInvalidJSON() {
	assert := assert.New(s.T())

	data := s.SuccessResponse()
	data.Data = "{{]"
	res, _ := httpmock.NewJsonResponder(200, data)
	req := &Request{
		Method: "GET",
		Url:    s.url,
	}

	httpmock.RegisterResponder("GET", s.url, res)

	var body []string
	err := s.s.Send(req, &body)

	assert.Error(err)
	assert.Nil(body)
}

func (s *ServiceSuite) EmptyResponse() *Response {
	return &Response{}
}

func (s *ServiceSuite) SuccessResponse() *Response {
	r := s.EmptyResponse()
	r.Success = true
	return r
}

func (s *ServiceSuite) ErrorResponse(msg string) *Response {
	r := s.EmptyResponse()
	r.Success = false
	r.Msg = msg
	return r
}
