package meiporo

import (
	"reflect"
	"testing"
)

func TestRouter(t *testing.T) {
	r := NewRouter()
	r.Add("GET", "/:id", func(c *Context) {})
	r.Add("GET", "/", func(c *Context) {})
	r.Start()
	func() {
		route, params := r.Match("GET", "/test")
		if !reflect.DeepEqual(route, &r.Routes["GET"][0]) {
			t.Errorf("Expected %v, but %v", &r.Routes["GET"][0], route)
		}
		if !reflect.DeepEqual(params, map[string]string{"id": "test"}) {
			t.Errorf("Expected %v, but %v", map[string]string{"id": "test"}, params)
		}
	}()
	func() {
		route, params := r.Match("GET", "/")
		if !reflect.DeepEqual(route, &r.Routes["GET"][1]) {
			t.Errorf("Expected %v, but %v", &r.Routes["GET"][1], route)
		}
		if !reflect.DeepEqual(params, map[string]string{}) {
			t.Errorf("Expected %v, but %v", map[string]string{}, params)
		}
	}()
}
