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

func (s *ServiceSuite) TestNewCustom() {
	assert := assert.New(s.T())

	service := NewCustom(&defaultResponder{})

	data := s.SuccessResponse(nil)
	res, _ := httpmock.NewJsonResponder(200, data)
	req := &Request{
		Method: "GET",
		URL:    s.url,
	}

	httpmock.RegisterResponder("GET", s.url, res)

	var body []string
	err := service.Send(req, &body)

	assert.NoError(err)
}

func (s *ServiceSuite) TestGetSuccess() {
	assert := assert.New(s.T())

	data := s.SuccessResponse(nil)
	res, _ := httpmock.NewJsonResponder(200, data)
	req := &Request{
		Method: "GET",
		URL:    s.url,
	}

	httpmock.RegisterResponder("GET", s.url, res)

	var body []string
	err := s.s.Send(req, &body)

	assert.NoError(err)
}

func (s *ServiceSuite) TestPostSuccess() {
	assert := assert.New(s.T())

	data := s.SuccessResponse(nil)
	res, _ := httpmock.NewJsonResponder(200, data)
	req := &Request{
		Method: "POST",
		URL:    s.url,
		Data:   s.SuccessResponse(nil),
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
		URL:    s.url,
	}

	httpmock.RegisterResponder("GET", s.url, res)

	var i interface{}
	err := s.s.Send(req, i)

	assert.Error(err)
}

func (s *ServiceSuite) TestGetUnknownError() {
	assert := assert.New(s.T())

	data := s.ErrorResponse("")
	res, _ := httpmock.NewJsonResponder(500, data)
	req := &Request{
		Method: "GET",
		URL:    s.url,
	}

	httpmock.RegisterResponder("GET", s.url, res)

	var i interface{}
	err := s.s.Send(req, i)

	assert.Error(err)
}

func (s *ServiceSuite) TestGetJsonParseError() {
	assert := assert.New(s.T())

	data := "{"
	res, _ := httpmock.NewJsonResponder(500, data)
	req := &Request{
		Method: "GET",
		URL:    s.url,
	}

	httpmock.RegisterResponder("GET", s.url, res)

	var i interface{}
	err := s.s.Send(req, i)

	assert.Error(err)
}

func (s *ServiceSuite) TestGetInvalidJSON() {
	assert := assert.New(s.T())

	data := s.SuccessResponse("{{]")
	res, _ := httpmock.NewJsonResponder(200, data)
	req := &Request{
		Method: "GET",
		URL:    s.url,
	}

	httpmock.RegisterResponder("GET", s.url, res)

	var body []string
	err := s.s.Send(req, &body)

	assert.Error(err)
	assert.Nil(body)
}

func (s *ServiceSuite) TestGetInvalidRequest() {
	assert := assert.New(s.T())

	data := s.SuccessResponse("{{]")
	res, _ := httpmock.NewJsonResponder(200, data)
	req := &Request{
		Method: "INVALID_METHOD",
		URL:    s.url,
	}

	httpmock.RegisterResponder("GET", s.url, res)

	var body []string
	err := s.s.Send(req, &body)

	assert.Error(err)
	assert.Nil(body)
}

func (s *ServiceSuite) EmptyResponse() *defaultResponse {
	return &defaultResponse{}
}

func (s *ServiceSuite) SuccessResponse(data interface{}) Response {
	r := s.EmptyResponse()
	r.Success = true
	r.Data = data
	return r
}

func (s *ServiceSuite) ErrorResponse(err string) Response {
	r := s.EmptyResponse()
	r.Success = false
	r.Msg = err
	return r
}
