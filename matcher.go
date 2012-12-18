/**
 * Created with IntelliJ IDEA.
 * User: azeem
 * Date: 12/12/12
 * Time: 9:50 AM
 * To change this template use File | Settings | File Templates.
 */
package router

import "strconv"

// An named object that matches a string and extracts
// a value.
type Matcher interface {
	Match(target string) (status bool, value interface{})
	GetName() string
}

// An object that can be named
type Named struct {
	name string
}

// returns the name of the object
func (named Named) GetName() string {
	return named.name
}

// A Matcher that checks whether the
// target string is exactly equal to another
// string
type ExactMatcher struct {
	Named
	Rhs string
}

// creates a new ExactMatcher object that
// checks strings agains rhs
func Exact(rhs string) ExactMatcher {
	return ExactMatcher{Rhs:rhs}
}

// matches target string exactly
func (matcher ExactMatcher)Match(lhs string) (status bool, value interface{}) {
	if matcher.Rhs == lhs {
    	return true, lhs
	}
    return false, nil
}

// a Matcher that checks for integer values
type IntegerMatcher struct {
	Named
	base int
}

// creates a new IntegerMatcher with base 10
func Integer() IntegerMatcher {
	return IntegerMatcher{base:10}
}

// sets the name of the matcher
func (matcher IntegerMatcher)Name(name string) IntegerMatcher {
	matcher.name = name
	return matcher
}

// sets the radix of integers to be matched
func (matcher IntegerMatcher)Base(base int) IntegerMatcher {
	matcher.base = base
	return matcher
}

// matches an Integer
func (matcher IntegerMatcher)Match(target string) (status bool, value interface {}) {
	i, err := strconv.ParseInt(target, matcher.base, 32)
	if err != nil {
		return false, nil
	}
	return true, int(i)
}
