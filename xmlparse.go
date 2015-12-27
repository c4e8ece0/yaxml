package yaxml

// File contains structs and funcs for parsing Yandex.XML to yaxml.Tree{}

import (
	//"charset"
	"encoding/xml"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
)

// Function NoAlarm is default dummy alarmer
// IT DOESN'T WORK
func NoAlarm(s string) {}

// Parse XML-file, closing io.Reader on the caller
func Parse(r io.Reader) (XMLTree, error) {
	x := xml.NewDecoder(r)
	x.CharsetReader = charset.NewReaderLabel
	v := XMLTree{}
	e := x.Decode(&v)
	return v, e
}

// InnerXML required for fields where <hlword> could be found
type InnerXML struct {
	Content string `xml:",innerxml"`
}

// Extract highlighted text
func (x *InnerXML) HL() []string {
	s := x.Content
	s = strings.Replace(s, "</hlword> <hlword>", " ", -1)

	f := regexp.MustCompile(`<hlword>(.*?)</hlword>`)
	r := []string{}
	for _, pair := range f.FindAllStringSubmatch(s, -1) {
		r = append(r, pair[1])
	}

	return r
}

// Extract clean text
func (x *InnerXML) Text() string {
	s := x.Content
	s = strings.Replace(s, "<hlword>", "", -1)
	s = strings.Replace(s, "</hlword>", "", -1)
	return s
}

// --------------------------------------------------------------------------
// XML Root
type XMLTree struct {
	Request  XMLRequest  `xml:"request"`
	Response XMLResponse `xml:"response"`
}

// --------------------------------------------------------------------------
// Helpers

// Common tag
type XMLRelevance struct { // We miss you...
	// <relevance priority="phrase">104363467</relevance>
	Relevance string `xml:",chardata"`
	Priority  string `xml:"priority,attr"`
}

// --------------------------------------------------------------------------
// <request> block - DONE
type XMLRequest struct {
	Query       string              `xml:"query"`
	Page        int                 `xml:"page"`
	SortBy      XMLRequestSortBy    `xml:"sortby"`
	MaxPassages int                 `xml:"maxpassages"`
	Groupings   XMLRequestGroupings `xml:"groupings"`
}

type XMLRequestSortBy struct {
	SortBy   string `xml:",chardata"`
	Order    string `xml:"order,attr"`
	Priority string `xml:"priority,attr"`
}

type XMLRequestGroupings struct {
	GroupBy []XMLRequestGroupingsGroupBy `xml:"groupby"`
}

type XMLRequestGroupingsGroupBy struct {
	Attr         string `xml:"attr,attr"`
	Mode         string `xml:"mode,attr"`
	GroupsOnPage int    `xml:"groups-on-page,attr"`
	DocsInGroup  int    `xml:"docs-in-group,attr"`
	CurCateg     int    `xml:"curcateg,attr"` // ALARM if > -1
}

// --------------------------------------------------------------------------
// <response> block
type XMLResponse struct {
	Date       string                     `xml:"date,attr"`
	ReqID      string                     `xml:"reqid"`
	Found      []XMLResponseFound         `xml:"found"`
	FoundHuman string                     `xml:"found-human"`
	Misspell   XMLResponseMisspellOrReask `xml:"misspell"`
	Reask      XMLResponseMisspellOrReask `xml:"reask"`
	Results    XMLResponseResults         `xml:"results"`
}

type XMLResponseFound struct {
	Found    int    `xml:",chardata"`
	Priority string `xml:"priority,attr"`
}

type XMLResponseMisspellOrReask struct {
	Rule       string   `xml:"rule"` // Misspell, KeyboardLayout, Volapyuk (ru translit)
	SourceText InnerXML `xml:"source-text"`
	Text       string   `xml:"text"`
}

type XMLResponseResults struct {
	Grouping []XMLResponseResultsGrouping `xml:"grouping"`
}

type XMLResponseResultsGrouping struct {
	Attr           string                                `xml:"attr,attr"`
	Mode           string                                `xml:"mode,attr"`
	GroupsOnPage   int                                   `xml:"groups-on-page,attr"`
	DocsInGroup    int                                   `xml:"docs-in-group,attr"`
	CurCateg       int                                   `xml:"curcateg,attr"`
	Found          []XMLResponseResultsGroupingFound     `xml:"found"`
	FoundDocs      []XMLResponseResultsGroupingFoundDocs `xml:"found-docs"`
	FoundDocsHuman string                                `xml:"found-docs-human"`
	Page           XMLResponseResultsGroupingPage        `xml:"page"`
	Group          []XMLResponseResultsGroupingGroup     `xml:"group"`
}

type XMLResponseResultsGroupingFound struct {
	XMLResponseFound
}

type XMLResponseResultsGroupingFoundDocs struct {
	XMLResponseFound
}

type XMLResponseResultsGroupingPage struct {
	Page  int `xml:",chardata"`
	First int `xml:"first,attr"`
	Last  int `xml:"last,attr"`
}

type XMLResponseResultsGroupingGroup struct {
	Categ     []XMLResponseResultsGroupingGroupCateg `xml:"categ"` // ALARM !<categ attr=">d<"
	Doccount  int                                    `xml:"doccount"`
	Relevance XMLRelevance                           `xml:"relevance"`
	Doc       []XMLResponseResultsGroupingGroupDoc   `xml:"doc"`
}

type XMLResponseResultsGroupingGroupCateg struct {
	Categ string `:",chardata"` // ALARM!!
	Attr  string `xml:"attr,attr"`
	Name  string `xml:"name,attr"` // hostname
}

type XMLResponseResultsGroupingGroupDoc struct {
	Id         string                                       `xml:"id,attr"`
	Relevance  XMLRelevance                                 `xml:"relevance"`
	URL        string                                       `xml:"url"`
	Domain     string                                       `xml:"domain"`
	Title      InnerXML                                     `xml:"title"`
	Headline   InnerXML                                     `xml:"headline"` // i.e. meta-description
	ModTime    string                                       `xml:"modtime"`
	Size       int                                          `xml:"size"`
	Charset    string                                       `xml:"charset"`
	Passages   XMLResponseResultsGroupingGroupDocPassages   `xml:"passages"`
	Properties XMLResponseResultsGroupingGroupDocProperties `xml:"properties"`
}

type XMLResponseResultsGroupingGroupDocPassages struct {
	Passage []InnerXML `xml:"passage"`
}

type XMLResponseResultsGroupingGroupDocProperties struct {
	// 0 — standart passage (on-page)
	// 1 — anchor based (on-link)
	PassagesType int    `xml:"_PassagesType"` // ALARM
	Lang         string `xml:"lang"`
}
