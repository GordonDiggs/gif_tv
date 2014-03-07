package main

import (
	"bufio"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var urls []string

const INDEX = `
<html>
<head><script src="http://code.jquery.com/jquery-1.11.0.min.js"></script>
<script type="text/javascript">
$(function() {
  setTimeout(function() {
    var num = parseInt(window.location.href.replace(/.*num=/, ''), 10) || 0;
    num += 1;
    window.location.href = '/?num='+num;
  }, 2000);
});
</script>
</head>
<body style="background-color:black;">
<img src="{{.Url}}" style="width:100%;height:100%;">
</body>
</html>
`

type Gif struct {
	Url string
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	urls = lines
}

func handler(w http.ResponseWriter, r *http.Request) {
	num, _ := strconv.ParseInt(r.FormValue("num"), 10, 0)

	if int(num) > len(urls) {
		num = int64(rand.Intn(len(urls)))
	}

	var gif_url = urls[num]

	var p = new(Gif)
	p.Url = gif_url
	t := template.New("gif_page")
	t, _ = t.Parse(INDEX)
	t.Execute(w, p)
}

func main() {
	readLines("gifs.txt")

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", handler)
	fmt.Println("Waiting for requests...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
