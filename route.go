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
	"net/url"
	"strings"
	"errors"
	"fmt"
)

var (
	ErrRouteNotFound = errors.New("A Route with the given name was not found")

	ErrNonameMatcherInPathGen = errors.New("Matcher without name was encountered in URL Path generation")
	ErrMissingValueInPathGen = errors.New("Value not provided for a named matcher in URL Path generation")
)

type Routes []*Route

// Matches http request against routes and returns
// the first non nil MatchResult, returns nil otherwise
func (routes Routes)MatchRoute(req *http.Request) (result *MatchResult) {
	for _, route := range(routes) {
		if result := route.Match(req); result != nil {
			return result
		}
	}
	return nil
}

// finds the route with the given name
func (routes Routes)Find(name string) (route *Route) {
	for _, route := range(routes) {
		if route.name == name {
			return route
		}
	}
	return nil
}

// generates the url, back from the route with the given name
func (routes Routes)Url(name string, vars map[string]interface {}) (genUrl *url.URL, err error) {
	route := routes.Find(name)
	if route == nil {
		return nil, ErrRouteNotFound
	}
	return route.Url(vars)
}

// An http request Route
type Route struct {
	reqMatchers []HttpRequestMatcher
	handle interface{}
	name string
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

func (route *Route)GetName() (name string) {
	return route.name
}

func (route *Route)Url(vars map[string]interface {}) (genUrl *url.URL, err error) {
	genUrl = &url.URL{}
	for _, reqMatcher := range(route.reqMatchers) {
		if genErr := reqMatcher.Generate(genUrl, vars); genErr != nil {
			return nil, genErr
		}
	}
	return genUrl, nil
}

//////////////////////////////////
// Chain Methods for Route object

func (route *Route)Name(name string) (*Route) {
	route.name = name
	return route
}

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

//creates a new host name matcher
func (route *Route)Host(elems ...interface {}) (*Route) {
	hostMatcher := &HostMatcher{}
	for _, elem := range(elems) {
		hostMatcher.elemMatchers = append(hostMatcher.elemMatchers, makeMatcher(elem))
	}
	route.reqMatchers = append(route.reqMatchers, hostMatcher)
	return route
}

// creates a new Method matcher
func (route *Route)Method(method string) (*Route) {
	methodMatcher := &MethodMatcher{method:method}
	route.reqMatchers = append(route.reqMatchers, methodMatcher)
	return route
}

func (route *Route)Scheme(scheme string) (*Route) {
	schemeMatcher := &SchemeMatcher{scheme:scheme}
	route.reqMatchers = append(route.reqMatchers, schemeMatcher)
	return route
}

// End Chain Methods
//////////////////////////////////


// A Matcher that matches HTTP request
type HttpRequestMatcher interface {
	Match(req *http.Request) (status bool, vars map[string]interface{})
	Generate(url *url.URL, vars map[string]interface{}) (err error)
}

// An HttpRequestMatcher that matches
// htt request host names
type HostMatcher struct {
	elemMatchers []Matcher
}

func (hostMatcher *HostMatcher) Match(req *http.Request) (status bool, vars map[string]interface{}) {
	vars = make(map[string]interface {})
	elems := strings.Split(req.Host, ".")
	for i := 0;i < len(hostMatcher.elemMatchers);i++ {
		if !extractMatchValue(hostMatcher.elemMatchers[i], elems[i], vars) {
			return false, nil
		}
	}
	return true, vars
}

func (hostMatcher *HostMatcher) Generate(genUrl *url.URL, vars map[string]interface {}) (err error) {
	hostElems, err := generateMatchValues(hostMatcher.elemMatchers, vars)
	if err != nil {
		return nil
	}
	genUrl.Host = strings.Join(hostElems, ".")
	return nil
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

// generates url path
func (path *PathMatcher)Generate(genUrl *url.URL, vars map[string]interface {}) (err error) {
	pathElems, err := generateMatchValues(path.elemMatchers, vars)
	if err != nil {
		return err
	}
	genUrl.Path = "/" + strings.Join(pathElems, "/")
	return nil
}


// An HttpMethodMatcher that matches
// http request method
type MethodMatcher struct {
	method string
}

// matches http.Request.Method
func (methodMatcher *MethodMatcher) Match(req *http.Request) (status bool, vars map[string]interface{}) {
	if strings.ToUpper(methodMatcher.method) == strings.ToUpper(req.Method) {
		return true, nil
	}
	return false, nil
}

// generates nothing
func (methodMatcher *MethodMatcher) Generate(genUrl *url.URL, vars map[string]interface {}) (err error) {
	return nil
}

// An HttpMethodMatcher that matches
// http request method
type SchemeMatcher struct {
	scheme string
}

// matches http.Request.Method
func (schemeMatcher *SchemeMatcher) Match(req *http.Request) (status bool, vars map[string]interface{}) {
	if schemeMatcher.scheme == req.URL.Scheme {
		return true, nil
	}
	return false, nil
}

// sets scheme in the url
func (schemeMatcher *SchemeMatcher) Generate(genUrl *url.URL, vars map[string]interface {}) (err error) {
	genUrl.Scheme = schemeMatcher.scheme
	return nil
}

///////////////////////
// some util functions

func generateMatchValues(matchers []Matcher, vars map[string]interface {}) (vals []string, err error) {
	for _, matcher := range(matchers) {
		if exactMatcher, ok := matcher.(ExactMatcher); ok {
			vals = append(vals, exactMatcher.Rhs)
		} else {
			matcherName := matcher.GetName()
			if matcherName == "" {
				return nil, ErrNonameMatcherInPathGen
			}
			value, exists := vars[matcherName]
			if !exists {
				return nil, ErrMissingValueInPathGen
			}
			vals = append(vals, fmt.Sprintf("%s", value))
		}
	}
	return vals, nil
}

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
