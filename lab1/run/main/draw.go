package main

import (
	"fmt"
	"lexicalanalysis/analyzer"
	"lexicalanalysis/utils"
	"log"
	"os"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func DrawAll(data analyzer.GobData) {
	os.RemoveAll("./output")
	os.Mkdir("./output", os.ModePerm)
	for _, nfa := range data.Nfas {
		if nfa.Name != "Variable" {
			drawGraph(nfa.NFA, "NFA_"+nfa.Name)
		}
	}
	for _, dfa := range data.Dfas {
		if dfa.Name != "Variable" {
			drawGraph(dfa.DFA, "DFA_"+dfa.Name)
		}
	}
}

func drawGraph(item []analyzer.Part, picName string) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()
	var start, end []*cgraph.Node
	var startSet, endSet utils.Set
	var n, m *cgraph.Node
	for _, i := range item {
		if !startSet.Data[i.Src] {
			n, _ = graph.CreateNode(fmt.Sprint(i.Src))
			start = append(start, n)
		} else {
			for _, j := range start {
				if j.Name() == fmt.Sprint(i.Src) {
					n = j
				}
			}
		}
		if !endSet.Data[i.Dst] {
			m, _ = graph.CreateNode(fmt.Sprint(i.Dst))
			end = append(end, m)
		} else {
			for _, j := range start {
				if j.Name() == fmt.Sprint(i.Dst) {
					m = j
				}
			}
		}
		e, _ := graph.CreateEdge(string(i.Edge), n, m)
		e.SetLabel(string(i.Edge))
	}
	if err := g.RenderFilename(graph, graphviz.PNG, fmt.Sprintf("./output/%s.png", picName)); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("output %s done\n", picName)
}
