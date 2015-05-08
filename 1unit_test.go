package main

import (
	"testing"
)

// func AlreadyNotUnitTestGetTopStat(t *testing.T) {
// 	type testSet struct {
// 		N             int
// 		statInfo      map[string]int
// 		count         int
// 		minCountToTop int
// 	}
// 	tests := [...]testSet{
// 		{5, map[string]int{"aaa": 3, "BBB": 5, "CCC": 50}, 58, 3},
// 		{2, map[string]int{"aaa": 3, "BBB": 5, "CCC": 50}, 58, 5},
// 		{5, lettersStat, 1267, 86},
// 		{5, wordsStat, 271, 5},
// 	}
// 	for i, setForTest := range tests {
// 		allCount, top := GetTopStat(setForTest.N, setForTest.statInfo)
// 		if allCount != setForTest.count {
// 			t.Errorf("TestGetTopStat #%d Count expect %d, got %d", i, setForTest.count, allCount)
// 		}
// 		var minTopCount int
// 		for _, res := range top {
// 			for _, frequency := range res {
// 				if minTopCount == 0 {
// 					minTopCount = frequency
// 				} else if minTopCount > frequency {
// 					minTopCount = frequency
// 				}
// 			}
// 		}
// 		if minTopCount != setForTest.minCountToTop {
// 			t.Errorf("TestGetTopStat #%d Minimum count of element in top expect %d, got %d", i, setForTest.minCountToTop, minTopCount)
// 		}
// 	}
// }

func TestSetInfo(t *testing.T) {
	type testSet struct {
		line          string
		resultWords   map[string]int
		resultLetters map[string]int
	}
	tests := [...]testSet{
		{
			"AAA BBB CCC DDD FFF AAA BBBB CCC DD F",
			map[string]int{
				"aaa":  2,
				"bbb":  1,
				"ccc":  2,
				"ddd":  1,
				"fff":  1,
				"bbbb": 1,
				"dd":   1,
				"f":    1,
			},
			map[string]int{
				"a": 6,
				"b": 7,
				"c": 6,
				"d": 5,
				"f": 4,
			},
		},
		{line, wordsStat, lettersStat},
	}
	for i, setForTest := range tests {
		syncho := routineSynchroniser{
			map[string]int{},
			map[string]int{},
			make(chan string),
			make(chan string),
			make(chan []map[string]int),
		}
		syncho.SetInfo(setForTest.line)
		for key, val := range setForTest.resultLetters {
			if syncho.lettersStat[key] != val {
				t.Errorf("TestSetInfo #%d  lettersStat on key '%s' expect %d, got %d", i, key, val, syncho.lettersStat[key])
			}
		}
		for key, val := range setForTest.resultWords {
			if syncho.wordsStat[key] != val {
				t.Errorf("TestSetInfo #%d  wordsStat on key '%s' expect %d, got %d", i, key, val, syncho.wordsStat[key])
			}
		}
	}
}

var line = "Go Coding Test We would like you to build a simple Go application. When started, it will listen on port 5555 (but this may be configurable through a command-line flag). Clients will be able to connect to this port and send arbitrary natural language over the wire. The purpose of the application is to process the text, and store some stats about the different words that it sees. The application will also expose an HTTP interface on port 8080 (configurable): clients hitting the /stats endpoint with an optional query string variable N will receive a JSON representation of the statistics about the words that the application has seen so far. Specifically, the JSON response for /stats should look like: ```javascript { \"count\": 42, \"top_5_words\": [\"lorem\", \"ipsum\", \"dolor\", \"sit\", \"amet\"], \"top_5_letters\": [\"e\", \"t\", \"a\", \"o\", \"i\"] } ``` Where `count` represents the total number of words seen, `top_5_words` contains the 5 words that have been seen with the highest frequency, and `top_5_letters` contains the 5 letters that have been seen with the highest frequency (you may choose to transform all letters to lowercase if you so wish). If N is provided, then its value should be used instead: ```javascript // /stats?N=3 { \"count\": 42, \"top_3_words\": [\"lorem\", \"ipsum\", \"dolor\"], \"top_3_letters\": [\"e\", \"t\", \"a\"] } ``` ## Things to look out for * The number of words to process may be large, although you may safely assume that they will fit within main memory. * The application should support a high degree of concurrency, whereby many clients would be sending text at the same time. * We would like to see your approach to automated testing for this type of Go program."

var lettersStat = map[string]int{
	"c": 36,
	"s": 92,
	"k": 5,
	"3": 3,
	"b": 19,
	"r": 71,
	"j": 4,
	"2": 2,
	"i": 79,
	"y": 21,
	"a": 86,
	"q": 3,
	"p": 43,
	"h": 58,
	"x": 3,
	"8": 2,
	"0": 2,
	"g": 22,
	"o": 112,
	"w": 28,
	"_": 12,
	"n": 72,
	"f": 24,
	"v": 9,
	"e": 142,
	"u": 39,
	"m": 25,
	"5": 10,
	"d": 33,
	"t": 147,
	"l": 61,
	"4": 2,
}
var wordsStat = map[string]int{
	"listen":         1,
	"may":            4,
	"about":          2,
	"3":              1,
	"coding":         1,
	"flag":           1,
	"language":       1,
	"safely":         1,
	"they":           1,
	"testing":        1,
	"that":           5,
	"represents":     1,
	"main":           1,
	"see":            1,
	"we":             2,
	"expose":         1,
	"n":              3,
	"an":             2,
	"top_5_letters":  2,
	"o":              1,
	"letters":        2,
	"provided":       1,
	"clients":        3,
	"ipsum":          2,
	"sit":            1,
	"out":            1,
	"program":        1,
	"you":            4,
	"it":             2,
	"but":            1,
	"interface":      1,
	"string":         1,
	"count":          3,
	"where":          1,
	"concurrency":    1,
	"test":           1,
	"when":           1,
	"able":           1,
	"arbitrary":      1,
	"is":             2,
	"some":           1,
	"variable":       1,
	"total":          1,
	"degree":         1,
	"natural":        1,
	"process":        2,
	"javascript":     2,
	"query":          1,
	"i":              1,
	"top_3_letters":  1,
	"support":        1,
	"time":           1,
	"http":           1,
	"its":            1,
	"type":           1,
	"would":          3,
	"then":           1,
	"your":           1,
	"json":           2,
	"transform":      1,
	"to":             10,
	"configurable":   2,
	"endpoint":       1,
	"e":              2,
	"choose":         1,
	"fit":            1,
	"will":           5,
	"sees":           1,
	"dolor":          2,
	"instead":        1,
	"started":        1,
	"command":        1,
	"text":           2,
	"hitting":        1,
	"value":          1,
	"many":           1,
	"send":           1,
	"so":             2,
	"have":           2,
	"highest":        2,
	"same":           1,
	"a":              6,
	"stats":          4,
	"top_3_words":    1,
	"this":           3,
	"connect":        1,
	"purpose":        1,
	"8080":           1,
	"with":           3,
	"store":          1,
	"words":          5,
	"seen":           4,
	"look":           2,
	"top_5_words":    2,
	"amet":           1,
	"lowercase":      1,
	"if":             2,
	"simple":         1,
	"line":           1,
	"been":           2,
	"on":             2,
	"representation": 1,
	"lorem":          2,
	"application":    5,
	"5555":           1,
	"frequency":      2,
	"sending":        1,
	"port":           3,
	"through":        1,
	"receive":        1,
	"far":            1,
	"wish":           1,
	"memory":         1,
	"build":          1,
	"has":            1,
	"42":             2,
	"assume":         1,
	"be":             5,
	"the":            19,
	"for":            3,
	"number":         2,
	"contains":       2,
	"all":            1,
	"used":           1,
	"large":          1,
	"optional":       1,
	"approach":       1,
	"like":           3,
	"and":            3,
	"different":      1,
	"also":           1,
	"statistics":     1,
	"response":       1,
	"should":         3,
	"within":         1,
	"go":             3,
	"over":           1,
	"specifically":   1,
	"t":              2,
	"although":       1,
	"wire":           1,
	"of":             6,
	"5":              2,
	"things":         1,
	"at":             1,
	"high":           1,
	"whereby":        1,
	"automated":      1,
}
