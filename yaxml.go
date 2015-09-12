// Packages yaxml parses results Yandex.XML XML
package main

import (
	"encoding/xml"
	//	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html/charset"
)

// ALARM: if > -1 RequestGroupingsGroupBy.CurCateg
// ALARM: not closed <categ /> ResponseResultsGroupingGroupCateg
// ALARM: exists yandex xml relevance <relevance priority="phrase">104363467</relevance>
// ALARM: len(ResponseResultsGrouping) > 1
// ALARM: len(ResponseResultsGroupingGroupDoc) > 1
// ALARM: ResponseResultsGroupingGroupDocProperties.Type
//
// MISSED: https://tech.yandex.ru/xml/doc/dg/concepts/response_response-el-docpage/#misspell-block
//
// NOTICE: headline Для формирования используется HTML-тег meta, содержащий атрибут name со значением «description».
// NOTICE: HLWORK (InnerXML) для объекта надо собирать отдельно

// Make struct from xml
func Parse(r io.Reader) (Tree, error) {
	x := xml.NewDecoder(r)
	x.CharsetReader = charset.NewReaderLabel
	v := Tree{}
	e := x.Decode(&v)
	return v, e
}

// TODO: oop version
func New() {}

// TODO: dummy alarmer for constructor
func NoAlarm() {}

// TODO: Merge pages
func Merge(to Tree, from Tree) {}
func MergeAll(arr ...Tree)     {}
func MergeList(arr []Tree)     {}

// --------------------------------------------------------------------------
// XML Root
type Tree struct {
	XMLName xml.Name `xml:"yandexsearch"`
	//Request  Request  `xml:"request"`
	Response Response `xml:"response"`
}

// --------------------------------------------------------------------------
// Common blocks
type InnerXML struct {
	Content string `xml:",innerxml"`
}

// Extract highlighted text
func (x *InnerXML) HL() []string { return make([]string, 0) }

// Extract clean text
func (x *InnerXML) Text() []string { return make([]string, 0) }

// --------------------------------------------------------------------------
// <response> block
type Response struct {
	Date       string          `xml:"date,attr"`
	ReqID      string          `xml:"reqid"`
	Found      []ResponseFound `xml:"found"`
	FoundHuman string          `xml:"found-human"`
	Results    ResponseResults `xml:"results"`
}
type ResponseFound struct {
	Found    int    `xml:",chardata"`
	Priority string `xml:"priority,attr"`
}
type ResponseResults struct {
	Grouping []ResponseResultsGrouping `xml:"grouping"`
}
type ResponseResultsGrouping struct {
	Attr           string                             `xml:"attr,attr"`
	Mode           string                             `xml:"mode,attr"`
	GroupsOnPage   int                                `xml:"groups-on-page,attr"`
	DocsInGroup    int                                `xml:"docs-in-group,attr"`
	CurCateg       int                                `xml:"curcateg,attr"`
	Found          []ResponseResultsGroupingFound     `xml:"found"`
	FoundDocs      []ResponseResultsGroupingFoundDocs `xml:"found-docs"`
	FoundDocsHuman string                             `xml:"found-docs-human"`
	Page           ResponseResultsGroupingPage        `xml:"page"`
	Group          []ResponseResultsGroupingGroup     `xml:"group"`
}
type ResponseResultsGroupingFound struct {
	ResponseFound
}
type ResponseResultsGroupingFoundDocs struct {
	ResponseFound
}
type ResponseResultsGroupingPage struct {
	Page  int `xml:",chardata"`
	First int `xml:"first,attr"`
	Last  int `xml:"last,attr"`
}
type ResponseResultsGroupingGroup struct {
	Categ     []ResponseResultsGroupingGroupCateg   `xml:"categ"` // ALARM !<categ attr=">d<"
	Doccount  int                                   `xml:"doccount"`
	Relevance ResponseResultsGroupingGroupRelevance `xml:"relevance"`
	Doc       []ResponseResultsGroupingGroupDoc     `xml:"doc"`
}
type ResponseResultsGroupingGroupCateg struct {
	Categ string `:",chardata"` // ALARM!!
	Attr  string `xml:"attr,attr"`
	Name  string `xml:"name,attr"` // hostname
}
type ResponseResultsGroupingGroupRelevance struct {
	// history <relevance priority="phrase">104363467</relevance>
	Relevance string `xml:",chardata"`     // ALARM <-  panic on int
	Priority  string `xml:"priority,attr"` // ALARM
}

type ResponseResultsGroupingGroupDoc struct {
	Id         string                                    `xml:"id,attr"`
	Relevance  ResponseResultsGroupingGroupDocRelevance  `xml:"relevance"`
	URL        string                                    `xml:"url"`
	Domain     string                                    `xml:"domain"`
	Title      InnerXML                                  `xml:"title"`
	Headline   InnerXML                                  `xml:"headline"`
	ModTime    string                                    `xml:"modtime"`
	Size       int                                       `xml:"size"`
	Charset    string                                    `xml:"charset"`
	Passages   ResponseResultsGroupingGroupDocPassages   `xml:"passages"`
	Properties ResponseResultsGroupingGroupDocProperties `xml:"properties"`
}
type ResponseResultsGroupingGroupDocRelevance struct {
	ResponseResultsGroupingGroupRelevance
}
type ResponseResultsGroupingGroupDocPassages struct {
	Passage []InnerXML `xml:"passage"`
}

type ResponseResultsGroupingGroupDocProperties struct {
	// «0» — стандартный пассаж (сформирован из текста документа);
	// «1» — пассаж на основе текста ссылки. Используется, если документ найден по ссылке.
	PassagesType int    `xml:"_PassagesType"` // ALARM
	Lang         string `xml:"lang"`
}

// --------------------------------------------------------------------------
// <request> block - DONE
type Request struct {
	Query       string           `xml:"query"`
	Page        int              `xml:"page"`
	SortBy      RequestSortBy    `xml:"sortby"`
	MaxPassages int              `xml:"maxpassages"`
	Groupings   RequestGroupings `xml:"groupings"`
}
type RequestSortBy struct {
	SortBy   string `xml:",chardata"`
	Order    string `xml:"order,attr"`
	Priority string `xml:"priority,attr"`
}
type RequestGroupings struct {
	GroupBy []RequestGroupingsGroupBy `xml:"groupby"`
}
type RequestGroupingsGroupBy struct {
	Attr         string `xml:"attr,attr"`
	Mode         string `xml:"mode,attr"`
	GroupsOnPage int    `xml:"groups-on-page,attr"`
	DocsInGroup  int    `xml:"docs-in-group,attr"`
	CurCateg     int    `xml:"curcateg,attr"` // ALARM if > -1
}

// --------------------------------------------------------------------------
// Debug values
const (
	DIR_SAMPLE = "sample/"
	CUR_SAMPLE = DIR_SAMPLE + "000a8adf88587086e242ca0303678cfa.xml.0"
)

//
func main() {
	t, e := singlefile(CUR_SAMPLE)
	fmt.Printf("%#v %#v", t, e)
}

//
func singlefile(path string) (Tree, error) {
	file, e := os.OpenFile(path, os.O_RDONLY, 0777)
	if e != nil {
		return Tree{}, e
	}

	return Parse(file)
}
