package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	h "net/http"
)

/*Responder marshals an http.Response into a service.Response that will
  return either a Body() or Error()*/
type Responder interface {
	Marshal(*h.Response) (Response, error)
}

/*Response contains relevant data from http server responses
It can be used in conjunction with the following methods like this:

data, _ := ServiceGet(url)
var c []*models.Content
err := json.Unmarshal(data, &c)

Since ServiceGet returns data as a []byte, we can unmarshal it
to whatever is needed in the calling method. Here, its []*models.Content
*/
type Response interface {
	Error() error
	Body() ([]byte, error)
}

/*Request contains the Method used in sending, the Url to request, and
  any Data to be sent*/
type Request struct {
	Method string
	URL    string
	Data   interface{}
}

/*Service has a method that will send a Request and put the response into
  the provided interface using json unmarshalling. The interface must be a
  pointer type */
type Service interface {
	Send(*Request, interface{}) error
}

type defaultResponder struct{}

type defaultResponse struct {
	Msg     string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success,omitempty"`
}

type http struct {
	client    *h.Client
	responder Responder
}

/*New returns a Service with the default http Client*/
func New() Service {
	return &http{
		client:    &h.Client{},
		responder: &defaultResponder{},
	}
}

/*NewCustom returns a BaseService with a custom responder*/
func NewCustom(r Responder) Service {
	return &http{
		client:    &h.Client{},
		responder: r,
	}
}

/*ServiceSend sends a request of type METHOD to the url with data as the
  JSON payload and puts the response into i*/
func (b *http) Send(req *Request, i interface{}) error {
	var r *bytes.Buffer
	var err error
	if req.Data != nil {
		b, err := json.Marshal(req.Data)
		if err != nil {
			return err
		}
		r = bytes.NewBuffer(b)
	} else {
		r = nil
	}

	var prepared *h.Request
	if r != nil {
		prepared, err = h.NewRequest(req.Method, req.URL, r)
	} else {
		prepared, err = h.NewRequest(req.Method, req.URL, nil)
	}

	if err != nil {
		return err
	}

	raw, err := b.do(prepared)
	if err != nil {
		return err
	}

	if i == nil {
		return nil
	}

	err = json.Unmarshal(raw, i)
	if err != nil {
		return err
	}

	return nil
}

func (b *http) do(req *h.Request) ([]byte, error) {
	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rData, err := b.handle(resp)
	if err != nil {
		return nil, err
	}

	return rData, nil
}

func (b *http) handle(resp *h.Response) ([]byte, error) {
	response, err := b.responder.Marshal(resp)
	if err != nil {
		return nil, err
	}

	err = response.Error()
	if err != nil {
		return nil, err
	}

	return response.Body()
}

func (r *defaultResponder) Marshal(res *h.Response) (Response, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response defaultResponse
	fmt.Println(string(body))
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *defaultResponse) Error() error {
	if !r.Success {
		if r.Msg != "" {
			return fmt.Errorf("%s", r.Msg)
		}

		return fmt.Errorf("ERROR: unknown error")
	}
	return nil
}

func (r *defaultResponse) Body() ([]byte, error) {
	return json.Marshal(r.Data)
}
