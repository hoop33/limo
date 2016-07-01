/* Example on using Snowball stemmer 

This program will read a file, then print "word -> stem(word)" for every word in file
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"bitbucket.org/tebeka/snowball"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s FILENAME\n", os.Args[0])
		flag.PrintDefaults()
	}
	lang := flag.String("lang", "english", "stemmer language")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "error: wrong number of arguments\n")
		os.Exit(1)
	}

	fmt.Println("Using snowball version", snowball.Version)

	stmr, err := snowball.New(*lang)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: can't open %s - %s\n", flag.Arg(0), err)
		os.Exit(1)
	}

	re := regexp.MustCompile("[a-zA-Z]+")

	for _, field := range re.FindAll(data, -1) {
		word := string(bytes.ToLower(field))
		fmt.Printf("%s -> %s\n", word, stmr.Stem(word))
	}
}
