package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"runtime"
	"sort"
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
		go func(conn net.Conn) {
			defer conn.Close()
			for {
				message, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					return
				}
				if message != "" {
					synchroniser.setQuery <- string(message)
				}
			}
		}(conn)
	}
}

type sortedMap struct {
	m map[string]int
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	var res bool
	if sm.m[sm.s[i]] == sm.m[sm.s[j]] {
		res = (sm.s[i] < sm.s[j])
	} else {
		res = sm.m[sm.s[i]] > sm.m[sm.s[j]]
	}
	return res
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[string]int) ([]string, int) {
	sm := new(sortedMap)
	count := 0
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		count += m[key]
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s, count
}

// Structure that help synchronize goroutines (show statistic and set words pasts)
type routineSynchroniser struct {
	wordsStat   map[string]int
	lettersStat map[string]int
	setQuery    chan string
	getQuery    chan string
	getRes      chan []map[string]int
}

// Method manage access to statistical info for show and set flows
func (info *routineSynchroniser) SyncStatisticInfo() {
	var message string
	for {
		select {
		case message = <-info.setQuery:
			info.SetInfo(message)
		case <-info.getQuery:
			wStat := map[string]int{}
			lStat := map[string]int{}
			for k, v := range info.wordsStat {
				wStat[k] = v
			}
			for k, v := range info.lettersStat {
				lStat[k] = v
			}
			info.getRes <- []map[string]int{wStat, lStat}
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
		word = strings.ToLower(word)
		if _, ok := info.wordsStat[word]; ok {
			info.wordsStat[word]++
		} else {
			info.wordsStat[word] = 1
		}
		for _, letterRune := range word {
			letter := strings.ToLower(string(letterRune))
			if info.lettersStat[letter] == 0 {
				info.lettersStat[letter] = 1
			} else {
				info.lettersStat[letter] += 1
			}
		}
	}
}

var synchroniser = routineSynchroniser{
	map[string]int{},
	map[string]int{},
	make(chan string, 100000),
	make(chan string),
	make(chan []map[string]int),
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

// Function get number of top letters and words to show and generate
// statistic answer result
func GetStatisticalResult(N int) map[string]interface{} {
	synchroniser.getQuery <- "start"
	statArray := <-synchroniser.getRes
	sortedWords, wordsCount := sortedKeys(statArray[0])
	sortedLetters, _ := sortedKeys(statArray[1])
	topWordsString := fmt.Sprintf("top_%d_words", N)
	topLettersString := fmt.Sprintf("top_%d_letters", N)
	wordsN, lettersN := len(sortedWords), len(sortedLetters)
	if N < wordsN {
		wordsN = N
	}
	if N < lettersN {
		lettersN = N
	}
	topWords, topLetters := sortedWords[:wordsN], sortedLetters[:lettersN]
	result := map[string]interface{}{
		"count":          wordsCount,
		topWordsString:   topWords,
		topLettersString: topLetters,
	}
	return result
}
