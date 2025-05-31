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

type Question struct {
	question      string
	correctAnswer string
}

func readCSV(filenName string) [][]string {
	log.Println("Opening CSV file from filepath", filenName)
	f, err := os.Open(filenName)
	if err != nil {
		log.Fatal("An error occured", err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("An error occured", err)
	}
	return records
}

func timer(seconds int, c chan int) {
	time.Sleep(time.Second * time.Duration(seconds))
	c <- -1
}

func getQuestionsFromCSV(fileName string) []Question {
	records := readCSV(fileName)

	var questions []Question

	for _, record := range records {
		question := Question{record[0], record[1]}
		questions = append(questions, question)

	}
	return questions
}

func askQuestion(question Question, c chan int) {
	var userAnswer int

	fmt.Printf("%s? :", question.question)
	fmt.Scan(&userAnswer)
	fmt.Print("\n")
	userAnswerString := strconv.Itoa(userAnswer)

	if userAnswerString == question.correctAnswer {
		fmt.Println("Answer is Correct!")
		c <- 1
	} else {
		fmt.Println("Answer is Incorrect :(")
		c <- 0
	}
}

func main() {
	csvFileName := flag.String("fileName", "problems.csv", "file name string")
	timeLimit := flag.Int("timeLimit", 30, "time in seconds")
	flag.Parse()
	questions := getQuestionsFromCSV(*csvFileName)

	score := 0
	c := make(chan int)
	go timer(*timeLimit, c)
	for _, question := range questions {
		go askQuestion(question, c)
		event := <-c

		if event == -1 {
			fmt.Printf("\n\nOops! Time is up!\n Final Score: %d\n", score)
			os.Exit(0)
		} else {
			score += event
			fmt.Printf("Score: %d\n", score)
		}
	}

}
