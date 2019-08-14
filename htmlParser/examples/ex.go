package main

import (
	"fmt"
	//"strings"
	"bufio"
	"flag"
	"os"

	"gophercises/htmlParser"
)

func main() {
	//r := strings.NewReader(htmlString)
	filePath := flag.String("file", "ex1.html", "HTML file to parse")
	flag.Parse()
	fmt.Println(*filePath)

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("Error while opening file: ", err)
	}
	r := bufio.NewReader(file)

	links, err := htmlParser.ParseHtml(r)
	if err != nil {
		fmt.Println("Error while parsing: ", err)
	}

	fmt.Printf("%+v\n", links)
}

var htmlString = `
<html>
<head>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
</head>
<body>
  <h1>Social stuffs</h1>
  <div>
    <a href="https://www.twitter.com/joncalhoun">
      Check me out on twitter
      <i class="fa fa-twitter" aria-hidden="true"></i>
    </a>
    <a href="https://github.com/gophercises">
      Gophercises is on <strong>Github</strong>!
    </a>
  </div>
</body>
</html>
`
