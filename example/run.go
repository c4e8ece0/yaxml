package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/c4e8ece0/yaxml"
	_ "github.com/go-sql-driver/mysql"
)

// Action types
const (
	ACTION_SINGLE = iota
	ACTION_READDIR
	ACTION_MASS
)

// Debug values
const (
	FILE_MISSPELL = "000020b9fd4ce0a85f62881628ef9137.xml.0"
	FILE_REASK    = "000228ace7127d8c78820e2fa463bc63.xml.0"
	FILE_FIRST    = "000a8adf88587086e242ca0303678cfa.xml.0"
	DIR_SAMPLE    = "../sample/"
	CUR_SAMPLE    = DIR_SAMPLE + FILE_REASK

	ACTION = ACTION_SINGLE

	DIR_MAIN = "u:/yaxml_dat_sample/"
)

// Globals
var db *sql.DB
var CH_NewWords = make(chan string, 100)

type Word struct {
	Word   string
	Strict bool

	Source string // нкря

	// Approximation of the number of documents found for the query.
	// -- Response.Found[0].Found
	Approx uint64

	// Estimated number of groups formed
	// -- Response.Results.Grouping.Found[0].Found
	Groups uint64

	// Approximation of the number of documents found for the query.
	// A more precise estimate compared to the value passed in the found tag for the block with general information about search results.
	// -- Response.Results.Grouping.FoundDocs[0].Found
	Docs uint64
}

//
func main() {
	switch ACTION {
	// ----------------------------------------------------------------------
	case ACTION_SINGLE:
		tree, e := ParseOne(CUR_SAMPLE)
		if e != nil {
			log.Fatalln(e)
		}

		j, e := json.MarshalIndent(tree, "", "\t")
		if e != nil {
			log.Fatalln(e)
		}

		ioutil.WriteFile("json", j, 0777)

		fmt.Printf("Approx: %d\n", tree.Response.Found[0].Found)
		fmt.Printf("Groups: %d\n", tree.Response.Results.Grouping[0].Found[0].Found)
		fmt.Printf("Docs: %d\n", tree.Response.Results.Grouping[0].FoundDocs[0].Found)

		fmt.Println("\n\nSimple Now:\n")
		ParseSimple(CUR_SAMPLE)
	// ----------------------------------------------------------------------
	case ACTION_MASS:
		db, e := sql.Open("mysql", "root:zxcvbnm321@tcp(192.168.185:3306)/rl")
		if e != nil {
			log.Fatalln(e)
		}
		defer db.Close()
	// ----------------------------------------------------------------------
	case ACTION_READDIR:
		dh, e := os.Open(DIR_MAIN)
		if e != nil {
			log.Fatalln(e)
		}
		defer dh.Close()

		for {
			files, e := dh.Readdirnames(1)
			if e != nil {
				break
			}

			filename := files[0]
			if !strings.HasSuffix(filename, ".xml.0") {
				continue
			}

			fmt.Printf("%v %v %v\n", files[0], len(files), e)
		}
	}
	// ----------------------------------------------------------------------

}

// Parse singe file with full parser
func ParseOne(path string) (yaxml.XMLTree, error) {
	file, e := os.Open(path)
	if e != nil {
		log.Fatalln(e)
	}
	defer file.Close()

	fmt.Println("ParseOne()")
	return yaxml.Parse(file)
}

// Parse single file with simplifications
func ParseSimple(path string) (*yaxml.Simple, error) {
	file, e := os.Open(path)
	if e != nil {
		log.Fatalln(e)
	}
	defer file.Close()

	sim := yaxml.New()
	e = sim.AppendReader(file)
	fmt.Println("ParseSimple()")
	return sim, e
}
