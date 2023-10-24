package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Traveller struct {
	id   int
	x, y int
}

const (
	m = 7
	n = 7
	k = 4
)

const (
	snapshotInterval       = 3 * time.Second
	travellerMoveInterval  = 1 * time.Second
	travellerSpawnInterval = 1 * time.Second
)

const (
	travellerSpawnProbability = 10
	travellerMoveProbability  = 50
)

const (
	North         = 0
	South         = 1
	East          = 2
	West          = 3
	DirectionSize = 4
)

var (
	grid             [m][n]*Traveller
	verticalTraces   [m - 1][n]bool
	horizontalTraces [m][n - 1]bool
	mutex            sync.Mutex
	waitGroup        sync.WaitGroup
)

func getRandomDirection() int {
	return rand.Intn(DirectionSize)
}

func getValueWithProbability(probability int) bool {
	return rand.Intn(100) < probability
}

func createTravellers() {
	for i := 0; i < k; i++ {
		go spawnTraveller(i)
	}
}

func spawnTraveller(i int) {
	for {
		if getValueWithProbability(travellerSpawnProbability) {
			addTraveller(i)
			return
		}
		time.Sleep(travellerSpawnInterval)
	}
}

func addTraveller(id int) {
	mutex.Lock()

	x, y := getEmptySquare()
	traveller := &Traveller{id, x, y}
	grid[x][y] = traveller

	mutex.Unlock()

	go traveller.activate()
}

func getEmptySquare() (int, int) {
	x, y := rand.Intn(m), rand.Intn(n)
	for grid[x][y] != nil {
		x, y = rand.Intn(m), rand.Intn(n)
	}

	return x, y
}

func (traveller *Traveller) activate() {
	for {
		if getValueWithProbability(travellerMoveProbability) {
			traveller.move()
		}
		time.Sleep(travellerMoveInterval)
	}
}

func (traveller *Traveller) move() {
	direction := getRandomDirection()

	switch direction {
	case North:
		traveller.moveNorth()
		break
	case South:
		traveller.moveSouth()
		break
	case East:
		traveller.moveEast()
		break
	case West:
		traveller.moveWest()
		break
	}
}

func (traveller *Traveller) moveNorth() {
	mutex.Lock()

	if traveller.x-1 >= 0 && grid[traveller.x-1][traveller.y] == nil {
		grid[traveller.x-1][traveller.y] = traveller
		grid[traveller.x][traveller.y] = nil
		verticalTraces[traveller.x-1][traveller.y] = true
		traveller.x--
	}

	mutex.Unlock()
}

func (traveller *Traveller) moveSouth() {
	mutex.Lock()

	if traveller.x+1 < m && grid[traveller.x+1][traveller.y] == nil {
		grid[traveller.x+1][traveller.y] = traveller
		grid[traveller.x][traveller.y] = nil
		verticalTraces[traveller.x][traveller.y] = true
		traveller.x++
	}

	mutex.Unlock()
}

func (traveller *Traveller) moveEast() {
	mutex.Lock()

	if traveller.y+1 < n && grid[traveller.x][traveller.y+1] == nil {
		grid[traveller.x][traveller.y+1] = traveller
		grid[traveller.x][traveller.y] = nil
		horizontalTraces[traveller.x][traveller.y] = true
		traveller.y++
	}

	mutex.Unlock()
}

func (traveller *Traveller) moveWest() {
	mutex.Lock()

	if traveller.y-1 >= 0 && grid[traveller.x][traveller.y-1] == nil {
		grid[traveller.x][traveller.y-1] = traveller
		grid[traveller.x][traveller.y] = nil
		horizontalTraces[traveller.x][traveller.y-1] = true
		traveller.y--
	}

	mutex.Unlock()
}

func clearSnapshot() {
	verticalTraces = [m - 1][n]bool{}
	horizontalTraces = [m][n - 1]bool{}
}

func printSnapshot() {
	counter := 1
	for {
		mutex.Lock()

		fmt.Printf("SNAPSHOT - %d\n", counter)
		println(getSnapshot())

		clearSnapshot()
		counter++

		mutex.Unlock()
		time.Sleep(snapshotInterval)
	}
}

func getSnapshot() string {
	snapshot := "\n"
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			snapshot += getSnapshotSquare(i, j)
			snapshot += getSnapshotHorizontalLine(i, j)
		}
		snapshot += "\n"

		if i < m-1 {
			for j := 0; j < n; j++ {
				snapshot += getSnapshotVerticalLine(i, j)
			}
		}
		snapshot += "\n"
	}

	return snapshot
}

func getSnapshotSquare(x, y int) string {
	if grid[x][y] != nil {
		return fmt.Sprintf("%02d", grid[x][y].id)
	}
	return "xx"
}

func getSnapshotHorizontalLine(x, y int) string {
	if y < n-1 && horizontalTraces[x][y] {
		return "|"
	}
	return " "
}

func getSnapshotVerticalLine(x, y int) string {
	if verticalTraces[x][y] {
		return "-- "
	}
	return "   "
}

func main() {
	waitGroup.Add(2)

	go createTravellers()
	go printSnapshot()

	waitGroup.Wait()
}
