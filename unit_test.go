package main

import (
	"testing"
)

func TestGetTopStat(t *testing.T) {
	type testSet struct {
		N             int
		statInfo      map[string]int
		count         int
		minCountToTop int
	}
	tests := [...]testSet{
		{5, map[string]int{"aaa": 3, "BBB": 5, "CCC": 50}, 58, 3},
		{2, map[string]int{"aaa": 3, "BBB": 5, "CCC": 50}, 58, 5},
		{5, lettersStat, 1267, 86},
		{5, wordsStat, 271, 5},
	}
	for i, setForTest := range tests {
		allCount, top := GetTopStat(setForTest.N, setForTest.statInfo)
		if allCount != setForTest.count {
			t.Errorf("TestGetTopStat #%d Count expect %d, got %d", i, setForTest.count, allCount)
		}
		var minTopCount int
		for _, res := range top {
			for _, frequency := range res {
				if minTopCount == 0 {
					minTopCount = frequency
				} else if minTopCount > frequency {
					minTopCount = frequency
				}
			}
		}
		if minTopCount != setForTest.minCountToTop {
			t.Errorf("TestGetTopStat #%d Minimum count of element in top expect %d, got %d", i, setForTest.minCountToTop, minTopCount)
		}
	}
}

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
				"A": 6,
				"B": 7,
				"C": 6,
				"D": 5,
				"F": 4,
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
			make(chan [2]map[string]int),
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
	"d": 33,
	"T": 8,
	"t": 139,
	"l": 61,
	"4": 2,
	"e": 142,
	"u": 39,
	"m": 25,
	"5": 10,
	"b": 19,
	"r": 71,
	"J": 2,
	"j": 2,
	"2": 2,
	"C": 2,
	"s": 89,
	"k": 5,
	"c": 34,
	"S": 3,
	"3": 3,
	"p": 42,
	"h": 57,
	"x": 3,
	"H": 1,
	"P": 1,
	"8": 2,
	"0": 2,
	"i": 78,
	"y": 21,
	"a": 86,
	"q": 3,
	"I": 1,
	"n": 67,
	"f": 24,
	"v": 9,
	"N": 5,
	"G": 3,
	"o": 110,
	"g": 19,
	"W": 4,
	"w": 24,
	"O": 2,
	"_": 12,
}
var wordsStat = map[string]int{
	"a":          6,
	"support":    1,
	"http":       1,
	"autom":      1,
	"type":       1,
	"would":      3,
	"purpos":     1,
	"also":       1,
	"seen":       4,
	"look":       2,
	"amet":       1,
	"frequenc":   2,
	"your":       1,
	"natur":      1,
	"json":       2,
	"transform":  1,
	"you":        4,
	"to":         10,
	"it":         3,
	"but":        1,
	"string":     1,
	"count":      3,
	"lorem":      2,
	"e":          2,
	"test":       2,
	"when":       1,
	"some":       1,
	"top_5_word": 2,
	"total":      1,
	"abl":        1,
	"about":      2,
	"mani":       1,
	"highest":    2,
	"valu":       1,
	"thei":       1,
	"concurr":    1,
	"simpl":      1,
	"that":       5,
	"i":          1,
	"all":        1,
	"main":       1,
	"time":       1,
	"we":         2,
	"n":          3,
	"approach":   1,
	"and":        3,
	"top_5_lett": 2,
	"o":          1,
	"then":       1,
	"stat":       4,
	"sit":        1,
	"out":        1,
	"although":   1,
	"program":    1,
	"where":      1,
	"safe":       1,
	"at":         1,
	"5555":       1,
	"word":       5,
	"dolor":      2,
	"high":       1,
	"far":        1,
	"3":          1,
	"memori":     1,
	"degre":      1,
	"build":      1,
	"flag":       1,
	"receiv":     1,
	"thing":      1,
	"code":       1,
	"applic":     5,
	"configur":   2,
	"see":        2,
	"number":     2,
	"interfac":   1,
	"specif":     1,
	"choos":      1,
	"us":         1,
	"like":       3,
	"an":         2,
	"hit":        1,
	"ha":         1,
	"should":     3,
	"if":         2,
	"go":         3,
	"line":       1,
	"ipsum":      2,
	"t":          2,
	"been":       2,
	"lowercas":   1,
	"provid":     1,
	"top_3_word": 1,
	"thi":        3,
	"wire":       1,
	"of":         6,
	"endpoint":   1,
	"5":          2,
	"top_3_lett": 1,
	"larg":       1,
	"wherebi":    1,
	"listen":     1,
	"mai":        4,
	"through":    1,
	"command":    1,
	"client":     3,
	"process":    2,
	"wish":       1,
	"queri":      1,
	"42":         2,
	"have":       2,
	"same":       1,
	"be":         5,
	"the":        19,
	"for":        3,
	"connect":    1,
	"8080":       1,
	"with":       3,
	"statist":    1,
	"store":      1,
	"expos":      1,
	"repres":     1,
	"within":     1,
	"over":       1,
	"represent":  1,
	"respons":    1,
	"contain":    2,
	"on":         2,
	"languag":    1,
	"option":     1,
	"assum":      1,
	"fit":        1,
	"will":       5,
	"arbitrari":  1,
	"is":         2,
	"differ":     1,
	"letter":     2,
	"instead":    1,
	"port":       3,
	"text":       2,
	"variabl":    1,
	"javascript": 2,
	"start":      1,
	"send":       2,
	"so":         2,
}
