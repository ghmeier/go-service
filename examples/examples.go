package examples

import (
	"github.com/ghmeier/go-service"
)

/*MyService is your http gateway to call external APIs that
  has the default Service embedded*/
type MyService struct {
	service.Service
}

//User is a custom data type with JSON bindings
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

/*NewMyService initializes an instance of MyService with the embedded
  default service*/
func NewMyService() *MyService {
	return &MyService{Service: service.New("http://some.service.com").Copy("user")}
}

/*Get sends a request using the service to handle preparing and
  unmarshalling the request*/
func (m *MyService) Get(name string) (*User, error) {
	// the value to be unmashalled
	var u User
	err := m.Send(&service.Request{
		// run a GET request
		Method: "GET",
	}, &u)

	if err != nil {
		return nil, err
	}
	return &u, nil
}
