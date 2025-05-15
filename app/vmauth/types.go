package main

import (
	"regexp"
)

type URLMap struct {
	SrcHosts     []*Regex
	SrcPaths     []*Regex
	SrcQueryArgs []*QueryArg
	SrcHeaders   []*Header
	URLPrefix    *URLPrefix
	HeadersConf  HeadersConf
}

type Regex struct {
	pattern *regexp.Regexp
}

func NewRegex(pattern string) (*Regex, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Regex{pattern: r}, nil
}

func (r *Regex) MatchString(s string) bool {
	return r.pattern.MatchString(s)
}

type QueryArg struct {
	Name  string
	Value *Regex
}

type Header struct {
	Name  string
	Value string
}

type URLPrefix struct {
	URL string
}

type HeadersConf struct {
	Headers          []Header
	RequestHeaders   []*Header
	ResponseHeaders  []*Header
	KeepOriginalHost *bool
}
