package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/c4e8ece0/yaxml"
)

// Debug values
const (
	FILE_MISSPELL = "000020b9fd4ce0a85f62881628ef9137.xml.0"
	FILE_REASK    = "000228ace7127d8c78820e2fa463bc63.xml.0"
	FILE_FIRST    = "000a8adf88587086e242ca0303678cfa.xml.0"
	DIR_SAMPLE    = "../sample/"

	CUR_SAMPLE = DIR_SAMPLE + FILE_REASK
)

// Try parse singe file
func main() {
	file, e := os.OpenFile(CUR_SAMPLE, os.O_RDONLY, 0777)
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
	defer file.Close()

	t, e := yaxml.Parse(file)
	// fmt.Printf("%#v %#v", t, e)
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}

	j, e := json.MarshalIndent(t, "", "\t")
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}

	ioutil.WriteFile("json", j, 0777)
	fmt.Println("done")
}
