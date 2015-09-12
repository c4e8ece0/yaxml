// Packages yaxml parses results Yandex.XML XML
package yaxml

import (
	"encoding/xml"
	"error"
	"io"
	"io/ioutil"
)

// FIRST SAMPLE: <grouping attr="d" mode="deep" groups-on-page="100" docs-in-group="1" curcateg="-1">

type Tree struct {
}

//
func Parse(io.Reader) (Tree, error) {
}

// Debug values
const (
	DIR_SAMLE = "sample/"
	CUR_SAMLE = DIR_SAMLE + "000a8adf88587086e242ca0303678cfa.xml.0"
)

func main() {
	singlefile(CUR_SAMPLE)
}

func singlefile(f) {
	file := ioutil.ReadFile(CUR_SAMLE)
	// TODO: bububu
}

/*

// XML example


// Bar XML structs
type BarTcy struct {
	Rang  string `xml:"rang,attr"`
	Value string `xml:"value,attr"`
}

type BarUrl struct {
	Domain string `xml:"domain,attr"`
}
type BarUrlinfo struct {
	XMLName xml.Name `xml:"urlinfo"`
	Url     BarUrl   `xml:"url"`
	Tcy     BarTcy   `xml:"tcy"`
}

// Fetch Yandex.Bar XML
func Bar(host string) (string, error) {
	resp, err := http.Get("http://bar-navig.yandex.ru/u?ver=2&show=32&url=http://" + host)
	if err != nil {
		return "", err
	}

	v, e := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return string(v), e
}

*/
