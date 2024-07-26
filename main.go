package main

import (
        "bufio"
        "flag"
        "fmt"
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
        fmt.Println("  go run main.go [-h] [-f <hyphen count>] <filename>")
        fmt.Println()
        fmt.Println("Flags:")
        fmt.Println("  -h    Show this help message and exit")
        fmt.Println("  -f    Filter subdomains by hyphen count")
        fmt.Println()
        fmt.Println("Run Modes:")
        fmt.Println("  Simple Run: go run main.go <filename>")
        fmt.Println("    Provides detailed information about the file, including")
        fmt.Println("    the minimum and maximum number of hyphens, and example subdomains.")
        fmt.Println()
        fmt.Println("  Filtered Run: go run main.go -f <hyphen count> <filename>")
        fmt.Println("    Provides a list of subdomains with the specified number of hyphens.")
        fmt.Println("github.com/rezauditore/hyphenlu")
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

        if len(flag.Args()) < 1 {
                fmt.Println("Usage: go run main.go [-h] [-f <hyphen count>] <filename>")
                return
        }
        filename := flag.Arg(0)
        file, err := os.Open(filename)
        if err != nil {
                fmt.Printf("Error opening file: %v\n", err)
                return
        }
        defer file.Close()

        hyphenCountMap := make(map[int][]string)
        var mu sync.Mutex
        var wg sync.WaitGroup

        subdomains := make(chan string, 100)

        numWorkers := 8
        for i := 0; i < numWorkers; i++ {
                wg.Add(1)
                go worker(subdomains, hyphenCountMap, &mu, &wg)
        }

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                subdomains <- scanner.Text()
        }
        close(subdomains)

        if err := scanner.Err(); err != nil {
                fmt.Printf("Error reading file: %v\n", err)
                return
        }

        wg.Wait()

        if filter != -1 {
                if _, ok := hyphenCountMap[filter]; ok {
                        for _, subdomain := range hyphenCountMap[filter] {
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
        fmt.Print("[INFO] Example Subdomains: ")
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
        fmt.Println("\n[*] Shutting Down at here.")
}
