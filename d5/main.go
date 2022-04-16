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

type AdventOfCodeDay3Solution struct {
	dataSource AdventOfCodeDataSource
	data       []*Line
}

type Point struct {
	x int
	y int
}

type Line struct {
	from *Point
	to   *Point
}

func (l Line) MaxY() int {
	return max(l.from.y, l.to.y)
}

func (l Line) MaxX() int {
	return max(l.from.x, l.to.x)
}

func (l Line) IsHorizontal() bool {
	return l.from.y == l.to.y
}

func (l Line) IsVertical() bool {
	return l.from.x == l.to.x
}

func toInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func parseLine(lineStr string) *Line {
	lines := strings.Split(lineStr, " -> ")
	point1 := strings.Split(lines[0], ",")
	point2 := strings.Split(lines[1], ",")
	x1 := toInt(point1[0])
	y1 := toInt(point1[1])
	x2 := toInt(point2[0])
	y2 := toInt(point2[1])
	from := &Point{x1, y1}
	to := &Point{x2, y2}
	return &Line{from, to}
}

func parseLines(data []string) []*Line {
	lines := []*Line{}
	for _, lineStr := range data {
		line := parseLine(lineStr)
		lines = append(lines, line)
	}
	return lines
}

func (solution *AdventOfCodeDay3Solution) Data() []*Line {
	if solution.data == nil {
		data, err := solution.dataSource.Read()
		lines := parseLines(data)

		if err != nil {
			log.Fatal(err)
		}

		solution.data = lines
	}

	return solution.data
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

type HydrothermalVentMap struct {
	overlappingVents int
	layout           [][]int
}

func (h HydrothermalVentMap) Print() {
	for _, line := range h.layout {
		fmt.Println(line)
	}
}

func (h HydrothermalVentMap) OverlappingVents() int {
	return h.overlappingVents
}

func initLayout(rows, cols int) [][]int {
	layout := [][]int{}

	for i := 0; i <= rows; i++ {
		row := []int{}
		for j := 0; j <= cols; j++ {
			row = append(row, 0)
		}
		layout = append(layout, row)
	}

	return layout
}

func getNumOfOverlappingVents(lines []*Line, layout [][]int) int {
	overlappingVents := 0

	// this is super inefficient if we don't care about the map layout

	for _, line := range lines {
		currX := line.from.y
		currY := line.from.x

		endX := line.to.y
		endY := line.to.x

		for {
			layout[currX][currY] += 1

			if layout[currX][currY] == 2 {
				overlappingVents += 1
			}

			if currX == endX && currY == endY {
				break
			}

			if currX > endX {
				currX -= 1
			} else if currX < endX {
				currX += 1
			}

			if currY > endY {
				currY -= 1
			} else if currY < endY {
				currY += 1
			}
		}
	}

	return overlappingVents
}

func NewHydrothermalVentMap(lines []*Line, rows int, cols int) *HydrothermalVentMap {
	layout := initLayout(rows, cols)
	overlappingVents := getNumOfOverlappingVents(lines, layout)
	ventmap := &HydrothermalVentMap{overlappingVents, layout}
	return ventmap
}

func (solution AdventOfCodeDay3Solution) Part1() int {
	data := solution.Data()
	horizontalOrVerticalLines := []*Line{}
	rows := 0
	cols := 0

	for _, line := range data {
		if line.IsHorizontal() || line.IsVertical() {
			horizontalOrVerticalLines = append(horizontalOrVerticalLines, line)
			rows = max(rows, line.MaxY())
			cols = max(cols, line.MaxX())
		}
	}

	ventMap := NewHydrothermalVentMap(horizontalOrVerticalLines, rows, cols)
	return ventMap.OverlappingVents()
}

func main() {
	fileDataSource := AdventOfCodeFileDataSourceDay3{"test_data2"}
	solution := AdventOfCodeDay3Solution{fileDataSource, nil}
	fmt.Println(solution.Part1())
}
