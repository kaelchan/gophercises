package main

/* The idea of timeout is borrowed from https://gobyexample.com/timeouts
 */

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// Problem contains both problem and answer
type Problem struct {
	description string
	answer      int
}

// return a set of problems from a csv file
func readCSV(filename string) ([]Problem, error) {
	// get the reader
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(file)

	// read rows, convert, and append
	ret := make([]Problem, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		answer, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}
		problem := Problem{
			description: record[0],
			answer:      answer,
		}
		ret = append(ret, problem)
	}
	return ret, nil
}

func displayResult(points int, problemCount int) {
	var correctRate float64
	if problemCount == 0 {
		correctRate = 0
	} else {
		correctRate = 100 * float64(points) / float64(problemCount)
	}
	fmt.Printf("Correct Rate: %v / %v = %f%%\n",
		points, problemCount, correctRate)
}

func examination(
	points *int, problemSet []Problem,
	finishNotification chan bool) {
	for index, problem := range problemSet {
		// query the user until a correct format is received
		var (
			answer int
			err    error
		)
		for {
			var stringAnswer string
			fmt.Printf("Problem # %v: %v, please enter your answer: ",
				index, problem.description)
			fmt.Scanln(&stringAnswer)

			answer, err = strconv.Atoi(stringAnswer)
			if err != nil {
				// redo if the format isn't correct
				fmt.Println("Please enter a integer")
				continue
			}
			break
		}
		// compare with gold truth
		if problem.answer == answer {
			*points++
		}
	}
	finishNotification <- true
}

func main() {
	filename := flag.String("filename", "problems.csv", "the input file name")
	timeLimit := flag.Int("time-limit", 30, "time limit in seconds")
	// Once all flags are declared, call flag.Parse() to execute the
	// command-line parsing.
	flag.Parse()
	log.Println(*filename)
	log.Println(*timeLimit)

	problemSet, err := readCSV(*filename)
	if err != nil {
		fmt.Println("Error while reading the csv")
		return
	}
	problemCount := len(problemSet)

	points := 0
	finishNotification := make(chan bool, 1)
	go examination(&points, problemSet, finishNotification)

	select {
	case <-finishNotification:
		fmt.Println("\n--------- Finished! ---------")
	case <-time.After(time.Duration(*timeLimit) * time.Second):
		fmt.Println("\n********* Time Out! *********")
	}

	displayResult(points, problemCount)
}
