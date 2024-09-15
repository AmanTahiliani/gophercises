// Authored by Aman Tahiliani
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Initial Algorithm design by Aman

// Read and Parse CSV
// Initialize Score variables
// Initialize Loop
// Print out Queston On to the terminal
// One Channel to Count time and return time over when done.
// One Channel to take in Input and return it.
// We either get the answer back or a timeout.
// If answer, then check for the correctness of the answer and update the score accordingly.
var score int = 0

func readCsv(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	fmt.Println(records)
	return records, nil
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func scanAnswer(scanChan chan int) {
	var i int
	fmt.Scan(&i)
	scanChan <- i
	close(scanChan)
}

func timesUp(n int, timeChan chan string) {
	time.Sleep(time.Duration(n) * time.Second)
	timeChan <- "Times Up!! Moving on."
}

func checkAnswer(correctAnswer int, inputAnswer int) {
	if inputAnswer == correctAnswer {
		score++
		fmt.Println("Correct Answer!")
	} else {
		fmt.Println("Wrong Answer :(")
	}
	fmt.Println("")
}

func main() {
	log.Println("Starting the Quiz")

	csvFileName := flag.String("csv", "quiz.csv", "A CSV file in the format of  'question, answer'")
	flag.Parse()

	records, err := readCsv(*csvFileName)

	if err != nil {
		exit(fmt.Sprintf("Error occurred while reading records: %s", err.Error()))
	}
	fmt.Println("Number of records: ", len(records))
	score = 0

	for index, record := range records[1:] {
		// startTime := time.Now().Format(time.RFC822)
		question := record[0]
		answer := record[1]

		fmt.Println("Question Number", index+1)
		fmt.Println(question, "?")

		answer_int, _ := strconv.Atoi(answer)

		scanChan := make(chan int)
		timeChan := make(chan string)

		go scanAnswer(scanChan)
		go timesUp(7, timeChan)

		select {
		case msg1 := <-scanChan:
			checkAnswer(answer_int, msg1)
		case msg2 := <-timeChan:
			fmt.Println(msg2)
		}

	}
	fmt.Println("Your score is :", score)
}
