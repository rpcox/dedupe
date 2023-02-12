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
	dupeCount    int
	maxDupeCount int
	dupeLines    int
	minLen       int
	maxLen       int
	meanLen      float64
	lengthSum    int
}

func newStats() *stats {
	s := stats{0, 0, 0, 1000000, 0, 0, 0}
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

func PresentStats(s *stats, totalLines int, t0 time.Time, t1 time.Time) {
	s.meanLen = float64(s.lengthSum) / float64(s.dupeLines)
	pct := (float64(s.dupeCount) / float64(totalLines)) * 100

	fmt.Fprintf(os.Stderr, "\n%10d\tTotal lines read in\n", totalLines)
	fmt.Fprintf(os.Stderr, "%10d\tDupe lines (%.2f%%)\n", s.dupeCount, pct)

	if s.maxDupeCount > 0 {
		fmt.Fprintf(os.Stderr, "%10d\tMax line dupes\n", s.maxDupeCount)
		fmt.Fprintf(os.Stderr, "%10d\tMin dupe line length (bytes)\n", s.minLen)
		fmt.Fprintf(os.Stderr, "%10d\tMax dupe line length (bytes)\n", s.maxLen)
		fmt.Fprintf(os.Stderr, "%10.2f\tMean dupe line length (bytes)\n", s.meanLen)
	}

	fmt.Fprintf(os.Stderr, "\n\telepsed: %v\n\n", t1.Sub(t0))
	t0 = t1
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

		// foreach file
		for _, file := range files {

			fmt.Fprintf(os.Stderr, "***  %v\n", file)

			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			defer f.Close()

			counts := make(map[string]int)
			t0 := time.Now()
			totalLines := countLines(f, counts)
			t1 := time.Now()
			s := newStats()

			CollectStats(&counts, s)
			PresentStats(s, totalLines, t0, t1)
		}
	}
}
