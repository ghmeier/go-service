# Go Service
HTTP service gateway in go.

A simple library to send requests to external services without marshaling JSON yourself.

## Installation

```
go get github.com/ghmeier/go-service
```

## Examples

Sending one `GET` request:
```
package example

import "github.com/ghmeier/go-service"

type Foo struct {
    name string `json:"name"`
}

func Get() {
    s := service.New()

    var f Foo
    err := s.Send(&service.Request{
        Method: "GET",
        URL:    "http://some.url.com/foo",
    }, &f)
}

```
The default response will fill `f` with the object in the `data` field of the response and return an error if Response.`success` is `false`.


Sending one `POST` request:
```
package example

import "github.com/ghmeier/go-service"

func Post() {
    s := service.New()

    f := &Foo{}
    err := s.Send(&service.Request{
        Method: "POST",
        URL:    "http://some.url.com/foo",
        Data:   f,
    }, nil)
}

```

Check out `examples.go` for an example of implementing a service gateway using `go-service`.

## Custom Response

If you want to handle responses from services with different repsponse types, implement the `service.Responder` and `service.Response` interfaces.

### Example Custom Responder/Response:
```
type CustomResponder struct{}

//Marshal should return a new response based on an http response. Usually this Unmarshals the body.
func (c *CustomResponder) Marshal(r *http.Response) (service.Response, error) {
    res := &CustomResponse{}

    res.OK = r.StatusCode == 200
    res.Status = r.Status

    return res
}

type CustomResponse struct {
    Status string `json:"status"`
    OK     bool   `json:"ok"`
}

//Error should return an error based on the CustomResponse state
func (c *CustomResponse) Error() error {
    if !c.OK {
        return fmt.Errorf(c.Status)
    }

    return nil
}

//Body should return a []byte that can be marshalled into the second argument of Service.Send()
func (c *CustomResponse) Body() ([]byte, error) {
    return json.Marshal(c)
}
```

Using Custom Responder:
```
s := service.NewCustom(&CustomResponder{})
```
