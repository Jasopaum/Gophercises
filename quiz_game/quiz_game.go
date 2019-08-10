package main

import (
    "fmt"
    "encoding/csv"
    "os"
    "bufio"
    "strings"
    "flag"
    "time"
    "math/rand"
)


func main() {
    problemFile := flag.String("problems", "problems.csv", "Path to the file containing the problems.")
    maxTime := flag.Int("time", 30, "Maximum time to answer questions.")
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
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(allQuestions),
                func(i, j int) {
                    allQuestions[i], allQuestions[j] = allQuestions[j], allQuestions[i]
                })

    fmt.Println("Press enter to start.")
    fmt.Scanln()

    timer := time.NewTimer(time.Duration(*maxTime) * time.Second)
    c := make(chan string)
    result := 0
    maxScore := len(allQuestions)

    problemloop:
    for i, q := range allQuestions {
        fmt.Printf("Question %d: %s\n", i, q[0])

        go waitAnswer(readIn, c)

        select {
            case answer := <-c:
                if answer == q[1] {
                    result++
                }

            case <-timer.C:
                fmt.Println()
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
