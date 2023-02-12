// simple tool for deduping files
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

func ShowVersion() {
	ver := "0.1.0"
	fmt.Println("dedupe v", ver)
	os.Exit(0)
}

func ShowUsage(code int) {
	fmt.Println("show usage here")
	os.Exit(code)
}

func countLines(f *os.File, counts map[string]int) int {
	n := 0
	input := bufio.NewScanner(f)
	for input.Scan() {
		n++
		counts[input.Text()]++
	}

	return n
}

type stats struct {
	fileName     string
	dupeCount    int
	maxDupeCount int
	dupeLines    int
	minLen       int
	maxLen       int
	meanLen      float64
	lengthSum    int
	totalLines   int
	elapsed      time.Duration
}

func newStats() *stats {
	s := stats{"", 0, 0, 0, 1000000, 0, 0, 0, 0, 0}
	return &s
}

func CollectStats(counts *map[string]int, s *stats) {
	for line, n := range *counts {
		if n > 1 {
			s.dupeLines++
			s.dupeCount += (n - 1)

			if n > s.maxDupeCount {
				s.maxDupeCount = (n - 1)
			}

			length := len(line)
			s.lengthSum += length
			if length > s.maxLen {
				s.maxLen = length
			}
			if length < s.minLen {
				s.minLen = length
			}
		}

	}
}

func PresentStats(s *stats) {
	//s.meanLen = float64(s.lengthSum) / float64(s.dupeLines)
	pct := (float64(s.dupeCount) / float64(s.totalLines)) * 100

	fmt.Fprintf(os.Stdout, "%8d     %8d         %3.1f%%      %8d      %s\n",
		s.totalLines, s.dupeCount, pct, s.maxDupeCount, s.fileName)
}

func main() {
	var (
		files   []string
		help    = flag.Bool("help", false, "Show usage")
		version = flag.Bool("version", false, "Show version")
	)

	flag.Parse()

	if *help {
		ShowUsage(0)
	}

	if *version {
		ShowVersion()
	}

	if flag.NArg() > 0 {
		files = os.Args[1:]
	}

	if len(files) == 0 {
		ShowUsage(0)
	} else {

		files := os.Args[1:]
		fmt.Println("   LINES     DUPE-LINES     PCT-DUPE     DUPE-MAX      FILE")

		// foreach file

		for _, file := range files {

			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%8v\n", err)
				continue
			}
			defer f.Close()

			s := newStats()
			s.fileName = file

			counts := make(map[string]int)
			t0 := time.Now()
			s.totalLines = countLines(f, counts)
			t1 := time.Now()
			s.elapsed = t1.Sub(t0)

			CollectStats(&counts, s)
			PresentStats(s)
		}
	}
}
