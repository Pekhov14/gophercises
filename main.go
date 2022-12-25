package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "format of question - answer")
	timeLimit := flag.Int("limit", 30, "The time limit for quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed open the CSV file: %s\n", *csvFileName))
	}
	//_ = file // Заглушка использования переменной

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

	// Ярлык цикла
problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)

		answerCh := make(chan string) // Канал с ответами
		// Процедура go, вызываем анонимную func
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemLoop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}
