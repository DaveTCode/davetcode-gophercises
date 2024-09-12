package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type question struct {
	q string
	a string
}

func parseCsvFile(filePath string) []question {
	quiz_file, err := os.Open(filePath)

	if err != nil {
		log.Fatal("Can't open file "+filePath, err)
	}

	defer quiz_file.Close()

	csvReader := csv.NewReader(quiz_file)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Can't parse "+filePath, err)
	}

	questions := make([]question, 0)

	for _, record := range records {
		if len(record) != 2 {
			log.Fatal("Invalid csv record", record)
		} else {
			questions = append(questions, question{q: record[0], a: record[1]})
		}
	}

	return questions
}

func askQuestion(question question, ch chan bool) {
	fmt.Printf("%s: ", question.q)
	var input string
	fmt.Scanln(&input)

	ch <- strings.TrimSpace(input) == question.a
}

var (
	quizFileName     string
	shuffle          bool
	timeLimitSeconds int
)

func init() {
	flag.StringVar(&quizFileName, "quiz", "quiz.csv", "The csv file containing one question per line")
	flag.BoolVar(&shuffle, "shuffle", false, "Shuffle the questions in the quiz file")
	flag.IntVar(&timeLimitSeconds, "timeLimit", 30, "The time limit in seconds")

	flag.Parse()
}

func main() {

	questions := parseCsvFile(quizFileName)

	if shuffle {
		rand.Shuffle(len(questions), func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
		})
	}

	fmt.Println("Press Enter To Start")
	var input string
	fmt.Scanln(&input)

	timer := time.NewTimer(time.Duration(timeLimitSeconds) * time.Second)

	ch := make(chan bool)

	correct := 0

Quiz:
	for _, question := range questions {
		go askQuestion(question, ch)
		select {
		case ans := <-ch:
			if ans {
				correct++
			}
		case <-timer.C:
			break Quiz
		}
	}

	fmt.Println("======")
	fmt.Printf("%d/%d Correct\n", correct, len(questions))
}
