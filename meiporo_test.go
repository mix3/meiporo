package meiporo

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

var result = ""

type FirstMiddleware struct{}

func (m *FirstMiddleware) Handler() Handler {
	return func(c *Context) {
		result += "foo"
		c.Next()
		result += "bah"
	}
}

type SecondMiddleware struct{}

func (m *SecondMiddleware) Handler() Handler {
	return func(c *Context) {
		result += "bar"
		c.Next()
		result += "baz"
	}
}

func TestMeiporo(t *testing.T) {
	func() {
		m := New()
		if m == nil {
			t.Error("meiporo.New() cannot return nil")
		}
	}()
	func() {
		go New().Run()
	}()
	func() {
		result = ""
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		m := New()
		m.Use(&FirstMiddleware{})
		m.Use(&SecondMiddleware{})
		m.Router.Add("GET", "/", func(c *Context) {
			result += "bat"
			c.Res.WriteHeader(http.StatusBadRequest)
		})
		m.Router.Start()
		m.ServeHTTP(rec, req)
		expect(t, result, "foobarbatbazbah")
		expect(t, rec.Code, http.StatusBadRequest)
	}()
	func() {
		result = ""
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		m := New()
		m.Use(&FirstMiddleware{})
		m.Router.Add("GET", "/",
			func(c *Context) {
				result += "batman!"
				c.Next()
			},
			func(c *Context) {
				result += "batman!"
				c.Next()
			},
			func(c *Context) {
				result += "batman!"
				c.Next()
			},
			func(c *Context) {
				result += "bat"
				c.Res.WriteHeader(http.StatusBadRequest)
			},
		)
		m.Router.Start()
		m.ServeHTTP(rec, req)
		expect(t, result, "foobatman!batman!batman!batbah")
		expect(t, rec.Code, http.StatusBadRequest)
	}()
	func() {
		result = ""
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		m := New()
		m.Router.Add("GET", "/", func(c *Context) {
			c.Res.WriteHeader(http.StatusOK)
		})
		m.Router.Start()
		c := m.CreateContext(rec, req)
		expect(t, c.Res.Written(), false)
		c.Run()
		expect(t, c.Res.Written(), true)
	}()
}
