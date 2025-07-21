package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
)

func countHyphens(s string) int {
	return strings.Count(s, "-")
}

func worker(subdomains chan string, hyphenCountMap map[int][]string, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for subdomain := range subdomains {
		count := countHyphens(subdomain)
		mu.Lock()
		hyphenCountMap[count] = append(hyphenCountMap[count], subdomain)
		mu.Unlock()
	}
}

func printHelp() {
	fmt.Println("Hyphen Detector Tool (Version 1)")
	fmt.Println("--------------------------------")
	fmt.Println("This tool detects hyphens in a list of subdomains.")
	fmt.Println("It is useful for making things easy to analyze.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  cat <filename> | go run main.go [-h] [-f <hyphen count>]")
	fmt.Println("  go run main.go [-h] [-f <hyphen count>] <filename>")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -h    Show this help message and exit")
	fmt.Println("  -f    Filter subdomains by hyphen count")
	fmt.Println()
	fmt.Println("GitHub: github.com/rezauditore/hyphenlu")
}

func readInput(r io.Reader, subdomains chan string) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		subdomains <- scanner.Text()
	}
	return scanner.Err()
}

func main() {
	var filter int
	var help bool
	flag.IntVar(&filter, "f", -1, "Filter subdomains by hyphen count")
	flag.BoolVar(&help, "h", false, "Show this help message and exit")
	flag.Parse()

	if help {
		printHelp()
		return
	}

	var input io.Reader
	if len(flag.Args()) > 0 {
		filename := flag.Arg(0)
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	} else {
		// Check if piped input exists
		stat, err := os.Stdin.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		if stat.Mode()&os.ModeCharDevice != 0 {
			fmt.Fprintln(os.Stderr, "Usage: cat <filename> | go run main.go [-f <hyphen count>]")
			fmt.Fprintln(os.Stderr, "       or: go run main.go [-f <hyphen count>] <filename>")
			os.Exit(1)
		}
		input = os.Stdin
	}

	hyphenCountMap := make(map[int][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	subdomains := make(chan string, 100)

	numWorkers := 8
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(subdomains, hyphenCountMap, &mu, &wg)
	}

	err := readInput(input, subdomains)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	close(subdomains)

	wg.Wait()

	if filter != -1 {
		if list, ok := hyphenCountMap[filter]; ok {
			for _, subdomain := range list {
				fmt.Println(subdomain)
			}
		} else {
			fmt.Printf("No subdomains found with %d hyphens.\n", filter)
		}
		return
	}

	var hyphenCounts []int
	for k := range hyphenCountMap {
		hyphenCounts = append(hyphenCounts, k)
	}
	sort.Ints(hyphenCounts)

	if len(hyphenCounts) == 0 {
		fmt.Println("No subdomains found.")
		return
	}

	fmt.Println("---------------------------------------------------")
	fmt.Printf("[INFO] Minimum number of hyphens: %d\n", hyphenCounts[0])
	fmt.Print("[INFO] Example Subdomains:\n")
	for i, subdomain := range hyphenCountMap[hyphenCounts[0]] {
		fmt.Println(subdomain)
		if i >= 5 {
			break
		}
	}

	maxHyphens := hyphenCounts[len(hyphenCounts)-1]
	fmt.Printf("\n[INFO] Maximum number of hyphens: %d\n", maxHyphens)
	fmt.Println("[INFO] Example Subdomains:")
	for i, subdomain := range hyphenCountMap[maxHyphens] {
		fmt.Println(subdomain)
		if i >= 5 {
			break
		}
	}
	fmt.Println("---------------------------------------------------")
}
