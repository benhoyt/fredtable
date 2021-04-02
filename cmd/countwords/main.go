// Simple hash table example program (print frequencies of unique words)
//
// To run, execute this command:
// $ go run ./cmd/countwords/main.go <cmd/countwords/test.txt

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/benhoyt/fredtable"
)

type StringKey = fredtable.StringKey

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	counts := fredtable.NewHashTable()
	var uniques []string
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		key := fredtable.StringKey(word)
		pvalue := counts.Get(key)
		if pvalue == nil {
			counts.Set(key, 1)
			uniques = append(uniques, word)
		} else {
			counts.Set(key, (*pvalue).(int)+1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	sort.Slice(uniques, func(i, j int) bool {
		iCount := (*counts.Get(fredtable.StringKey(uniques[i]))).(int)
		jCount := (*counts.Get(fredtable.StringKey(uniques[j]))).(int)
		return iCount > jCount
	})

	for _, word := range uniques {
		count := (*counts.Get(fredtable.StringKey(word))).(int)
		fmt.Println(word, count)
	}

	fmt.Fprintln(os.Stderr, "Hash table dump:")
	fmt.Fprintln(os.Stderr, "----------------")
	counts.Dump(os.Stderr)
}
