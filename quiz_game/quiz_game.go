package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	problemFile := flag.String("problems", "problems.csv", "Path to the file containing the problems.")
	maxTime := flag.Int("time", 30, "Maximum time to answer questions.")
	doShuffle := flag.Bool("shuffle", false, "Shuffle questions if flag on.")
	flag.Parse()

	csvFile, err := os.Open(*problemFile)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to open the CSV file: %s\n", *problemFile))
		os.Exit(1)
	}

	reader := csv.NewReader(csvFile)
	readIn := bufio.NewReader(os.Stdin)

	allQuestions, _ := reader.ReadAll()

	// Shuffle questions
	if *doShuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(allQuestions),
			func(i, j int) {
				allQuestions[i], allQuestions[j] = allQuestions[j], allQuestions[i]
			})
	}

	fmt.Println("Press enter to start.")
	fmt.Scanln()

	timer := time.NewTimer(time.Duration(*maxTime) * time.Second)
	c := make(chan string)
	result := 0
	maxScore := len(allQuestions)

problemloop:
	for i, q := range allQuestions {
		fmt.Printf("Question %d: %s\n", i+1, q[0])

		go waitAnswer(readIn, c)

		select {
		case answer := <-c:
			if answer == q[1] {
				result++
			}

		case <-timer.C:
			fmt.Println("Time out!")
			break problemloop
		}
	}

	fmt.Printf("You scored %d/%d\n", result, maxScore)
}

func waitAnswer(readIn *bufio.Reader, c chan string) {
	answer, _ := readIn.ReadString('\n')
	answer = strings.TrimSuffix(answer, "\n")

	c <- answer
}
