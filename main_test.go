package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"testing"
)

type possibleResult struct {
	count   int
	words   []string
	letters []string
}
type testAllSet struct {
	insideLine string
	result     map[int]possibleResult
}

func TestAllProcess(t *testing.T) {
	go main()
	tests := [...]testAllSet{
		testAllSet{
			"My name is John Doe",
			map[int]possibleResult{
				3: possibleResult{
					5,
					[]string{"my", "name", "is", "john", "doe"},
					[]string{"n", "e", "o"},
				},
				5: possibleResult{
					5,
					[]string{"my", "name", "is", "john", "doe"},
					[]string{"M", "y", "n", "a", "m", "e", "i", "s", "J", "o", "h", "D", "e"},
				},
				6: possibleResult{
					5,
					[]string{"my", "name", "is", "john", "doe"},
					[]string{"M", "y", "n", "a", "m", "e", "i", "s", "J", "o", "h", "D", "e"},
				},
			},
		},
		testAllSet{
			"Hello John! My name is Joan Johnes",
			map[int]possibleResult{
				1: possibleResult{
					12,
					[]string{"john"},
					[]string{"n", "o"},
				},
				3: possibleResult{
					12,
					[]string{"my", "name", "is", "john"},
					[]string{"n", "e", "o"},
				},
				5: possibleResult{
					12,
					[]string{"my", "name", "is", "john", "doe", "hello", "joan"},
					[]string{"M", "y", "n", "a", "m", "e", "i", "s", "J", "o", "h", "D", "e"},
				},
			},
		},
		testAllSet{
			"Hello Joan Johnes!",
			map[int]possibleResult{
				5: possibleResult{
					15,
					[]string{"my", "name", "is", "john", "doe"},
					[]string{"n", "e", "o", "J", "a"},
				},
				7: possibleResult{
					15,
					[]string{"my", "name", "is", "john", "doe"},
					[]string{"M", "y", "n", "a", "m", "e", "i", "s", "J", "o", "h", "D", "e"},
				},
			},
		},
	}
	for i, test := range tests {
		RunTest(i, test, t)
	}
}

func RunTest(i int, test testAllSet, t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:5555")
	for err != nil {
		conn, err = net.Dial("tcp", "127.0.0.1:5555")
	}
	defer conn.Close()
	fmt.Fprintln(conn, test.insideLine)
	bufio.NewReader(conn).ReadString('\n')
	for N, res := range test.result {
		request := fmt.Sprintf("http://localhost:8080/?N=%d", N)
		resp, err := http.Get(request)
		if err != nil {
			t.Errorf("Test #%d, N=%d: HTTP connect error: %s", i, N, err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Test #%d, N=%d: HTTP read body error: %s", i, N, err)
		}
		mapBody := map[string]interface{}{}
		json.Unmarshal(body, &mapBody)
		if int(mapBody["count"].(float64)) != res.count {
			t.Errorf("Test #%d, N=%d: count expect %g, got %d", i, N, mapBody["count"], res.count)
		}
		topWordsString := fmt.Sprintf("top_%d_words", N)
		topLettersString := fmt.Sprintf("top_%d_letters", N)
		for _, topElement := range mapBody[topWordsString].([]interface{}) {
			sort.Strings(res.words)
			if sort.SearchStrings(res.words, topElement.(string)) == 0 && res.words[0] != topElement.(string) {
				t.Errorf("Test #%d, N=%d: can't find '%s' in %s", i, N, topElement, res.words)
			}
		}
		for _, topElement := range mapBody[topLettersString].([]interface{}) {
			sort.Strings(res.letters)
			if sort.SearchStrings(res.letters, topElement.(string)) == 0 && res.letters[0] != topElement.(string) {
				t.Errorf("Test #%d, N=%d: can't find '%s' in %s", i, N, topElement, res.letters)
			}
		}
	}

}
