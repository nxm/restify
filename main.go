package main

import (
	"fmt"
	"net/http"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

var (
	version = "dev"
	commit = "master"
)

var (
	url = kingpin.Arg("url", "A URL to RESTify into JSON").Required().URL()
	byClass = kingpin.Flag("class", "If specified, first-level elements encountered with this class will be extracted.").String()
	byId = kingpin.Flag("id", "If specified, the element with this id will be extracted.").String()
	showVersion = kingpin.Flag("version", "Print version and exit").Bool()
)

func main() {

	kingpin.Parse()

	if *showVersion {
		log.Printf("Version: %s, Commit: %s\n", version, commit)
		os.Exit(0)
	}

	resp, err := http.Get((*url).String())
	if err != nil {
		log.Fatal("Failed to get from URL", err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal("Unable to parse HTML", err)
	}

	var asJson []byte
	var subset []*html.Node
	if *byId != "" {
		elem, ok := scrape.Find(root, scrape.ById(*byId))
		if !ok {
			log.Fatalf("Unable to find an element with the ID '%s'\n", *byId)
		}
		subset = append(subset, elem)
	} else if *byClass != "" {
		subset = scrape.FindAll(root, scrape.ByClass(*byClass))
		if len(subset) == 0 {
			log.Fatalf("Unable to find an element with the class '%s'\n", *byClass)
		}
	} else {
		subset = append(subset, root)
	}
	if asJson, err = convertToJson(subset); err != nil {
		log.Fatal("Failed to parse HTML into JSON", err)
	}

	fmt.Print(string(asJson))
}
