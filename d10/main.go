package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

func readLines(reader io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(reader)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

type AdventOfCodeDataSource interface {
	Read() ([]string, error)
}

type AdventOfCodeFileDataSourceDay10 struct {
	filepath string
}

func (dataSource AdventOfCodeFileDataSourceDay10) Read() ([]string, error) {
	f, err := os.Open(dataSource.filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readLines(f)
}

type AdventOfCodeDay10Solution struct {
	dataSource AdventOfCodeDataSource
	data       []string
}

func (solution *AdventOfCodeDay10Solution) Data() []string {
	if solution.data == nil {
		data, err := solution.dataSource.Read()

		if err != nil {
			log.Fatal(err)
		}

		solution.data = data
	}

	return solution.data
}

func getFirstIllegalRune(line string) (rune, []rune) {
	stack := []rune{}

	for _, ch := range line {
		if ch == '(' || ch == '[' || ch == '{' || ch == '<' {
			stack = append(stack, ch)
		} else {
			if len(stack) == 0 {
				return ch, nil
			}

			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if ch == ')' && top == '(' {
				continue
			}

			if ch == ']' && top == '[' {
				continue
			}

			if ch == '}' && top == '{' {
				continue
			}

			if ch == '>' && top == '<' {
				continue
			}

			return ch, nil
		}
	}

	return 0, stack
}

func (solution AdventOfCodeDay10Solution) Part1() int {
	lines := solution.Data()

	/*

		): 3 points
		]: 57 points
		}: 1197 points
		>: 25137 points

	*/

	points := make(map[rune]int)

	points[')'] = 3
	points[']'] = 57
	points['}'] = 1197
	points['>'] = 25137

	total := 0

	for _, line := range lines {
		r, _ := getFirstIllegalRune(line)
		total += points[r]
	}

	return total
}

func (solution AdventOfCodeDay10Solution) Part2() int {
	lines := solution.Data()

	points := make(map[rune]int)

	points['('] = 1
	points['['] = 2
	points['{'] = 3
	points['<'] = 4

	scores := []int{}

	for _, line := range lines {
		_, incompleteLine := getFirstIllegalRune(line)

		if incompleteLine != nil {
			score := 0

			for i := len(incompleteLine) - 1; i >= 0; i-- {
				ch := incompleteLine[i]
				score *= 5
				score += points[ch]
			}

			scores = append(scores, score)
		}
	}

	sort.Ints(scores)

	index := len(scores) / 2

	return scores[index]
}

func main() {
	fileDataSource := AdventOfCodeFileDataSourceDay10{"test_data2"}
	solution := AdventOfCodeDay10Solution{fileDataSource, nil}
	fmt.Println(solution.Part1())
	fmt.Println(solution.Part2())
}
