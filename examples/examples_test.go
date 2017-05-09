package examples

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

type mockResponse struct {
	Msg     string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success,omitempty"`
}

func TestExampleGetSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert := assert.New(t)
	s := NewMyService()

	b := &User{Name: "test", Email: "someemail"}
	res, _ := httpmock.NewJsonResponder(200, mockResponse{Success: true, Data: b})

	httpmock.RegisterResponder("GET", "http://some.service.com/user", res)

	u, err := s.Get("test")

	assert.NoError(err)
	assert.NotNil(u)
	assert.EqualValues("test", u.Name)
}

func TestExampleGetError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert := assert.New(t)
	s := NewMyService()

	res, _ := httpmock.NewJsonResponder(500, &mockResponse{Success: false, Msg: "error"})

	httpmock.RegisterResponder("GET", "http://some.service.com/user", res)

	u, err := s.Get("Test")

	assert.Error(err)
	assert.Nil(u)
}
