// A test application for applying brotli compression to files
package main

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/kothar/brotli-go.v0/dec"
	"gopkg.in/kothar/brotli-go.v0/enc"
)

// Flags
var (
	compress   string
	decompress string
	output     string
	quality    int
)

func main() {
	// Configure flags
	flag.StringVar(&compress, "c", "", "compress a file")
	flag.StringVar(&decompress, "d", "", "decompress a file")
	flag.StringVar(&output, "o", "", "output file")
	flag.IntVar(&quality, "q", 9, "compression quality (1-11)")
	flag.Parse()

	// Read input
	var input string
	if compress != "" {
		input = compress
	} else if decompress != "" {
		input = decompress
	} else {
		log.Fatal("You must specify either compress or decompress")
	}

	inputData, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}

	// Perform compression or decompression
	var outputData []byte
	if compress != "" {
		params := enc.NewBrotliParams()
		params.SetQuality(quality)
		outputData, err = enc.CompressBuffer(params, inputData, nil)
	} else if decompress != "" {
		outputData, err = dec.DecompressBuffer(inputData, nil)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Write output
	if output == "" {
		if compress != "" {
			output = input + ".bro"
		} else if decompress != "" {
			output = input + ".unbro"
		}
	}
	err = ioutil.WriteFile(output, outputData, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
