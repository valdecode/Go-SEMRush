package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"encoding/csv"
	"strings"
	"io"
	"net/url"
	"strconv"
	"os"
)

type csvEntry struct {
	url string
	csv []string
}

var scoreMap = map[string]int{}

func main() {
	file := flag.String("f", "", "CSV file with the url in the first column")
	header := flag.Bool("header", false, "First row of file is header")
	apiKey := flag.String("api", "", "API Key given by SEMRush")
	lang := flag.String("lang", "us", "SEMRush language database")
	domain := flag.Bool("domain", false, "From file the first column instead the url is the domain name")
	flag.Parse()
	if *file == "" || *apiKey == "" {
		flag.Usage()
	} else {
		csvEntryList := parseFile(file, *header)
		processUrlList(csvEntryList, apiKey, lang, *domain)
	}
}

// This parses a CSV file and returns an array with every url
func parseFile(file *string, header bool) []csvEntry {
	dat, err := ioutil.ReadFile(*file)
	CheckError(err)
	r := csv.NewReader(strings.NewReader(string(dat)))
	var csvEntryList []csvEntry
	for {
		record, err := r.Read()
		if header {
			fmt.Println(strings.Join(record, ",") + ",Traffic By Semrush")
			header = false
			continue
		}
		if err == io.EOF {
			break
		}
		CheckError(err)
		csvEntryList = append(csvEntryList, csvEntry{record[0], record})
	}
	return csvEntryList
}

// This will use SEMRush API for every url in urlList to get a domain score
func processUrlList(csvEntryList []csvEntry, apiKey *string, lang *string, isDomain bool) {
	w := csv.NewWriter(os.Stdout)
	for i := 0; i < len(csvEntryList); i++ {
		domain := csvEntryList[i].url
		if !isDomain {
			domain = getDomainFromUrl(csvEntryList[i].url)
		}
		score, ok := scoreMap[domain];
		if !ok {
			score = GetDomainScore(domain, apiKey, lang)
			scoreMap[domain] = score
		}
		csvEntryList[i].csv = append(csvEntryList[i].csv, strconv.Itoa(score))
		err := w.Write(csvEntryList[i].csv);
		CheckError(err)
	}
	w.Flush()
	CheckError(w.Error())
}

// Extract the domain name from the url address
func getDomainFromUrl(address string) string {
	u, err := url.Parse(address)
	CheckError(err)
	return u.Host
}

// A common way to treat errors
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
