package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)


type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	filename := flag.String("csv", "problems.csv", "problem file in 'question,answer' format")
	timeLimit := flag.Int("time", 30, "time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file: %s\n", *filename))
	}
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		exit("Failed to parse the CSV file")
	}
	problems := parseLines(lines)
	askQuestions(problems, *timeLimit)

}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: strings.TrimSpace(line[0]),
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return ret

}

func askQuestions(problems []problem, timeLimit int) {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.answer {
				correct++
			}
		}
	}
	fmt.Printf("\nYou scored %d out of %d\n", correct, len(problems))
}

