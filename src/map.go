package main

import "fmt"
import "sync"
//import "math"


import(
  //  "os"
    //"bufio"
	"encoding/json"
	//"time"
    "io/ioutil"
)

//import "time"

const MINHEAT = 10
const MEDIUMHEAT = 20
const MAXHEAT = 30

type tile struct {
	xCoord int
	yCoord int

	heat      int //how hot a tile is before fire
	fireLevel int //strength of the fire

	wall bool
	door bool

	occupied *Person
	personID int

	outOfBounds bool

	neighborNorth *tile
	neighborEast  *tile
	neighborSouth *tile
	neighborWest  *tile

	neighborNW *tile
	neighborNE *tile
	neighborSE *tile
	neighborSW *tile
}

//Initializes the fire
func SetFire(thisTile *tile) {
	thisTile.heat = MINHEAT
	thisTile.fireLevel = 1
}

func FireSpread(tileMap [][]tile) {
	for x := 0; x < len(tileMap); x++ {
		for y := 0; y < len(tileMap[0]); y++ {
			tl := &(tileMap[x][y])
			// tl.heat < 30 {fireSpreadTile(tl)}
			fireSpreadTile(tl)
		}
	}

}

func fireSpreadTile(thisTile *tile) { //TODO: fixa tiles på riktigt!!
	if thisTile.heat >= MINHEAT {
		thisTile.fireLevel = 1
	}
	if thisTile.heat >= MEDIUMHEAT {
		thisTile.fireLevel = 2
	}
	if thisTile.heat >= MAXHEAT {
		thisTile.fireLevel = 3
	}

	if thisTile.neighborNorth != nil && !thisTile.neighborNorth.wall && thisTile.fireLevel != 0 {
		(thisTile.neighborNorth.heat) += thisTile.fireLevel
	}
	if thisTile.neighborEast != nil && !thisTile.neighborEast.wall && thisTile.fireLevel != 0 {
		(thisTile.neighborEast.heat) += thisTile.fireLevel
	}
	if thisTile.neighborWest != nil && !thisTile.neighborWest.wall && thisTile.fireLevel != 0 {
		(thisTile.neighborWest.heat) += thisTile.fireLevel
	}
	if thisTile.neighborSouth != nil && !thisTile.neighborSouth.wall && thisTile.fireLevel != 0 {
		(thisTile.neighborSouth.heat) += thisTile.fireLevel
	}
}

func assignNeighbor(thisTile *tile, x int, y int, maxX int, maxY int, tileMap [][]tile) {
	if x > 0 {
		thisTile.neighborNorth = &tileMap[x-1][y]
	}

	if y > 0 {
		thisTile.neighborWest = &tileMap[x][y-1]
	}

	if x < maxX-1 {
		thisTile.neighborSouth = &tileMap[x+1][y]
	}

	if y < maxY-1 {
		thisTile.neighborEast = &tileMap[x][y+1]
	}

	if x > 0 && y > 0 {
		thisTile.neighborNW = &tileMap[x-1][y-1]
	}
	if x > 0 && y < maxY-1 {
		thisTile.neighborNE = &tileMap[x-1][y+1]
	}
	if x < maxX-1 && y < maxY-1 {
		thisTile.neighborSE = &tileMap[x+1][y+1]
	}
	if x < maxX-1 && y > 0 {
		thisTile.neighborSW = &tileMap[x+1][y-1]
	}
}

func makeNewTile(thisPoint int, x int, y int) tile {

	//makes a basic floor tile with no nothin on it
	//and also no neighbors
	newTile := tile{x, y, 0, 0, false, false, nil, 0, false, nil, nil, nil, nil, nil, nil, nil, nil}

	if thisPoint == 0 {
		//make normal floor
		//helt normalt flour

		//append to tilemap
	} else if thisPoint == 1 {
		//wall
		newTile.wall = true
	} else if thisPoint == 2 {
		//door
		newTile.door = true
	} else if thisPoint == 3 {
		//out of bounds
		newTile.outOfBounds = true
	}

	return newTile
}

func TileConvert(inMap [][]int) [][]tile {
	mapXSize := len(inMap)
	mapYSize := len(inMap[0])

	//Initiates a slice of tile slices (2D tile slice)
	tileMap := make([][]tile, mapXSize)

	for x := 0; x < mapXSize; x++ {
		//initiates slice of tiles
		tileMap[x] = make([]tile, mapYSize)

		for y := 0; y < mapYSize; y++ {
			//constructs a new tile
			newTile := makeNewTile(inMap[x][y], x, y)

			//inserts tile into 2d slice
			tileMap[x][y] = newTile

		}
	}

	//Assigns 4 neighbors to each tile
	for x := 0; x < mapXSize; x++ {
		for y := 0; y < mapYSize; y++ {
			assignNeighbor(&(tileMap[x][y]), x, y, mapXSize, mapYSize, tileMap)
		}
	}

	return tileMap

}

func GetTile(inMap [][]tile, x int, y int) *tile {
	for i := range inMap {
		for j := range inMap[i] {
			if inMap[i][j].xCoord == x && inMap[i][j].yCoord == y {
				return &inMap[i][j]
			}
		}
	}
	return nil
}

func PeopleInit(inMap [][]tile, peopleList [][]int) []*Person {
	size := len(peopleList)
	peopleArray := make([]*Person, size)
	for i, person := range peopleList {

		tile := GetTile(inMap, person[1], person[0])
		peopleArray[i] = makePerson(tile)
	}
	return peopleArray
}

func StatsInit(size int) [][]int {
	statsList := make([][]int, size)
	for i := range statsList {
		statsList[i] = make([]int, 2) //change 2 ->3
	}
	return statsList
}

func StatsStart(statsList [][]int, peopleArray []*Person) [][]int{
	for i := range statsList {
		//peopleArray[i].GetStats(statsList[i])
		statsList[i] = peopleArray[i].GetStats()
	}
	return statsList
}


func StartStats(peopleList [][]int) [][]int {
	lst := [][]int{}
	for _, pers := range peopleList {
		templst := []int{}
		templst[0] = pers[0]
		templst[1] = pers[1]
		lst = append(lst, templst)
	}
	return lst
}


func Run(m *[][]tile, ppl []*Person, statsList *[][]int) {
	//sList := []Stats{}

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	wg.Add(len(ppl))
	//	*statsList = [][]int{}
	for _, pers := range ppl {

		go func(p *Person) { //, ind int){

			//	ind := i
			//	p := pers
			defer wg.Done()

			//	sList :=  []int{}// append(sList, p.getStats())
			p.MovePerson(m)
			//sList := &[]int{}
			//p.GetStats(statsList[i]) //(statsList[ind])
			
			mutex.Lock()
			//statsList = append(*statsList, *sList)
			mutex.Unlock()
			//	fmt.Println(len(statsList))

		}(pers) //, i)
	}

	step++
	wg.Wait()
	*statsList = StatsStart(*statsList, ppl)
	FireSpread(*m)
	//	}

	/*	// go run ruitnes for concurrency
		for _, person := range peopleArray {
			person.MovePerson(&inMap)
		}*/
	//	return sList
}

/**
func RunGo(inMap *[][]tile, peopleArray []*Person) []*tile { // OBS: not working
	movement := make([]*tile, len(peopleArray))
	var wg sync.WaitGroup

	wg.Add(len(peopleArray))
	for i, person := range peopleArray {
		go func(currentPerson *Person, ind int) {
			defer wg.Done()
			if !currentPerson.DiagonalStep() {
				fmt.Println("not diagonal move!")
				currentPerson.MovePerson(inMap)
			}
			if currentPerson.IsWaiting() {
				movement[ind] = nil
			} else {
				movement[ind] = currentPerson.path[len(currentPerson.path)-1]
			}
		}(person, i)
	}
	return movement
=======
	FireSpread(*m) // MOVED!

//	if math.Mod(float64(step), 40) == 0 {InitPlans(m)}
	// TODO: takes up tiiime!!1 fixy-changy

	//if math.Mod(float64(step), 5) == 0 {InitPlans(m)}
>>>>>>> jump
}
*/
func CheckFinish(peopleArray []*Person) bool {
	for i := 0; i < len(peopleArray); i++ {
		if peopleArray[i].safe == false && peopleArray[i].alive == true {
			return false
		}
	}
	return true
}

/*
func mainMap() {
	//	mainPath()
	//	MainPeople()
	//      testRedirect()
	//	testDiagPpl()
	//	testDiag()

	testDiagonally()
	//	testMovePeople()
	//	Whut()
	//	testJP()
	//	GLoop()

	//	testDiagonally()
	testMovePeople()
=======
 func CheckFinish (peopleArray []*Person) bool {
 	for i := 0; i < len(peopleArray); i++ {
 		if (peopleArray[i].safe == false && peopleArray[i].alive == true) {
 			return false
 		}
 	}
 	return true
 }

func mainMap() {
	if math.Mod(float64(step), 20) == 0 {}
//	mainPath()
//	MainPeople()
//      testRedirect()
//	testDiagPpl()
//	testDiag()
	debugging()
//	testSame()
//	testDiagonally()
//	testMovePeople()
//	Whut()
//	testJP()
	//	GLoop()

//	testDiagonally()
//	testMovePeople()
>>>>>>> jump

}

func testRedirect() {
	matrix := [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{1, 0, 1, 1, 1, 0, 1},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 2, 0, 0, 0, 0}}

	list := [][]int{
		{1, 1},
		{1, 2},
		{0, 2},
		{0, 0},
		{2, 1},
		{3, 1},
		{0, 3},
		{0, 1},
		{1, 0},
		{1, 3}}

	tryThis(matrix, list, -1, -1)
}

func testDiag() {
	matrix := [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 2}}

	list := [][]int{
		{0, 0},
		{0, 6}}
	tryThis(matrix, list, -1, -1)
}

func testDiagPpl() {
	matrix := [][]int{
		{0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 2}}

	list := [][]int{
		{0, 0},
		{0, 1},
		{1, 0},
		{0, 6},
		{2, 4}}

	tryThis(matrix, list, -1, -1)
}

func testDiagonally() {
	matrix := [][]int{
		{0, 1, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0, 0, 0},
		{0, 1, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 2}}

	list := [][]int{
		{0, 0},
		{0, 6},
		{2, 4}}

	tryThis(matrix, list, 1, 2)
}

func testMovePeople() {
	matrix := [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0},
		{1, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 2}}

	list := [][]int{
		{0, 0}}
	//	{0, 6},
	//	{2, 4}}

	tryThis(matrix, list, -1, -1)
	// Note: it takes 1 timeunit to take a step from the door and away
}

func testJP() {
	matrix := [][]int{
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{2, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 1}}

	list := [][]int{
<<<<<<< HEAD
		{0, 0}}
=======
		{0, 2}}
>>>>>>> jump
	//	{0, 6},
	//	{2, 4}}

	testmap := TileConvert(matrix)
	pplArray := PeopleInit(testmap, list)

	MovePeople(&testmap, pplArray)

	for i, p := range pplArray {
		fmt.Println("Person", i, "time:  ", p.time, "\n         health:", p.hp)
	}

}
*/
func tryThis(matrix [][]int, ppl [][]int, x, y int) {
	testmap := TileConvert(matrix)
	pplArray := PeopleInit(testmap, ppl)
	InitPlans(&testmap)

	if x >= 0 && y >= 0 {
		SetFire(&testmap[x][y])
	}
	MovePeople(&testmap, pplArray)

	for i, p := range pplArray {
		fmt.Println("Person", i, "time:  ", p.time, "\n         health:", p.hp)
	}
}


func testSame() {
	matrix := [][]int {
		{1,1,1,1,1,1,1,1,1,1,1,1,1},
		{1,0,0,0,0,0,1,0,0,0,0,0,1},
		{1,0,0,0,0,0,1,0,0,1,0,0,1},
		{1,0,0,0,0,0,1,1,1,1,0,0,1},
		{1,1,1,0,0,0,1,0,0,0,0,0,1},
		{2,0,0,0,0,0,0,0,0,0,0,0,1},
		{1,1,1,1,1,1,1,1,1,1,1,1,1}}

	list := [][]int{{3,2}, {5,5}}

	tryThis(matrix, list, -1, -1)
}



func debugging() {

	b, err3 := ioutil.ReadFile("../src/mapfile.txt")
	if err3 != nil{
		panic(err3)
	}

	var m = [][]int{}
	err := json.Unmarshal(b, &m)
	if err != nil{
		panic(err)
	}
	//m[8][1] = 2
	//m[13][0] = 2

	//	testmap := TileConvert(m)


	c, err4 := ioutil.ReadFile("../src/playerfile.txt")
	if err4 != nil{
		panic(err4)
	}

	var mm = [][]int{}
	err5 := json.Unmarshal(c, &mm)
	if err5 != nil{
		panic(err5)
	}

//	ppl := PeopleInit(testmap, mm)

//	list := [][]int{{89,33}}//{104, 28}, {105, 29}}  // lr tvärtom?
	tryThis(m, mm, 20, 20)
}



//send heat
func FireStats(m *[][]tile) [][]int{
	//func FireStats(start []*tile, dir Direction) [][]int{
	fire := [][]int{}

	for i, list := range *m {
		for j, _ := range list {
			//tl := GetTile(*m, i, j)
			tl := &(*m)[i][j]
			if tl.heat  > 0 {fire = append(fire, []int{tl.yCoord, tl.xCoord, tl.heat})}
		}
	}
	return fire
}
