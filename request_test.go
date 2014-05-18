package meiporo

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestRequest(t *testing.T) {
	func() {
		r, _ := http.NewRequest("GET", "http://example.com?f=foo&b=bar", nil)
		request := NewRequest(r)
		request.ParseForm()
		params := make(url.Values)
		params.Set("f", "foo")
		params.Set("b", "bar")
		if !reflect.DeepEqual(request.Form, params) {
			t.Errorf("Expected %v, but %v", request.Form, params)
		}
	}()
}
