// ahocorasick_test.go: test suite for ahocorasick
//
// Copyright (c) 2013 CloudFlare, Inc.

package ahocorasick

import (
	"regexp"
	"strings"
	"testing"
)

func assert(t *testing.T, b bool) {
	if !b {
		t.Fail()
	}
}

func TestNoPatterns(t *testing.T) {
	m := NewMatcher()
	hits := m.Match([]byte("foo bar baz"))
	assert(t, len(hits) == 0)
}

func TestNoData(t *testing.T) {
	m := NewMatcher()
	m.Append([]byte("foo"), 0)
	m.Append([]byte("bar"), 1)
	m.Append([]byte("baz"), 2)
	m.Finalize()
	hits := m.Match([]byte(""))
	assert(t, len(hits) == 0)
}

func TestSuffixes(t *testing.T) {
	m := NewMatcher()
	m.Append([]byte("Superman"), 0)
	m.Append([]byte("uperman"), 1)
	m.Append([]byte("perman"), 2)
	m.Append([]byte("erman"), 3)
	m.Finalize()
	hits := m.Match([]byte("The Man Of Steel: Superman"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 2)
	assert(t, hits[3] == 3)
}

func TestPrefixes(t *testing.T) {
	m := NewMatcher()
	m.Append([]byte("Superman"), 0)
	m.Append([]byte("uperman"), 1)
	m.Append([]byte("perman"), 2)
	m.Append([]byte("erman"), 3)
	m.Finalize()

	hits := m.Match([]byte("The Man Of Steel: Superman"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == 3)
	assert(t, hits[1] == 2)
	assert(t, hits[2] == 1)
	assert(t, hits[3] == 0)
}

func TestInterior(t *testing.T) {
	m := NewMatcher()
	m.Append([]byte("Steel"), 0)
	m.Append([]byte("tee"), 1)
	m.Append([]byte("e"), 2)
	m.Finalize()

	hits := m.Match([]byte("The Man Of Steel: Superman"))
	assert(t, len(hits) == 3)
	assert(t, hits[2] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[0] == 2)
}

func TestMatchAtStart(t *testing.T) {
	m := NewMatcher()
	m.Append([]byte("The"), 0)
	m.Append([]byte("Th"), 1)
	m.Append([]byte("he"), 2)
	m.Finalize()

	hits := m.Match([]byte("The Man Of Steel: Superman"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == 1)
	assert(t, hits[1] == 0)
	assert(t, hits[2] == 2)
}

func TestMatchAtEnd(t *testing.T) {
	m := NewMatcher()
	m.Append([]byte("teel"), 0)
	m.Append([]byte("eel"), 1)
	m.Append([]byte("el"), 2)
	m.Finalize()

	hits := m.Match([]byte("The Man Of Steel"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 2)
}

func TestOverlappingPatterns(t *testing.T) {
	m := NewMatcher()
	m.Append([]byte("Man"), 0)
	m.Append([]byte("n Of"), 1)
	m.Append([]byte("Of S"), 2)
	m.Finalize()

	hits := m.Match([]byte("The Man Of Steel"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 2)
}

func TestMultipleMatches(t *testing.T) {
	m := NewMatcher()
	m.Append([]byte("The"), 0)
	m.Append([]byte("Man"), 1)
	m.Append([]byte("an"), 2)
	m.Finalize()

	hits := m.Match([]byte("A Man A Plan A Canal: Panama, which Man Planned The Canal"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == 1)
	assert(t, hits[1] == 2)
	assert(t, hits[2] == 0)
}

func TestSingleCharacterMatches(t *testing.T) {
	m := NewMatcher()
	m.Append([]byte("a"), 0)
	m.Append([]byte("M"), 1)
	m.Append([]byte("z"), 2)
	m.Finalize()

	hits := m.Match([]byte("A Man A Plan A Canal: Panama, which Man Planned The Canal"))
	assert(t, len(hits) == 2)
	assert(t, hits[0] == 1)
	assert(t, hits[1] == 0)
}

func TestNothingMatches(t *testing.T) {
	m := NewMatcher()
	m.Append([]byte("baz"), 0)
	m.Append([]byte("bar"), 1)
	m.Append([]byte("foo"), 2)
	m.Finalize()

	hits := m.Match([]byte("A Man A Plan A Canal: Panama, which Man Planned The Canal"))
	assert(t, len(hits) == 0)
}

func TestWikipedia(t *testing.T) {
	m := NewMatcher()
	m.Append([]byte("a"), 0)
	m.Append([]byte("ab"), 1)
	m.Append([]byte("bc"), 2)
	m.Append([]byte("bca"), 3)
	m.Append([]byte("c"), 4)
	m.Append([]byte("caa"), 5)
	m.Finalize()

	hits := m.Match([]byte("abccab"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 2)
	assert(t, hits[3] == 4)

	hits = m.Match([]byte("bccab"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == 2)
	assert(t, hits[1] == 4)
	assert(t, hits[2] == 0)
	assert(t, hits[3] == 1)

	hits = m.Match([]byte("bccb"))
	assert(t, len(hits) == 2)
	assert(t, hits[0] == 2)
	assert(t, hits[1] == 4)
}

func TestMatch(t *testing.T) {
	m := NewMatcher()
	m.Append([]byte("Mozilla"), 0)
	m.Append([]byte("Mac"), 1)
	m.Append([]byte("Macintosh"), 2)
	m.Append([]byte("Safari"), 3)
	m.Append([]byte("Sausage"), 4)
	m.Finalize()

	hits := m.Match([]byte("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 2)
	assert(t, hits[3] == 3)

	hits = m.Match([]byte("Mozilla/5.0 (Mac; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 3)

	hits = m.Match([]byte("Mozilla/5.0 (Moc; Intel Computer OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36"))
	assert(t, len(hits) == 2)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 3)

	hits = m.Match([]byte("Mozilla/5.0 (Moc; Intel Computer OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Sofari/537.36"))
	assert(t, len(hits) == 1)
	assert(t, hits[0] == 0)

	hits = m.Match([]byte("Mazilla/5.0 (Moc; Intel Computer OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Sofari/537.36"))
	assert(t, len(hits) == 0)
}

var bytes = []byte("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36")
var sbytes = string(bytes)
var dictionary = []string{"Mozilla", "Mac", "Macintosh", "Safari", "Sausage"}
var dictionary2 = []string{"Googlebot", "bingbot", "msnbot", "Yandex", "Baiduspider"}
var dictionary3 = []string{"Mozilla", "Mac", "Macintosh", "Safari", "Phoenix"}
var dictionary4 = []string{"12343453", "34353", "234234523", "324234", "33333"}
var dictionary5 = []string{"12343453", "34353", "234234523", "324234", "33333", "experimental", "branch", "of", "the", "Mozilla", "codebase", "by", "Dave", "Hyatt", "Joe", "Hewitt", "and", "Blake", "Ross", "mother", "frequently", "performed", "in", "concerts", "around", "the", "village", "uses", "the", "Gecko", "layout", "engine"}
var precomputed = NewMatcher()
var precomputed2 = NewMatcher()
var precomputed3 = NewMatcher()
var precomputed4 = NewMatcher()
var precomputed5 = NewMatcher()

func fillTestData(m *Matcher, dict []string) {
	for i, s := range dict {
		m.Append([]byte(s), i)
	}
	m.Finalize()
}

func init() {
	fillTestData(precomputed, dictionary)
	fillTestData(precomputed2, dictionary2)
	fillTestData(precomputed3, dictionary3)
	fillTestData(precomputed4, dictionary4)
	fillTestData(precomputed5, dictionary5)
}

func BenchmarkMatchWorks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precomputed.Match(bytes)
	}
}

func BenchmarkContainsWorks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hits := make([]int, 0)
		for i, s := range dictionary {
			if strings.Contains(sbytes, s) {
				hits = append(hits, i)
			}
		}
	}
}

var re = regexp.MustCompile("(" + strings.Join(dictionary, "|") + ")")

func BenchmarkRegexpWorks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re.FindAllIndex(bytes, -1)
	}
}

func BenchmarkMatchFails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precomputed2.Match(bytes)
	}
}

func BenchmarkContainsFails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hits := make([]int, 0)
		for i, s := range dictionary2 {
			if strings.Contains(sbytes, s) {
				hits = append(hits, i)
			}
		}
	}
}

var re2 = regexp.MustCompile("(" + strings.Join(dictionary2, "|") + ")")

func BenchmarkRegexpFails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re2.FindAllIndex(bytes, -1)
	}
}

var bytes2 = []byte("Firefox is a web browser, and is Mozilla's flagship software product. It is available in both desktop and mobile versions. Firefox uses the Gecko layout engine to render web pages, which implements current and anticipated web standards. As of April 2013, Firefox has approximately 20% of worldwide usage share of web browsers, making it the third most-used web browser. Firefox began as an experimental branch of the Mozilla codebase by Dave Hyatt, Joe Hewitt and Blake Ross. They believed the commercial requirements of Netscape's sponsorship and developer-driven feature creep compromised the utility of the Mozilla browser. To combat what they saw as the Mozilla Suite's software bloat, they created a stand-alone browser, with which they intended to replace the Mozilla Suite. Firefox was originally named Phoenix but the name was changed so as to avoid trademark conflicts with Phoenix Technologies. The initially-announced replacement, Firebird, provoked objections from the Firebird project community. The current name, Firefox, was chosen on February 9, 2004.")
var sbytes2 = string(bytes2)

func BenchmarkLongMatchWorks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precomputed3.Match(bytes2)
	}
}

func BenchmarkLongContainsWorks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hits := make([]int, 0)
		for i, s := range dictionary3 {
			if strings.Contains(sbytes2, s) {
				hits = append(hits, i)
			}
		}
	}
}

var re3 = regexp.MustCompile("(" + strings.Join(dictionary3, "|") + ")")

func BenchmarkLongRegexpWorks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re3.FindAllIndex(bytes2, -1)
	}
}

func BenchmarkLongMatchFails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precomputed4.Match(bytes2)
	}
}

func BenchmarkLongContainsFails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hits := make([]int, 0)
		for i, s := range dictionary4 {
			if strings.Contains(sbytes2, s) {
				hits = append(hits, i)
			}
		}
	}
}

var re4 = regexp.MustCompile("(" + strings.Join(dictionary4, "|") + ")")

func BenchmarkLongRegexpFails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re4.FindAllIndex(bytes2, -1)
	}
}

func BenchmarkMatchMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precomputed5.Match(bytes)
	}
}

func BenchmarkContainsMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hits := make([]int, 0)
		for i, s := range dictionary4 {
			if strings.Contains(sbytes, s) {
				hits = append(hits, i)
			}
		}
	}
}

var re5 = regexp.MustCompile("(" + strings.Join(dictionary5, "|") + ")")

func BenchmarkRegexpMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re5.FindAllIndex(bytes, -1)
	}
}

func BenchmarkLongMatchMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precomputed5.Match(bytes2)
	}
}

func BenchmarkLongContainsMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hits := make([]int, 0)
		for i, s := range dictionary4 {
			if strings.Contains(sbytes2, s) {
				hits = append(hits, i)
			}
		}
	}
}

func BenchmarkLongRegexpMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re5.FindAllIndex(bytes2, -1)
	}
}
