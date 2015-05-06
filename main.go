package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/reiver/go-porterstemmer"
	"net"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// Main function:
// Read ports for show and set stat info from command line flags get_port and set_port.
// And run 2 goroutine, that listen and serve this ports
func main() {
	getPort := flag.String("get_port", "8080", "Set port to part which show stat info with json")
	setPort := flag.String("set_port", "5555", "Set port to part which set strings")
	flag.Parse()
	numcpu := runtime.NumCPU()
	fmt.Println("NumCPU", numcpu)
	runtime.GOMAXPROCS(numcpu)

	go func() {
		http.ListenAndServe(":"+*getPort, &getHandler{})
	}()
	go RunTCPServer(":" + *setPort)
	synchroniser.SyncStatisticInfo()

}

//Run TCP server that listen host (default: localhost:5555) and get messages from the clients
func RunTCPServer(host string) {
	ln, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			continue
		}
		if message != "" {
			synchroniser.setQuery <- string(message)
			conn.Write([]byte("done \n"))
		}
		conn.Close()
	}
}

// Structure that help synchronize goroutines (show statistic and set words pasts)
type routineSynchroniser struct {
	wordsStat   map[string]int
	lettersStat map[string]int
	setQuery    chan string
	getQuery    chan string
	getRes      chan [2]map[string]int
}

// Method manage access to statistical info for show and set flows
func (info *routineSynchroniser) SyncStatisticInfo() {
	var message string
	for {
		select {
		case message = <-info.setQuery:
			info.SetInfo(message)
		case <-info.getQuery:
			info.getRes <- [2]map[string]int{info.wordsStat, info.lettersStat}
		}
	}
}

// Method get string line, that client push to server. Line is split on words.
// Every word is stemmed by The Porter Stemming Algorithm (more info: http://tartarus.org/martin/PorterStemmer/)
// and saved in wordsStat, split on letters and every letter is saved in lettersStat
func (info *routineSynchroniser) SetInfo(line string) {
	rx, err := regexp.Compile("\\W")
	if err != nil {
		fmt.Println(err)
		return
	}
	line = rx.ReplaceAllString(line, " ")
	words := strings.Fields(line)
	for _, word := range words {
		if _, ok := info.wordsStat[porterstemmer.StemString(word)]; ok {
			info.wordsStat[porterstemmer.StemString(word)]++
		} else {
			info.wordsStat[porterstemmer.StemString(word)] = 1
		}
		for _, letterRune := range word {
			letter := string(letterRune)
			if info.lettersStat[letter] == 0 {
				info.lettersStat[letter] = 1
			} else {
				info.lettersStat[letter] += 1
			}
		}
	}
	// fmt.Println("{")
	// for k, v := range info.lettersStat {
	// 	fmt.Printf("\"%s\": %d", k, v)
	// 	fmt.Println(", ")
	// }
	// fmt.Println("}")

}

var synchroniser = routineSynchroniser{
	map[string]int{},
	map[string]int{},
	make(chan string),
	make(chan string),
	make(chan [2]map[string]int),
}

// Structure that serve connection to show static info part of app
type getHandler struct {
}

// Method serve requests for showing stats info
// get in string and process to int number of top letters and words in statistic to show.
// Also generate and show json string of statistic
func (m *getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		N   int
		err error
	)
	if nString := r.FormValue("N"); nString == "" {
		N = 5
	} else {
		N, err = strconv.Atoi(nString)
		if err != nil || N < 1 {
			json.NewEncoder(w).Encode(map[string]string{"status": fmt.Sprintf("Wrong N value")})
			return
		}
	}

	resValue := GetStatisticalResult(N)
	json.NewEncoder(w).Encode(resValue)
}

// map of pair: grammatical symbol(word or letter) and frequency of this symbol
type SymbolPair map[string]int

// Function give top N SymbolPairs and count of all symbols in statistic
func GetTopStat(N int, statInfo map[string]int) (int, []SymbolPair) {
	allCount, j := 0, 0
	popular := make([]SymbolPair, N)
	for symbol, frequency := range statInfo {
		allCount += frequency
		if j < N {
			popular[j] = SymbolPair{symbol: frequency}
			j++
		} else {
			for i := 0; i < N; i++ {
				for buf, bufFrequency := range popular[i] {
					if frequency > bufFrequency {
						popular[i] = SymbolPair{symbol: frequency}
						symbol, frequency = buf, bufFrequency
					}
				}
			}
		}
	}
	return allCount, popular
}

// Function get number of top letters and words to show and generate
// statistic answer result
func GetStatisticalResult(N int) map[string]interface{} {
	synchroniser.getQuery <- "start"
	statArray := <-synchroniser.getRes
	wordsStat, lettersStat := statArray[0], statArray[1]
	topWordsString := fmt.Sprintf("top_%d_words", N)
	topLettersString := fmt.Sprintf("top_%d_letters", N)
	allCount, topWordsPairs := GetTopStat(N, wordsStat)
	_, topLettersPairs := GetTopStat(N, lettersStat)
	topWords, topLetters := []string{}, []string{}
	for i := 0; i < N; i++ {
		for word, _ := range topWordsPairs[i] {
			topWords = append(topWords, word)
		}
		for letter, _ := range topLettersPairs[i] {
			topLetters = append(topLetters, letter)
		}
	}
	result := map[string]interface{}{
		"count":          allCount,
		topWordsString:   topWords,
		topLettersString: topLetters,
	}
	return result
}