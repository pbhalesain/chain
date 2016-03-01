package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"chain/fedchain/bc"
)

const help = `
Command decode reads a data item from stdin, decodes it,
and prints its JSON representation to stdout.

On Mac OS X, to decode an item from the pasteboard,

	pbpaste|decode tx
	pbpaste|decode block
	pbpaste|decode blockheader
`

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func prettyPrint(obj interface{}) {
	j, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fatalf("error json-marshaling: %s", err)
	}
	fmt.Println(string(j))
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, help)
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println(strings.TrimSpace(help))
		return
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatalf("%v", err)
	}

	switch strings.ToLower(args[0]) {
	case "blockheader":
		b := make([]byte, len(data)/2)
		_, err := hex.Decode(b, data)
		if err != nil {
			fatalf("err decoding hex: %s", err)
		}

		var bh bc.BlockHeader
		err = bh.Scan(b)
		if err != nil {
			fatalf("error decoding: %s", err)
		}

		// The struct doesn't have the hash, so calculate it and print it
		// before pretty printing the header.
		fmt.Printf("Block Hash: %s\n", bh.Hash())
		prettyPrint(bh)
	case "block":
		b := make([]byte, len(data)/2)
		_, err := hex.Decode(b, data)
		if err != nil {
			fatalf("err decoding hex: %s", err)
		}

		var block bc.Block
		err = block.Scan(b)
		if err != nil {
			fatalf("error decoding: %s", err)
		}

		// The struct doesn't have the hash, so calculate it and print it
		// before pretty printing the block
		fmt.Printf("Block Hash: %s\n", block.Hash())
		prettyPrint(block)
	case "tx":
		var tx bc.Tx
		err := tx.UnmarshalText(data)
		if err != nil {
			fatalf("error decoding: %s", err)
		}
		prettyPrint(tx)
	default:
		fatalf("unrecognized entity `%s`", args[0])
	}
}