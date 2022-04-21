package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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

type AdventOfCodeFileDataSourceDay9 struct {
	filepath string
}

func (dataSource AdventOfCodeFileDataSourceDay9) Read() ([]string, error) {
	f, err := os.Open(dataSource.filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readLines(f)
}

type AdventOfCodeDay9Solution struct {
	dataSource AdventOfCodeDataSource
	data       [][]int
}

func parseGrid(data []string) [][]int {
	grid := [][]int{}

	for _, row := range data {
		rowValStrs := strings.Split(row, "")
		rowVals := []int{}
		for _, valStr := range rowValStrs {
			val, err := strconv.Atoi(valStr)
			if err != nil {
				log.Fatal(err)
			}
			rowVals = append(rowVals, val)
		}
		grid = append(grid, rowVals)

	}
	return grid
}

func (solution *AdventOfCodeDay9Solution) Data() [][]int {
	if solution.data == nil {
		data, err := solution.dataSource.Read()
		lines := parseGrid(data)

		if err != nil {
			log.Fatal(err)
		}

		solution.data = lines
	}

	return solution.data
}

func (solution AdventOfCodeDay9Solution) Part1() int {
	data := solution.Data()

	directions := [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	total := 0

	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			val := data[i][j]

			isLowPoint := true

			for _, direction := range directions {
				nextRow := i + direction[0]
				nextCol := j + direction[1]

				if nextRow < 0 || nextCol < 0 {
					continue
				}

				if nextRow == len(data) || nextCol == len(data[0]) {
					continue
				}

				if val >= data[nextRow][nextCol] {
					isLowPoint = false
					break
				}
			}

			if isLowPoint {
				total += (val + 1)
			}
		}
	}

	return total
}

func (solution AdventOfCodeDay9Solution) Part2() int {
	data := solution.Data()

	directions := [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	total := 1

	largestBasins := [3]int{0, 0, 0}

	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			if data[i][j] != -1 {
				val := data[i][j]

				isLowPoint := true

				for _, direction := range directions {
					nextRow := i + direction[0]
					nextCol := j + direction[1]

					if nextRow < 0 || nextCol < 0 {
						continue
					}

					if nextRow == len(data) || nextCol == len(data[0]) {
						continue
					}

					if val >= data[nextRow][nextCol] {
						isLowPoint = false
						break
					}
				}

				if isLowPoint {
					queue := [][3]int{{i, j, val}}
					count := 0

					for len(queue) > 0 {
						count += 1
						row := queue[0][0]
						col := queue[0][1]
						val := queue[0][2]
						queue = queue[1:]

						for _, direction := range directions {
							nextRow := row + direction[0]
							nextCol := col + direction[1]

							if nextRow < 0 || nextCol < 0 {
								continue
							}

							if nextRow == len(data) || nextCol == len(data[0]) {
								continue
							}

							nextVal := data[nextRow][nextCol]

							if nextVal != 9 && nextVal != -1 && nextVal > val {
								queue = append(queue, [3]int{nextRow, nextCol, nextVal})
								data[nextRow][nextCol] = -1
							}
						}
					}

					if count > largestBasins[0] {
						largestBasins[2] = largestBasins[1]
						largestBasins[1] = largestBasins[0]
						largestBasins[0] = count
					} else if count > largestBasins[1] {
						largestBasins[2] = largestBasins[1]
						largestBasins[1] = count
					} else if count > largestBasins[2] {
						largestBasins[2] = count
					}
				}
			}
		}
	}

	for _, val := range largestBasins {
		if val != 0 {
			total *= val
		}
	}

	return total
}

func main() {
	fileDataSource := AdventOfCodeFileDataSourceDay9{"test_data2"}
	solution := AdventOfCodeDay9Solution{fileDataSource, nil}
	fmt.Println(solution.Part1())
	fmt.Println(solution.Part2())
}
