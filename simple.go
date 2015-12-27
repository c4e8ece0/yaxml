package yaxml

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// New() create new bucket for XML-files
// Each call to $.Parse(io.Reader) will append data to current bucket
func New() *Simple {
	s := &Simple{}
	s.List = make([]*SimpleDoc, 0)
	return s
}

// Root element
type Simple struct {
	Query string // query text
	Stat  SimpleStat
	List  []*SimpleDoc

	Reask     bool   // fail flag
	ReaskType string // fail name (Misspell, KeyboardLayout, Volapyuk)
	Variant   string // new query
}

// HL collect highlighted text from all fields
func (sm *Simple) HL() []string {
	r := make([]string, 0)
	for i, _ := range sm.List {
		r = append(r, sm.List[i].HL.Total()...)
	}
	return r
}

// AppendString appends new data from string
// to current bucket
func (sm *Simple) AppendString(s string) error {
	return sm.AppendReader(strings.NewReader(s))
}

// AppendReader appends new data from io.Reader
// to current bucket
func (sm *Simple) AppendReader(r io.Reader) error {
	tree, err := Parse(r)
	if err != nil {
		return err
	}

	if len(sm.List) == 0 {
		sm.Query = tree.Request.Query
		sm.Stat.Approx = uint64(tree.Response.Found[0].Found)
		sm.Stat.Docs = uint64(tree.Response.Results.Grouping[0].FoundDocs[0].Found)
		sm.Stat.Hosts = uint64(tree.Response.Results.Grouping[0].Found[0].Found)
		if tree.Response.Reask.Rule != "" {
			sm.Reask = true
			sm.ReaskType = tree.Response.Reask.Rule
			sm.Variant = tree.Response.Reask.Text
		} else if tree.Response.Misspell.Rule != "" {
			sm.Reask = true
			sm.ReaskType = tree.Response.Misspell.Rule
			sm.Variant = tree.Response.Misspell.Text
		}
	}

	for i, _ := range tree.Response.Results.Grouping[0].Group {
		t := &SimpleDoc{}
		t.URL = tree.Response.Results.Grouping[0].Group[i].Doc[0].URL
		t.Domain = tree.Response.Results.Grouping[0].Group[i].Doc[0].Domain
		t.Charset = tree.Response.Results.Grouping[0].Group[i].Doc[0].Charset
		t.Lang = tree.Response.Results.Grouping[0].Group[i].Doc[0].Properties.Lang
		t.Docs = uint64(tree.Response.Results.Grouping[0].Group[i].Doccount)

		t.Title = tree.Response.Results.Grouping[0].Group[i].Doc[0].Title.Text()
		t.HL.Title = append(t.HL.Title, tree.Response.Results.Grouping[0].Group[i].Doc[0].Title.HL()...)

		t.Headline = tree.Response.Results.Grouping[0].Group[i].Doc[0].Headline.Text()
		t.HL.Headline = append(t.HL.Headline, tree.Response.Results.Grouping[0].Group[i].Doc[0].Headline.HL()...)

		for _, str := range tree.Response.Results.Grouping[0].Group[i].Doc[0].Passages.Passage {
			t.Passages = append(t.Passages, str.Text())
			t.HL.Passages = append(t.HL.Passages, str.HL()...)
		}

		sm.List = append(sm.List, t)
	}

	t, _ := json.MarshalIndent(sm, "", "   ")
	fmt.Printf("%v", string(t))
	fmt.Printf("%#v", sm.HL())

	return nil
}

// Statistics of first parsed XML
type SimpleStat struct {
	Approx uint64
	Hosts  uint64
	Docs   uint64
}

// Description of each element in XML
type SimpleDoc struct {
	URL     string
	Domain  string
	Charset string
	Lang    string

	Title    string
	Headline string
	Passages []string

	Docs uint64
	HL   SimpleHL
}

// Highlighted words in fields
type SimpleHL struct {
	Title    []string
	Headline []string
	Passages []string
}

// Full List of highlighted words
func (s SimpleHL) Total() []string {
	r := make([]string, 0)
	r = append(r, s.Title...)
	r = append(r, s.Headline...)
	r = append(r, s.Passages...)
	return r
}

// // HL collect highlighted text from field "title"
// func (sm *Simple) TitleHL() []string {
// 	return []string{}
// }

// // HL collect highlighted text from field "headline"
// func (sm *Simple) HeadlineHL() []string {
// 	return []string{}
// }

// // HL collect highlighted text from field "passage"
// func (sm *Simple) PassageHL() []string {
// 	return []string{}
// }
