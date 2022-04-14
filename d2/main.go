package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Direction string

const (
	Forward Direction = "forward"
	Up                = "up"
	Down              = "down"
)

type SubmarineVector struct {
	direction Direction
	value     int
}

func (submarineVector SubmarineVector) String() string {
	return string(submarineVector.direction) + " " + strconv.Itoa(submarineVector.value)
}

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

func parseSubmarineVectors(lines []string) ([]*SubmarineVector, error) {
	submarineVectors := []*SubmarineVector{}
	for _, line := range lines {
		directionValue := strings.Split(line, " ")
		direction := directionValue[0]
		value, err := strconv.Atoi(directionValue[1])
		if err != nil {
			return nil, err
		}
		submarineVectors = append(submarineVectors, &SubmarineVector{Direction(direction), value})
	}
	return submarineVectors, nil
}

type AdventOfCodeDataSource interface {
	Read() ([]string, error)
}

type AdventOfCodeFileDataSourceDay2 struct {
	filepath string
}

func (dataSource AdventOfCodeFileDataSourceDay2) Read() ([]string, error) {
	f, err := os.Open(dataSource.filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readLines(f)
}

type AdventOfCodeNetworkDataSourceDay2 struct {
	cookieStr string
	url       string
}

func (dataSource AdventOfCodeNetworkDataSourceDay2) Cookies() []*http.Cookie {
	cookieStrs := strings.Split(dataSource.cookieStr, "; ")
	cookies := []*http.Cookie{}
	for _, cookieStr := range cookieStrs {
		keyValue := strings.Split(cookieStr, "=")
		key := keyValue[0]
		value := keyValue[1]
		cookies = append(cookies, &http.Cookie{Name: key, Value: value})
	}
	return cookies
}

func (dataSource AdventOfCodeNetworkDataSourceDay2) Read() ([]string, error) {
	client := http.Client{}
	cookies := dataSource.Cookies()
	req, err := http.NewRequest("GET", dataSource.url, nil)
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return readLines(resp.Body)
}

type AdventOfCodeDay2Solution struct {
	dataSource AdventOfCodeDataSource
	data       []*SubmarineVector
}

func (solution *AdventOfCodeDay2Solution) Data() ([]*SubmarineVector, error) {
	if solution.data != nil {
		return solution.data, nil
	}
	lines, err := solution.dataSource.Read()
	if err != nil {
		return nil, err
	}
	submarineVectors, err := parseSubmarineVectors(lines)
	if err != nil {
		return nil, err
	}
	solution.data = submarineVectors
	return submarineVectors, nil
}

func (solution AdventOfCodeDay2Solution) Part1() int {
	submarineVectors, err := solution.Data()

	if err != nil {
		log.Fatal(err)
	}

	horizontalPosition := 0
	depth := 0

	for _, submarineVector := range submarineVectors {
		direction := submarineVector.direction
		value := submarineVector.value

		if direction == Forward {
			horizontalPosition += value
		} else if direction == Up {
			depth -= value
		} else if direction == Down {
			depth += value
		}
	}

	return horizontalPosition * depth
}

func (solution AdventOfCodeDay2Solution) Part2() int {
	submarineVectors, err := solution.Data()

	if err != nil {
		log.Fatal(err)
	}

	aim := 0
	horizontalPosition := 0
	depth := 0

	for _, submarineVector := range submarineVectors {
		direction := submarineVector.direction
		value := submarineVector.value

		if direction == Forward {
			horizontalPosition += value
			depth += (aim * value)
		} else if direction == Up {
			aim -= value
		} else if direction == Down {
			aim += value
		}
	}

	return horizontalPosition * depth
}

func NewAdventOfCodeDay2Solution(url string, cookieStr string) AdventOfCodeDay2Solution {
	networkDataSource := AdventOfCodeNetworkDataSourceDay2{cookieStr, url}
	solution := AdventOfCodeDay2Solution{networkDataSource, nil}
	return solution
}

func main() {
	fileDataSource := AdventOfCodeFileDataSourceDay2{"test_data"}
	solution := AdventOfCodeDay2Solution{fileDataSource, nil}
	fmt.Println(solution.Part1())
	fmt.Println(solution.Part2())
}
