package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

//dataMap here is to help remove duplicates
var m map[string]int

/*getData scrapes given url, matches to a regex and puts answers within a map*/
func getData(url string, dataMap map[string]int) int {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR ON RESPONSE WITH : " + url)
	}

	dataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	hrefMatch, _ := regexp.Compile("(https://[a-z|A-Z]*.[a-z|A-Z]*.[a-z|A-Z]*)/[a-z|A-Z]*(/?)")
	stringOfData := string(dataBytes[:])

	found := hrefMatch.FindAllString(stringOfData, -1)
	if len(found) <= 0 {
		fmt.Println("no matches found on current..." + url)
	}

	addToSet(found, dataMap)
	resp.Body.Close()
	return (len(found))
}

/*addToSet is taking in a []string and then mapping said data to a var*/
func addToSet(data []string, mData map[string]int) {
	for i := 0; i < len(data); i++ {
		if mData[data[i]] <= 1 {
			mData[data[i]] = mData[data[i]] + 1
		} else {
			mData[data[i]] = 1
		}
	}
}

func main() {
	m = make(map[string]int)
	read := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter starting site...")
	fmt.Println("Format as https://www.example.com/")
	text, _ := read.ReadString('\n')
	//a URL to start the process off with
	getData(text, m)

	for i := 0; i < 3; i++ {
		for key := range m {
			getData(key, m)
		}
	}

	file, err := os.OpenFile("matches.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	for key := range m {
		file.WriteString("\n" + key + "\n")
	}
	file.Close()
}
