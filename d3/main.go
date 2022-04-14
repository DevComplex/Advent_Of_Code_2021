package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
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

type AdventOfCodeFileDataSourceDay3 struct {
	filepath string
}

func (dataSource AdventOfCodeFileDataSourceDay3) Read() ([]string, error) {
	f, err := os.Open(dataSource.filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readLines(f)
}

type AdventOfCodeNetworkDataSourceDay3 struct {
	cookieStr string
	url       string
}

func (dataSource AdventOfCodeNetworkDataSourceDay3) Cookies() []*http.Cookie {
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

func (dataSource AdventOfCodeNetworkDataSourceDay3) Read() ([]string, error) {
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

type AdventOfCodeDay3Solution struct {
	dataSource AdventOfCodeDataSource
	data       []string
}

func (solution *AdventOfCodeDay3Solution) Data() []string {
	if solution.data == nil {
		data, err := solution.dataSource.Read()

		if err != nil {
			log.Fatal(err)
		}

		solution.data = data
	}

	return solution.data
}

func binaryToDecimal(binaryStr string) int {
	total := 0.0
	index := 0

	for i := len(binaryStr) - 1; i >= 0; i-- {
		digit := binaryStr[i]

		if digit == '1' {
			total += math.Pow(2, float64(index))
		}

		index += 1
	}

	return int(total)
}

func initMatrix(strs []string) [][]rune {
	matrix := [][]rune{}

	for _, str := range strs {
		matrix = append(matrix, []rune(str))
	}

	return matrix
}

func (solution AdventOfCodeDay3Solution) Part2() int {
	data := solution.Data()
	matrix := initMatrix(data)

	co2ScrubberGroup := []int{}
	oxygenGeneratorGroup := []int{}

	for i := 0; i < len(matrix); i++ {
		co2ScrubberGroup = append(co2ScrubberGroup, i)
		oxygenGeneratorGroup = append(oxygenGeneratorGroup, i)
	}

	col := 0

	for len(co2ScrubberGroup) > 1 {
		co2ScrubberZeroGroup := []int{}
		co2ScrubberOneGroup := []int{}

		for _, row := range co2ScrubberGroup {
			if matrix[row][col] == '0' {
				co2ScrubberZeroGroup = append(co2ScrubberZeroGroup, row)
			} else {
				co2ScrubberOneGroup = append(co2ScrubberOneGroup, row)
			}
		}

		if len(co2ScrubberZeroGroup) <= len(co2ScrubberOneGroup) {
			co2ScrubberGroup = co2ScrubberZeroGroup
		} else {
			co2ScrubberGroup = co2ScrubberOneGroup
		}

		col += 1
	}

	col = 0

	for len(oxygenGeneratorGroup) > 1 {
		oxygenGeneratorZeroGroup := []int{}
		oxygenGeneratorOneGroup := []int{}

		for _, row := range oxygenGeneratorGroup {
			if matrix[row][col] == '0' {
				oxygenGeneratorZeroGroup = append(oxygenGeneratorZeroGroup, row)
			} else {
				oxygenGeneratorOneGroup = append(oxygenGeneratorOneGroup, row)
			}
		}

		if len(oxygenGeneratorOneGroup) >= len(oxygenGeneratorZeroGroup) {
			oxygenGeneratorGroup = oxygenGeneratorOneGroup
		} else {
			oxygenGeneratorGroup = oxygenGeneratorZeroGroup
		}

		col += 1
	}

	oxygenGeneratorGroupRow := oxygenGeneratorGroup[0]
	co2ScrubberGroupRow := co2ScrubberGroup[0]

	oxygenGeneratorBinaryRating := string(matrix[oxygenGeneratorGroupRow])
	co2ScrubberBinaryRating := string(matrix[co2ScrubberGroupRow])

	return binaryToDecimal(oxygenGeneratorBinaryRating) * binaryToDecimal(co2ScrubberBinaryRating)
}

func (solution AdventOfCodeDay3Solution) Part1() int {
	data := solution.Data()

	oneCounts := make(map[int]int)
	zeroCounts := make(map[int]int)

	binaryNumLen := 0

	for _, binaryNum := range data {
		if binaryNumLen == 0 {
			binaryNumLen = len(binaryNum)
		}

		for index, bit := range binaryNum {
			if bit == '0' {
				zeroCounts[index] += 1
			} else {
				oneCounts[index] += 1
			}
		}
	}

	gamma := []rune{}
	epsilon := []rune{}

	for i := 0; i < binaryNumLen; i++ {
		zeroCount := zeroCounts[i]
		oneCount := oneCounts[i]

		if zeroCount > oneCount {
			gamma = append(gamma, '0')
			epsilon = append(epsilon, '1')
		} else {
			gamma = append(gamma, '1')
			epsilon = append(epsilon, '0')
		}
	}

	gammaStr := string(gamma)
	epsilonStr := string(epsilon)

	return binaryToDecimal(gammaStr) * binaryToDecimal(epsilonStr)
}

func main() {
	fileDataSource := AdventOfCodeFileDataSourceDay3{"test_data"}
	solution := AdventOfCodeDay3Solution{fileDataSource, nil}

	fmt.Println(solution.Part1())
	fmt.Println(solution.Part2())
}
