package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"strings"

	"lexicalanalysis/analyzer"
	"syntacticanalysis/LL1"
)

func decode(Data analyzer.GobData) {
	file, _ := os.Create("gob.data")
	enc := gob.NewEncoder(file)
	enc.Encode(Data)
}
func encode() analyzer.GobData {
	file, _ := os.Open("gob.data")
	dec := gob.NewDecoder(file)
	data := analyzer.GobData{}
	dec.Decode(&data)
	return data
}

func main() {
	var data analyzer.GobData
	var key string
	fmt.Printf("input 1 to generate serialized data or any to skip: ")
	fmt.Scanf("%s\n", &key)
	if key == "1" {
		data := analyzer.GetNFADFA()
		decode(data)
	}
	data = encode()
	for {
		// test := "tg(-123.5) * (_a5)\n"
		fmt.Printf("input expression to match: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\r')
		text = strings.ReplaceAll(text, "\r", "")
		text = strings.ReplaceAll(text, "\\n", "\n")
		text = strings.ReplaceAll(text, "\\t", "\t")
		keys, tmps := analyzer.Match(text, data)
		if keys != nil {
			if LL1.Analyzer(keys, tmps) {
				fmt.Println("syntactic analysis: parse success")
			} else {
				fmt.Println("syntactic analysis: parse fail")
			}
			fmt.Println()
		}
	}
}
