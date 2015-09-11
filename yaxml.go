// Packages yaxml parses results Yandex.XML XML
package yaxml

import (
	"encoding/xml"
	"error"
	"io"
)

type Tree struct{}

func Parse(io.Reader) (Tree, error) {
}

func main() {
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
