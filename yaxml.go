// Package yaxml parses and fetches xml-data from Yandex.XML
package yaxml

import (
	"fmt"
)

type Opts struct {
	User    string
	Key     string
	Perpage int    // results per request
	Timeout int    // seconds
	LR      int    // 213 = Moscow
	Mode    string // deep
	Depth   int    // maximum depth

	OnFailRepeats      int // 3
	OnFailRepeatsTotal int // 30
}

// Default request for interner-wide SERP
var OptsDefault = &Opts{
	User:               "",
	Key:                "",
	Perpage:            100,
	Timeout:            5,
	LR:                 213,
	Mode:               "deep",
	Depth:              100,
	OnFailRepeats:      3,
	OnFailRepeatsTotal: 30,
}

func main() {
	fmt.Print(1)
}
