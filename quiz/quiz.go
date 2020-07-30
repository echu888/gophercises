// implementation of Gophercise Exercise #1 
//   https://github.com/gophercises/quiz/
//   
// primary features:
//	read from CSV file
//	channel timer
//	command line arguments handled by flags library

package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"io"
	"log"
	"time"
)

type Problem struct {
	Problem string
	ExpectedAnswer string
	ActualAnswer string
	Answered bool
	Correct bool
}



func readFile(filename string) []Problem {
	file, _ := os.Open(filename)
	reader  := csv.NewReader(bufio.NewReader(file))
	var problems []Problem
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		problems = append(problems, Problem{
			Problem:         line[0],
			ExpectedAnswer:  line[1],
			Answered:        false,
			Correct:         false,
		})
	}
	return problems

}

func askQuestions(problems []Problem) {
	for i:=0; i < len(problems); i++ {
		fmt.Printf("Question #%d: %s = ", i, problems[i].Problem)
		var input string
		fmt.Scanln(&input)
		if input != "" {
			problems[i].Answered     = true
			problems[i].ActualAnswer = input
		}
		problems[i].Correct      = (problems[i].ActualAnswer == problems[i].ExpectedAnswer)
	}


}

func printResults(problems []Problem) {
	fmt.Println()
	fmt.Println("Quiz results:")

	var totalAnswered = 0
	var totalCorrect = 0
	//for _, problem := range problems {
	for i, problem := range problems {
		fmt.Printf("Question #%d: %s = %s, Answered:%s, Correct:%t\n",
			i, problem.Problem, problem.ExpectedAnswer, problem.ActualAnswer, problem.Correct)

		if problem.Answered {
			totalAnswered++
		}
		if problem.Correct {
			totalCorrect++
		}
	}

	fmt.Println()
	fmt.Println("Quiz score:")
	fmt.Printf(" -- Answered %d out of %d total questions\n", totalAnswered, len(problems))
	fmt.Printf(" -- Correct %d\n", totalCorrect)
	fmt.Printf(" -- Incorrect %d\n", totalAnswered - totalCorrect)
}


func main() {
	const PROBLEMS_FILENAME = "problems.csv"
	const TIMES_UP = 30
	//var programArgs = os.Args[1:]

	filenamePtr  := flag.String("csv", PROBLEMS_FILENAME, "CSV file with problems")
	timeLimitPtr := flag.Int("limit", TIMES_UP, "time limit in seconds")
	flag.Parse()


	fmt.Println()
	fmt.Println("Welcome to quiz!")
	fmt.Println("Timer is set to", *timeLimitPtr, "seconds")
	fmt.Println("Reading quiz set from", *filenamePtr)
	fmt.Println("-------------------------------")
	problems := readFile(*filenamePtr)
	fmt.Println("Press enter when you're ready to start...")
	fmt.Scanln()


	channel := make(chan int)
	go func() {
		askQuestions(problems[:])
		channel <- 1
	}()

	select {
		case <-channel:
			printResults(problems)

		case <-time.After(time.Duration(*timeLimitPtr) * time.Second):
			fmt.Println("Times up!")
			printResults(problems)
	}

}


