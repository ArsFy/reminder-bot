package main

import "regexp"

var Dtime, Htime, Mtime *regexp.Regexp

func init() {
	Dtime, _ = regexp.Compile(`([0-9][0-9]|[0-9])\-([0-9][0-9]|[0-9]),([0-9][0-9]|[0-9]):([0-9][0-9]|[0-9]),(.*)`)
	Htime, _ = regexp.Compile(`([0-9][0-9]|[0-9]):([0-9][0-9]|[0-9]),(.*)`)
	Mtime, _ = regexp.Compile(`([0-9][0-9]|[0-9]),([0-9][0-9]|[0-9]):([0-9][0-9]|[0-9]),(.*)`)
}
