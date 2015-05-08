package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sort"
	"testing"
	"time"
)

func TestOverload(t *testing.T) {
	type testSet struct {
		numOfTests  int
		oneIter     int
		resultCount int
		runTime     int
		tests       []string
		words       []string
		letters     []string
	}
	tests := []string{
		"Introduction to Go The latest Go release, version, arrives as scheduled six months after. It contains only one tiny language change, in the form of a backwards-compatible simple variant of for-range loop, and a possibly breaking change to the compiler involving methods on pointers-to-pointers. The release focuses primarily on implementation work, improving the garbage collector and preparing the ground for a fully concurrent collector to be rolled out in the next few releases. Stacks are now contiguous, reallocated when necessary rather than linking on new \"segments\"; this release therefore eliminates the notorious \"hot stack split\" problem. There are some new tools available including support in the go command for build-time source code generation. The release also adds support for ARM processors on Android and Native Client (NaCl) and for AMD64 on Plan 9. As always, Go keeps the promise of compatibility, and almost everything will continue to compile and run without change when moved to. Changes to the language For-range loops Up until Go, for-range loop had two forms for i, v := range x { ... } and for i := range x { ... } If one was not interested in the loop values, only the iteration itself, it was still necessary to mention a variable (probably the blank identifier, as in for _ = range x), because the form for range x  was not syntactically permitted. This situation seemed awkward, so as of Go the variable-free form is now legal. The pattern arises rarely but the code can be cleaner when it does. Updating: The change is strictly backwards compatible to existing Go programs, but tools that analyze Go parse trees may need to be modified to accept this new form as the Key field of RangeStmt may now be nil. Method calls on **T Given these declarations, type T int func (T) M() {} var x **T both gc and gccgo accepted the method call x.M() which is a double dereference of the pointer-to-pointer x. The Go specification allows a single dereference to be inserted automatically, but not two, so this call is erroneous according to the language definition. It has therefore been disallowed in Go, which is a breaking change, although very few programs will be affected. Updating: Code that depends on the old, erroneous behavior will no longer compile but is easy to fix by adding an explicit dereference. Changes to the supported operating systems and architectures Android Go 1.4 can build binaries for ARM processors running the Android operating system. It can also build a .so library that can be loaded by an Android application using the supporting packages in the mobile subrepository. A brief description of the plans for this experimental port are available here. NaCl on ARM The previous release introduced Native Client (NaCl) support for the 32-bit x86 (GOARCH=386) and 64-bit x86 using 32-bit pointers (=amd64p32). The 1.4 release adds NaCl support for ARM (GOARCH=arm). Plan9 on AMD64 This release adds support for the Plan 9 operating system on AMD64 processors, provided the kernel supports the nsec system call and uses 4K pages. Changes to the compatibility guidelines The unsafe package allows one to defeat Go's type system by exploiting internal details of the implementation or machine representation of data. It was never explicitly specified what use of unsafe meant with respect to compatibility as specified in the Go compatibility guidelines. The answer, of course, is that we can make no promise of compatibility for code that does unsafe things. We have clarified this situation in the documentation included in the release. The Go compatibility guidelines and the docs for the unsafe package are now explicit that unsafe code is not guaranteed to remain compatible. Updating: Nothing technical has changed; this is just a clarification of the documentation. Changes to the implementations and tools Changes to the runtime Prior to Go 1.4, the runtime (garbage collector, concurrency support, interface management, maps, slices, strings) was mostly written in C, with some assembler support. In Go, much of the code has been translated to Go so that the garbage collector can scan the stacks of programs in the runtime and get accurate information about what variables are active. This change was large but should have no semantic effect on programs. This rewrite allows the garbage collector in to be fully precise, meaning that it is aware of the location of all active pointers in the program. This means the heap will be smaller as there will be no false positives keeping non-pointers alive. Other related changes also reduce the heap size, which is smaller by 10-30 overall relative to the previous release. A consequence is that stacks are no longer segmented, eliminating the \"hot split\" problem. When a stack limit is reached, a new, larger stack is allocated, all active frames for the goroutine are copied there, and any pointers into the stack are updated. Performance can be noticeably better in some cases and is always more predictable. Details are available in the design document. The use of contiguous stacks means that stacks can start smaller without triggering performance issues, so the default starting size for a goroutine's stack in 1.4 has been reduced from bytes to 2048 bytes. As preparation for the concurrent garbage collector scheduled for the 1.5 release, writes to pointer values in the heap are now done by a function call, called a write barrier, rather than directly from the function updating the value. In this next release, this will permit the garbage collector to mediate writes to the heap while it is running. This change has no semantic effect on programs in, but was included in the release to test the compiler and the resulting performance. The implementation of interface values has been modified. In earlier releases, the interface contained a word that was either a pointer or a one-word scalar value, depending on the type of the concrete object stored",
	}
	letters := []string{"e", "t", "a", "o", "i"}
	words := []string{"the", "to", "in", "for", "of"}
	set := testSet{1000, 1000, 1000000, 20, tests, words, letters}
	startTime := time.Now()
	for i := 0; i < set.numOfTests/set.oneIter; i++ {
		for j := 0; j < set.oneIter; j++ {
			for _, test := range set.tests {
				go SendString(test)
			}
		}
	}
	timeDelta := int(time.Since(startTime) / 1000000000)
	for timeDelta < set.runTime {
		fmt.Println("Wait", set.runTime-timeDelta)
		time.Sleep(10000000000)
		timeDelta = int(time.Since(startTime) / 1000000000)
	}
	result := GetResult(5, 0, t)
	if int(result["count"].(float64)) != set.resultCount {
		t.Errorf("Test result: count expect %d, got %d", set.resultCount, int(result["count"].(float64)))
	}
	for _, topElement := range result[fmt.Sprintf("top_%d_words", 5)].([]interface{}) {
		sort.Strings(set.words)
		if sort.SearchStrings(set.words, topElement.(string)) == 0 && set.words[0] != topElement.(string) {
			t.Errorf("Can't find '%s' in %s", topElement, set.words)
		}

	}
	for _, topElement := range result[fmt.Sprintf("top_%d_letters", 5)].([]interface{}) {
		sort.Strings(set.letters)
		if sort.SearchStrings(set.letters, topElement.(string)) == 0 && set.letters[0] != topElement.(string) {
			t.Errorf("Can't find '%s' in %s", topElement, set.letters)
		}

	}
}

func SendString(test string) {
	conn, err := net.Dial("tcp", "127.0.0.1:5555")
	if err != nil {
		log.Panic("Fuck", err)
		return
	}
	defer conn.Close()
	n, err := fmt.Fprintln(conn, test)
	if err != nil || n != 5964 {
		log.Panic("Fuck", err, n)
		return
	}
}

func GetResult(N int, i int, t *testing.T) map[string]interface{} {
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
	return mapBody
}
