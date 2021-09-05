package main

import (
	"fmt"
	"regexp"
	"sort"
)

type matcher []*regexpMatcher

func (rms matcher) Match(s string) []Range {
	var ms []match
	for _, m := range rms {
		ms = m.AppendMatches(s, ms)
	}
	ms = rms.removeOverlaps(ms)

	rs := make([]Range, len(ms))
	for i, m := range ms {
		rs[i] = m.Sel
	}
	return rs
}

func (rms matcher) removeOverlaps(ms []match) []match {
	if len(ms) < 2 {
		return ms
	}

	// Sort in ascending order by:
	// - Starts earliest
	// - Runs longest
	sort.Slice(ms, func(i, j int) bool {
		l, r := ms[i].Full, ms[j].Full

		if l.Start < r.Start {
			return true
		}
		if l.Start > r.Start {
			return false
		}

		return l.Len() > r.Len()
	})

	out := ms[:1]
	for _, m := range ms[1:] {
		if m.Full.Start < out[len(out)-1].Full.End {
			continue
		}
		out = append(out, m)
	}

	return out
}

type regexpMatcher struct {
	name   string
	regex  *regexp.Regexp
	subexp int
}

// compileRegexpMatcher builds a regexpMatcher with the provided name and
// regular expression.
func compileRegexpMatcher(name, s string) (*regexpMatcher, error) {
	if len(s) == 0 {
		return &regexpMatcher{name: name}, nil
	}

	re, err := regexp.Compile(s)
	if err != nil {
		return nil, err
	}

	n := 0
	if re.NumSubexp() > 0 {
		n++
	}

	return &regexpMatcher{
		regex:  re,
		subexp: n,
		name:   name,
	}, nil
}

func (rm *regexpMatcher) Name() string {
	return rm.name
}

func (rm *regexpMatcher) String() string {
	return fmt.Sprintf("%v:%v", rm.name, rm.regex)
}

type match struct {
	// Full matched area.
	Full Range

	// Selected portion that will be copied.
	Sel Range
}

func (rm *regexpMatcher) AppendMatches(s string, ms []match) []match {
	if rm.regex == nil {
		return ms
	}
	for _, m := range rm.regex.FindAllStringSubmatchIndex(s, -1) {
		ms = append(ms, match{
			Full: Range{Start: m[0], End: m[1]},
			Sel:  Range{Start: m[2*rm.subexp], End: m[2*rm.subexp+1]},
		})
	}
	return ms
}