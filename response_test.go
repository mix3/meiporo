package meiporo

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lestrrat/go-xslate"
)

func TestResponse(t *testing.T) {
	func() {
		rec := httptest.NewRecorder()
		rw := NewResponseWriter(rec)
		rw.Write([]byte("Hello world"))
		expect(t, rec.Code, rw.Status())
		expect(t, rec.Body.String(), "Hello world")
		expect(t, rw.Status(), http.StatusOK)
		expect(t, rw.Size(), 11)
		expect(t, rw.Written(), true)
	}()
	func() {
		rec := httptest.NewRecorder()
		rw := NewResponseWriter(rec)

		rw.Write([]byte("Hello world"))
		rw.Write([]byte("foo bar bat baz"))

		expect(t, rec.Code, rw.Status())
		expect(t, rec.Body.String(), "Hello worldfoo bar bat baz")
		expect(t, rw.Status(), http.StatusOK)
		expect(t, rw.Size(), 26)
		expect(t, rw.Written(), true)
	}()
	func() {
		rec := httptest.NewRecorder()
		rw := NewResponseWriter(rec)

		rw.WriteHeader(http.StatusNotFound)

		expect(t, rec.Code, rw.Status())
		expect(t, rec.Body.String(), "")
		expect(t, rw.Status(), http.StatusNotFound)
		expect(t, rw.Size(), 0)
		expect(t, rw.Written(), true)
	}()
	func() {
		rec := httptest.NewRecorder()
		rw := NewResponseWriter(rec)
		data := map[string]string{
			"f": "foo",
			"b": "bar",
		}
		rw.WriteJSON(data)
		expect_body := `{"b":"bar","f":"foo"}
`
		expect(t, rec.Code, rw.Status())
		expect(t, rec.Body.String(), expect_body)
		expect(t, rw.Status(), http.StatusOK)
		expect(t, rw.Size(), len(expect_body))
		expect(t, rw.Written(), true)
	}()
	func() {
		rec := httptest.NewRecorder()
		rw := NewResponseWriter(rec)
		type Sample struct {
			F string `xml:F`
			B string `xml:B`
		}
		data := Sample{
			F: "foo",
			B: "bar",
		}
		rw.WriteXML(data)
		expect_body := `<Sample><F>foo</F><B>bar</B></Sample>`
		expect(t, rec.Code, rw.Status())
		expect(t, rec.Body.String(), expect_body)
		expect(t, rw.Status(), http.StatusOK)
		expect(t, rw.Size(), len(expect_body))
		expect(t, rw.Written(), true)
	}()
	func() {
		rec := httptest.NewRecorder()
		rw := NewResponseWriter(rec)
		rw.WriteHTML("Hello %s", "world")
		expect_body := `Hello world`
		expect(t, rec.Code, rw.Status())
		expect(t, rec.Body.String(), expect_body)
		expect(t, rw.Status(), http.StatusOK)
		expect(t, rw.Size(), len(expect_body))
		expect(t, rw.Written(), true)
	}()
	func() {
		rec := httptest.NewRecorder()
		rw := NewResponseWriter(rec)
		rw.WriteText("Hello %s", "world")
		expect_body := `Hello world`
		expect(t, rec.Code, rw.Status())
		expect(t, rec.Body.String(), expect_body)
		expect(t, rw.Status(), http.StatusOK)
		expect(t, rw.Size(), len(expect_body))
		expect(t, rw.Written(), true)
	}()
	func() {
		dir, _ := ioutil.TempDir("", "xslate-test-")
		ioutil.WriteFile(dir+"/sample.tx", []byte("Hello [% name.to %]"), 0644)
		defer func() {
			os.RemoveAll(dir)
		}()
		Xslate.Configure(xslate.Args{
			"Loader": xslate.Args{
				"LoadPaths": []string{dir},
			},
		})
		rec := httptest.NewRecorder()
		rw := NewResponseWriter(rec)
		data := map[string]interface{}{
			"name": map[string]string{
				"to": "meiporo",
			},
		}
		rw.Render("sample.tx", data)
		expect_body := `Hello meiporo`
		expect(t, rec.Code, rw.Status())
		expect(t, rec.Body.String(), expect_body)
		expect(t, rw.Status(), http.StatusOK)
		expect(t, rw.Size(), len(expect_body))
		expect(t, rw.Written(), true)
	}()
}
