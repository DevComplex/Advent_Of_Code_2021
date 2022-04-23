package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readLinesFromFile(filename string) ([]string, error) {
	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func transformLinesToGrid(lines []string) ([][]int, error) {
	grid := [][]int{}

	for _, line := range lines {
		values := strings.Split(line, "")
		row := []int{}

		for _, value := range values {
			valueInt, err := strconv.Atoi(value)

			if err != nil {
				return nil, err
			}

			row = append(row, valueInt)
		}

		grid = append(grid, row)
	}

	return grid, nil
}

func readGridFromFile(filename string) ([][]int, error) {
	lines, err := readLinesFromFile(filename)

	if err != nil {
		return nil, err
	}

	grid, err := transformLinesToGrid(lines)

	if err != nil {
		return nil, err
	}

	return grid, nil
}

type Grid struct {
	data [][]int
}

func (g Grid) String() string {
	lines := []string{}

	for _, row := range g.data {
		line := ""

		for _, val := range row {
			valStr := strconv.Itoa(val)
			line += valStr
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (g Grid) DataCopy() [][]int {
	rows := len(g.data)
	cols := len(g.data[0])

	copy := [][]int{}

	for i := 0; i < rows; i++ {
		row := []int{}
		for j := 0; j < cols; j++ {
			row = append(row, g.data[i][j])
		}
		copy = append(copy, row)
	}

	return copy
}

func (g Grid) Step() (*Grid, int) {
	dataCopy := g.DataCopy()

	directions := [8][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}, {1, 1}, {-1, -1}, {1, -1}, {-1, 1}}
	queue := [][2]int{}
	totalFlashes := 0

	for i := 0; i < len(dataCopy); i++ {
		for j := 0; j < len(dataCopy[i]); j++ {
			dataCopy[i][j] += 1

			if dataCopy[i][j] > 9 {
				totalFlashes += 1
				queue = append(queue, [2]int{i, j})
				dataCopy[i][j] = 0
			}
		}
	}

	for len(queue) > 0 {
		coord := queue[0]
		queue = queue[1:]
		row := coord[0]
		col := coord[1]

		for _, dir := range directions {
			nextRow := row + dir[0]
			nextCol := col + dir[1]

			if nextRow < 0 || nextRow == len(dataCopy) {
				continue
			}

			if nextCol < 0 || nextCol == len(dataCopy[0]) {
				continue
			}

			if dataCopy[nextRow][nextCol] == 0 {
				continue
			}

			dataCopy[nextRow][nextCol] += 1

			if dataCopy[nextRow][nextCol] > 9 {
				totalFlashes += 1
				queue = append(queue, [2]int{nextRow, nextCol})
				dataCopy[nextRow][nextCol] = 0
			}
		}
	}

	return &Grid{dataCopy}, totalFlashes
}

type FlashSimulation struct {
	steps []*Grid
}

func (simulation FlashSimulation) Print() {
	for _, step := range simulation.steps {
		fmt.Println(step.String())
		fmt.Println()
	}
}

func NewFlashSimulation(filename string) (*FlashSimulation, error) {
	grid, err := readGridFromFile(filename)
	if err != nil {
		return nil, err
	}
	steps := []*Grid{{grid}}
	return &FlashSimulation{steps}, nil
}

func (simulation *FlashSimulation) Simulate(steps int) int {
	step := simulation.steps[0]
	totalFlashes := 0

	for i := 0; i < steps; i++ {
		s, flashesFromStep := step.Step()
		simulation.steps = append(simulation.steps, s)
		totalFlashes += flashesFromStep
		step = s
	}

	return totalFlashes
}

func (simulation *FlashSimulation) FirstStepWithAllFlash() int {
	step := simulation.steps[0]

	rows := len(step.data)
	cols := len(step.data[0])
	targetFlashes := rows * cols

	i := 1

	for {
		s, flashesFromStep := step.Step()

		if flashesFromStep == targetFlashes {
			return i
		}

		step = s
		i += 1
	}
}

func main() {
	filename := "test_data2"
	simulation, err := NewFlashSimulation(filename)

	if err != nil {
		log.Fatal(err)
	}

	flashes := simulation.Simulate(100)
	fmt.Println(flashes)

	step := simulation.FirstStepWithAllFlash()
	fmt.Println(step)
}
