/**
 * Created with IntelliJ IDEA.
 * User: azeem
 * Date: 13/12/12
 * Time: 9:47 AM
 * To change this template use File | Settings | File Templates.
 */
package router

import (
	"net/http"
	"strings"
)


// helper function that checks an http request against an
// array of routes and returns the first non nil MatchResult
// returns nil otherwise
func MatchRoute(routes []*Route, req *http.Request) (result *MatchResult) {
	for _, route := range(routes) {
		if result := route.Match(req); result != nil {
			return result
		}
	}
	return nil
}

// An http request Route
type Route struct {
	reqMatchers []HttpRequestMatcher
	handle interface{}
}

// a Match result object
type MatchResult struct {
	Vars map[string]interface{}
	Handle interface{}
}

// creates a new Route
func NewRoute() *Route {
	return &Route{}
}

// matches an http.Request against this Route
func (route *Route)Match(req *http.Request) (result *MatchResult) {
	match := &MatchResult{}
	match.Vars = make(map[string]interface {})
	for _, reqMatcher := range(route.reqMatchers) {
		if status, vars := reqMatcher.Match(req); status {
			for key, value := range(vars) {
				match.Vars[key] = value
			}
		} else {
			return nil
		}
	}
	match.Handle = route.handle
	return match
}

//////////////////////////////////
// Chain Methods for Route object

// sets the route Handle
func (route *Route)Handle(handle interface{}) (*Route) {
	route.handle = handle
	return route
}

// creates a new path matcher. elems are Matcher
// objects for each path component. strings are converted
// to ExactMatcher objects
func (route *Route)Path(elems ...interface{}) (*Route) {
	path := &PathMatcher{}
	for _, elem := range(elems) {
		path.elemMatchers = append(path.elemMatchers, makeMatcher(elem))
	}
	route.reqMatchers = append(route.reqMatchers, path)
	return route
}

// creates a new Method matcher
func (route *Route)Method(method interface{}) (*Route) {
	methodMatcher := &MethodMatcher{}
	methodMatcher.matcher = makeMatcher(method)
	route.reqMatchers = append(route.reqMatchers, methodMatcher)
	return route
}

// End Chain Methods
//////////////////////////////////


// A Matcher that matches HTTP request
type HttpRequestMatcher interface {
	Match(req *http.Request) (status bool, vars map[string]interface{})
}

// An HttpRequestMatcher that matches
// http request uri path
type PathMatcher struct {
	elemMatchers []Matcher
}

// matches http.Request.RequestURI
func (path *PathMatcher) Match(req *http.Request) (status bool, vars map[string]interface{}) {
	vars = make(map[string]interface {})
	elems := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	for i := 0;i < intMin(len(path.elemMatchers), len(elems)); i++ {
		if !extractMatchValue(path.elemMatchers[i], elems[i], vars) {
			return false, nil
		}
	}
	return true, vars
}


// An HttpMethodMatcher that matches
// http request method
type MethodMatcher struct {
	matcher Matcher
}

func (methodMatcher *MethodMatcher) Match(req *http.Request) (status bool, vars map[string]interface{}) {
	vars = make(map[string]interface{})
	if !extractMatchValue(methodMatcher.matcher, req.Method, vars) {
		return false, nil
	}
	return true, vars
}

///////////////////////
// some util functions

func extractMatchValue(matcher Matcher, target string, vars map[string]interface {}) bool {
	if status, value := matcher.Match(target); status == true {
		matcherName := matcher.GetName()
		if matcherName != "" {
			vars[matcherName] = value
		}
		return true
	}
	return false
}

func makeMatcher(value interface{}) (matcher Matcher) {
	if m, ok := value.(Matcher); ok {
		return m
	} else if strValue, ok := value.(string); ok {
		return Exact(strValue)
	}
	panic("Invalid path element type")
}

// finds the minimum between two ints
func intMin(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}
