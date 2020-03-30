/*
   Copyright © 2019, 2020 M.Watermann, 10247 Berlin, Germany
              All rights reserved
          EMail : <support@mwat.de>
*/

package db

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"fmt"
	"regexp"
	"strings"
)

/*
 * This file provides helper functions and methods for database searches.
 */


// `escapeQuery()` returns a string with some characters escaped.
//
// see: https://github.com/golang/go/issues/18478#issuecomment-357285669
func escapeQuery(aSource string) string {
	sLen := len(aSource)
	if 0 == sLen {
		return ``
	}
	var (
		character byte
		escape    bool
		j         int
	)
	result := make([]byte, sLen<<1)
	for i := 0; i < sLen; i++ {
		character = aSource[i]
		switch character {
		case '\n', '\r', '\\', '"':
			// We do not escape the apostrophe since it can be
			// a legitimate part of the search term and we use
			// double quotes to enclose those terms.
			escape = true
		case '\032': // Ctrl-Z
			escape = true
			character = 'Z'
		default:
		}

		if escape {
			escape = false
			result[j] = '\\'
			result[j+1] = character
			j += 2
		} else {
			result[j] = character
			j++
		}
	}

	return string(result[0:j])
} // escapeQuery()


/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

type (
	// TSearch provides text search capabilities.
	TSearch struct {
		raw   string // the raw (unprocessed) search expression
		where string // used to build the WHERE clause
		next  string
	}

	tExpression struct {
		entity  string // the DB field to lookup
		matcher string // how to lookup
		not     bool   // flag negating the search result
		op      string // how to concat with the next expression
		term    string // what to lookup
	}
)

// `allSQL()` returns a WHERE clause to match the current term
// in all suitable tables.
func (exp *tExpression) allSQL() (rWhere string) {
	exp.matcher, exp.op = "~", "OR"

	exp.entity = "authors"
	rWhere = exp.buildSQL()
	exp.entity = "comment"
	rWhere += exp.buildSQL()
	exp.entity = "format"
	rWhere += exp.buildSQL()
	exp.entity = "language"
	rWhere += exp.buildSQL()
	exp.entity = "publisher"
	rWhere += exp.buildSQL()
	exp.entity = "series"
	rWhere += exp.buildSQL()
	exp.entity = "tags"
	rWhere += exp.buildSQL()
	exp.entity, exp.op = "title", ""
	rWhere += exp.buildSQL()

	return
} // allSQL()

// `buildSQL()` returns an SQL clause based on `exp` properties
// suitable for the Calibre database.
func (exp *tExpression) buildSQL() (rWhere string) {
	b := 2 // number of brackets to close
	switch exp.entity {
	case "authors", "author": // accept (wrong) "author"
		rWhere = `(b.id IN (SELECT ba.book FROM books_authors_link ba JOIN authors a ON(ba.author = a.id) WHERE (a.name`

	case "comment":
		rWhere = `(b.id IN (SELECT c.book FROM comments c WHERE (c.text`

	case "format":
		rWhere = `(b.id IN (SELECT d.book FROM data d WHERE (d.format`

	case "languages", "language": // accept (wrong) "language"
		rWhere = `(b.id IN (SELECT bl.book FROM books_languages_link bl JOIN languages l ON(bl.lang_code = l.id) WHERE (l.lang_code`

	case "publisher":
		rWhere = `(b.id IN (SELECT bp.book FROM books_publishers_link bp JOIN publishers p ON(bp.publisher = p.id) WHERE (p.name`

	case "series":
		rWhere = `(b.id IN (SELECT bs.book FROM books_series_link bs JOIN series s ON(bs.series = s.id) WHERE (s.name`

	case "tags", "tag": // accept (wrong) "tag"
		rWhere = `(b.id IN (SELECT bt.book FROM books_tags_link bt JOIN tags t ON(bt.tag = t.id) WHERE (t.name`

	case "title":
		b, rWhere = 1, `(b.title`

	default:
		if 0 == len(exp.entity) {
			return
		}
		field := strings.ToLower(exp.entity)
		if '#' != exp.entity[0] {
			field = "#" + field
		}
		if isCustom, err := MetaFieldValue(field, "is_custom"); (nil != err) || (true != isCustom) {
			return // no user-defined field
		}
		if isCategory, err := MetaFieldValue(field, "is_category"); (nil != err) || (true != isCategory) {
			return
		}
		iTable, err := MetaFieldValue(field, "table")
		if nil != err {
			return
		}
		table, ok := iTable.(string)
		if !ok {
			return
		}
		rWhere = fmt.Sprintf(`(b.id IN (SELECT lct.book FROM books_%s_link lct JOIN %s ct ON(lct.value = ct.id) WHERE (ct.value`, table, table) // #nosec G201
	}

	term := escapeQuery(exp.term)
	if "=" == exp.matcher {
		if exp.not {
			rWhere += ` != "`
		} else {
			rWhere += ` = "`
		}
		rWhere += term + `")`
	} else {
		if exp.not {
			rWhere += ` NOT`
		}
		rWhere += ` LIKE "%` + term + `%")`
	}
	if 1 < b {
		rWhere += `))`
	}
	if 0 < len(exp.op) {
		rWhere += exp.op
	}

	return
} // buildSQL()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// Clause returns the produced WHERE clause.
func (so *TSearch) Clause() (rWhere string) {
	if 0 < len(so.raw) {
		so.Parse()
	}
	if 0 < len(so.where) {
		rWhere = ` WHERE ` + so.where // #nosec G202
	}

	return
} // Clause()

/*
There are several forms to recognise:

"just a search term" => lookup ALL book entities;
`entity:"=searchterm"` => lookup exact match of `searchterm` in `entity`;
`entity:"~searchterm"` => lookup `searchterm` contained in `entity`.

All three expressions can be combined by AND and OR.
All three expressions can be negated by a leading `!`.
*/

var (
	// Lookup table for missing comparison values
	soReplacementLookup = map[string]string{
		`AND`: `(1=1)`,
		`OR`:  `(1=0)`,
		``:    ``,
	}

	// RegEx to find a search expression
	soSearchExpressionRE = regexp.MustCompile(
		`(?i)((!?)(#?\w+):)"([=~]?)([^"]*)"(\s*(AND|OR))?`)
	//       12222333333311 44444445555555 6666777777776

	soSearchRemainderRE = regexp.MustCompile(
		`\s*(!?)\s*([\w ]+)`)
	//      1      2
)

// `p1` returns the parsed search term(s).
func (so *TSearch) p1() *TSearch {
	op, p, s, w := "", 0, "", strings.TrimSpace(so.raw)
	for matches := soSearchExpressionRE.FindStringSubmatch(w); 7 < len(matches); matches = soSearchExpressionRE.FindStringSubmatch(w) {
		if 0 == len(matches[4]) {
			matches[4] = `=` // defaults to exact match
		}
		exp := &tExpression{
			entity:  strings.ToLower(matches[3]),
			not:     (`!` == matches[2]),
			matcher: matches[4],
			op:      strings.ToUpper(matches[7]),
			term:    matches[5],
		}
		s = exp.buildSQL()
		if 0 == len(s) {
			s = soReplacementLookup[op]
		}
		w = strings.Replace(w, matches[0], s, 1)
		p = strings.Index(w, s) + len(s)
		op = exp.op // save the latest operant for below
	}
	if p < len(w) { // check whether there's something behind the last expression
		matches := soSearchRemainderRE.FindStringSubmatch(w[p:])
		if 2 < len(matches) {
			exp := &tExpression{
				not:  (`!` == matches[1]),
				term: matches[2],
			}
			if 0 == len(op) {
				s = `OR ` + exp.allSQL()
			} else {
				s = exp.allSQL()
			}
			w = strings.Replace(w, matches[0], s, 1)
		}
	}
	so.next, so.raw, so.where = ``, ``, w

	return so
} // p1()

// Parse returns the parsed search term(s).
func (so *TSearch) Parse() *TSearch {
	so.raw = strings.TrimSpace(so.raw)
	if 0 == len(so.raw) {
		so.next, so.where = "", ""
		return so
	}
	if soSearchExpressionRE.MatchString(so.raw) {
		// This is moved to a separate method for easier testing.
		return so.p1()
	}

	exp := &tExpression{term: so.raw}
	so.where, so.raw = exp.allSQL(), ""

	return so
} // Parse()

// String returns a string field representation.
func (so *TSearch) String() string {
	return `raw: '` + so.raw +
		`' | where: '` + so.where +
		`' | next: '` + so.next + `'`
} // String()

// Where returns the SQL code to use in the WHERE clause.
func (so *TSearch) Where() string {
	return so.where
} // Where()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// NewSearch returns a new `TSearch` instance.
func NewSearch(aSearchTerm string) *TSearch {
	return &TSearch{raw: aSearchTerm}
} // NewSearch()

/* _EoF_ */
