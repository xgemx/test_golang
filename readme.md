# Go Coding Test

We would like you to build a simple Go application. When started, it
will listen on port 5555 (but this may be configurable through a
command-line flag). Clients will be able to connect to this port and
send arbitrary natural language over the wire. The purpose of the
application is to process the text, and store some stats about the
different words that it sees.

The application will also expose an HTTP interface on port 8080
(configurable): clients hitting the /stats endpoint with an optional query  
string variable N will receive a JSON representation of the statistics about
the words that the application has seen so far.

Specifically, the JSON response for /stats should look like:

```javascript
{
  "count": 42,
  "top_5_words": ["lorem", "ipsum", "dolor", "sit", "amet"],
  "top_5_letters": ["e", "t", "a", "o", "i"]
}
```

Where `count` represents the total number of words seen, `top_5_words`
contains the 5 words that have been seen with the highest frequency, and
`top_5_letters` contains the 5 letters that have been seen with the
highest frequency (you may choose to transform all letters to lowercase
if you so wish).

If N is provided, then its value should be used instead:

```javascript
// /stats?N=3
{
  "count": 42,
  "top_3_words": ["lorem", "ipsum", "dolor"],
  "top_3_letters": ["e", "t", "a"]
}
```



## Things to look out for

* The number of words to process may be large, although you may
  safely assume that they will fit within main memory.

* The application should support a high degree of concurrency, whereby
  many clients would be sending text at the same time.

* We would like to see your approach to automated testing for this type
  of Go program.
