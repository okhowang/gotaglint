package main

import (
	"github.com/okhowang/gotaglint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(&gotaglint.JsonNameAnalyzer)
}
