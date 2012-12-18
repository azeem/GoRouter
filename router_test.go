/**
 * Created with IntelliJ IDEA.
 * User: azeem
 * Date: 18/12/12
 * Time: 12:55 PM
 * To change this template use File | Settings | File Templates.
 */
package router

import (
	"net/http"
	"testing"
	"fmt"
)

func testRoute(t *testing.T, title string, route *Route, req *http.Request, vars map[string]interface{}, handle interface{}) {
	if res := route.Match(req); res != nil {
		if handle != res.Handle {
			t.Error(fmt.Sprintf("Route('%s') incorrect result handle", title))
		}
		for key, value := range (vars) {
			if res.Vars[key] != value {
				t.Error(fmt.Sprintf("Route('%s') result variable '%s' mismatch", title, key))
			}
		}
	} else {
		t.Error(fmt.Sprintf("Route('%s') match failed", title))
	}
}

func makeRequest(method string, urlStr string) *http.Request {
	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func TestRoute(t *testing.T) {
	testRoute(t, "Exact Match",
		NewRoute().Path("abc", "def").Handle(123),
		makeRequest("GET", "http://example.com/abc/def"),
		nil, 123)

	testRoute(t, "Integer Test",
		NewRoute().Path("abc", "def", Integer().Name("Test"), "end").Handle("Hello World"),
		makeRequest("GET", "http://example.com/abc/def/234/end"),
		map[string]interface{} {"Test":234}, "Hello World")

	testRoute(t, "Method Test",
		NewRoute().Path("abc", "def").Method("POST").Handle("handle"),
		makeRequest("POST", "http://example.com/abc/def"),
		nil, "handle")
}
