package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
	"strconv"
)

type Rules struct {
	Name     string `json:"name"`
	Birth    []int  `json:"birth"`
	Survival []int  `json:"survival"`
}

type AutomataChallenge struct {
	Cells       [][]int `json:"cells"`
	Rules       Rules   `json:"rules"`
	Generations int     `json:"generations"`
}

type AutomataChallengeResponse struct {
	Challenge     AutomataChallenge `json:"challenge"`
	ChallengePath string            `json:"challengePath"`
}

type AutomataSolutionResponse struct {
	res string `json:"result"`
	msg string `json:"message"`
}

func main() {
	if(true){
		correctCount := 0
		loopCount := 10
		for i := 0; i < loopCount; i++ {		
			var response = getChallenge()

			printChallenge(response.Challenge)
			printBoard("Before", response.Challenge.Cells)

			var cells = getSolution(response.Challenge)

			printBoard("After", cells)

			if(sendSolution(formatSolution(cells), response.ChallengePath)) {
				print("The solution was correct!\n")
				correctCount++
			} else {
				print("The soluton was incorrect :(\n")
			}
		}
		print("Passed " + strconv.Itoa(correctCount) + " out of " + strconv.Itoa(loopCount) + ".\n")

	} else {	
		// Used for debugging
		var challenge AutomataChallenge 
		var rules Rules

		rules.Name = string("conway")
		rules.Birth = []int{3}
		rules.Survival = []int{2,3}

		// challenge.Cells = [][]int{
		// 	{0,0,0,0,0,0,0,0},
		// 	{0,0,0,0,0,0,0,0},
		// 	{0,0,1,1,1,0,0,0},
		// 	{0,0,1,0,0,0,0,0},
		// 	{0,0,0,0,0,0,0,0}}

		// challenge.Cells = [][]int{
		// 	{1,0,0,0,0,0},
		// 	{1,0,0,0,0,0},
		// 	{0,0,0,0,0,0},
		// 	{0,0,0,0,0,0},
		// 	{1,0,0,0,0,0}}

		challenge.Cells = [][]int{
			{0,0,0,1,0,0},
			{0,1,1,0,0,0},
			{0,0,1,0,0,0},
			{0,0,0,0,0,0},
			{0,0,0,0,0,0}}

		// challenge.Cells = [][]int{
		// 	{0,0,0,0,0,0,0,0,0},
		// 	{0,0,0,0,1,1,0,0,0},
		// 	{0,0,0,1,1,0,0,0,0},
		// 	{0,0,0,0,1,0,0,0,0},
		// 	{0,0,0,0,0,0,0,0,0},
		// 	{0,0,0,0,0,0,0,0,0}}

		challenge.Rules = rules
		challenge.Generations = 10

		printChallenge(challenge)
		printBoard("Before", challenge.Cells)

		var cells = getSolution(challenge)

		printBoard("After", cells)

		print("\n" + formatSolution(cells) + "\n")
	}
}

func formatSolution(cells [][]int ) string {
	var str string = "["
	for y := range cells {
		str += " [ "
		for x := range cells[y] {
			if cells[y][x] == 1 {
				str += "1"
			} else {
				str += "0"
			}
			if x < len(cells[y]) - 1{
				str += " ,"
			} else {
				str += " "
			}
		}
		if y < len(cells) - 1{
			str += "],"
		} else {
			str += "]"
		}
	}
	str += "]"
	return str
}

func is_in(x int, arr []int) bool {
	for i := range arr {
		if(arr[i] == x){
			return true
		}
	}
	return false
}

func to_index(val int, max int) int {
	if(val == max) {
		val = 0
	}
	if(val < 0) {
		val = max-1
	}
	return val
}

func get_cell(cells [][]int, x int, y int) int {
	max_x := len(cells[0])
	max_y := len(cells)

	if(x >= max_x || x < 0 || y >= max_y || y < 0) {
		return 0
	}
	return cells[y][x]
}

func getSolution(challenge AutomataChallenge) [][]int {
	var rules Rules = challenge.Rules
	var generations int = challenge.Generations
	var count int = 0
	var i int
	cells := challenge.Cells
	numRows := len(cells)
	numCols := len(cells[0])
	a_cells := make([][]int, numRows)
	b_cells := make([][]int, numRows)
	for i := 0; i < numRows; i++ {
		a_cells[i] = make([]int, numCols)
		b_cells[i] = make([]int, numCols)
	}
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			a_cells[i][j] = cells[i][j]
			b_cells[i][j] = cells[i][j]
		}
	}

	for i = 0; i < generations; i++ {
		for y := range b_cells {
			for x := range b_cells[y] {
				// Get the number of living neighbours, note that we are using a wrapped cell space.
				count = get_cell(b_cells, x+1, y+1)  	+//top right
						get_cell(b_cells, x, y+1)  		+//top cent
						get_cell(b_cells, x-1, y+1) 	+//top left
						get_cell(b_cells, x+1, y) 		+//cent right
					  	get_cell(b_cells, x-1, y) 		+//cent left
					  	get_cell(b_cells, x+1, y-1) 	+//bot right
					 	get_cell(b_cells, x, y-1) 		+//bot cent
					  	get_cell(b_cells, x-1, y-1) 	 //bot left
				// Use this for wrapped cell space: 
				// count = b_cells[to_index(y+1, len(b_cells))][to_index(x+1, len(b_cells[y]))]  	+//top right
				// 		b_cells[to_index(y+1, len(b_cells))][x] 								+//top cent
				// 		b_cells[to_index(y+1, len(b_cells))][to_index(x-1, len(b_cells[y]))] 	+//top left
				// 		b_cells[(y)][to_index(x+1, len(b_cells[y]))]							+//cent right
				// 	  	b_cells[(y)][to_index(x-1, len(b_cells[y]))]           					+//cent left
				// 	  	b_cells[to_index(y-1, len(b_cells))][to_index(x+1, len(b_cells[y]))] 	+//bot right
				// 	  	b_cells[to_index(y-1, len(b_cells))][x] 								+//bot cent
				// 	  	b_cells[to_index(y-1, len(b_cells))][to_index(x-1, len(b_cells[y]))]	 //bot left
				
				// To debug the assignments:
				// print("b_cells: \n")
				// for k := 0; k < numRows; k++ {
				// 	for j := 0; j < numCols; j++ {
				// 		print(strconv.Itoa(b_cells[k][j]))
				// 	}
				// 	print("\n")
				// }
				// print("(" + strconv.Itoa(y) + ", " + strconv.Itoa(x) + "): " + strconv.Itoa(count) + "\n")
				// print("b_cells[" + strconv.Itoa(to_index(y+1, len(b_cells))) + "][" + strconv.Itoa(to_index(x+1, len(b_cells[y]))) + "] = " + strconv.Itoa(b_cells[to_index(y+1, len(b_cells))][to_index(x+1, len(b_cells[y]))]) + "\n")
				// print("b_cells[" + strconv.Itoa(to_index(y+1, len(b_cells))) + "][" + strconv.Itoa(x) + "] = " + strconv.Itoa(b_cells[to_index(y+1, len(b_cells))][x]) + "\n")
				// print("b_cells[" + strconv.Itoa(to_index(y+1, len(b_cells))) + "][" + strconv.Itoa(to_index(x-1, len(b_cells[y]))) + "] = " + strconv.Itoa(b_cells[to_index(y+1, len(b_cells))][to_index(x-1, len(b_cells[y]))]) + "\n")
				// print("b_cells[" + strconv.Itoa(y) + "][" + strconv.Itoa(to_index(x+1, len(b_cells[y]))) + "] = " + strconv.Itoa(b_cells[y][to_index(x+1, len(b_cells[y]))]) + "\n")
				// print("b_cells[" + strconv.Itoa(y) + "][" + strconv.Itoa(to_index(x-1, len(b_cells[y]))) + "] = " + strconv.Itoa(b_cells[y][to_index(x-1, len(b_cells[y]))]) + "\n")
				// print("b_cells[" + strconv.Itoa(to_index(y-1, len(b_cells))) + "][" + strconv.Itoa(to_index(x+1, len(b_cells[y]))) + "] = " + strconv.Itoa(b_cells[to_index(y-1, len(b_cells))][to_index(x+1, len(b_cells[y]))]) + "\n")
				// print("b_cells[" + strconv.Itoa(to_index(y-1, len(b_cells))) + "][" + strconv.Itoa(x) + "] = " + strconv.Itoa(b_cells[to_index(y-1, len(b_cells))][x]) + "\n")
				// print("b_cells[" + strconv.Itoa(to_index(y-1, len(b_cells))) + "][" + strconv.Itoa(to_index(x-1, len(b_cells[y]))) + "] = " + strconv.Itoa(b_cells[to_index(y-1, len(b_cells))][to_index(x-1, len(b_cells[y]))]) + "\n")
				if b_cells[y][x] == 1 {	
					// The cell was alive, check if it survives
					if(is_in(count, rules.Survival)){
						// The neighbour count is in the survival list, the cell lives!
						a_cells[y][x] = 1
					} else {
						// The cell dies, not enough friends!
						a_cells[y][x] = 0
					}
				} else {
					// The cell was dead, check for births
					if(is_in(count, rules.Birth)){
						// The neighbour count is in the birth list, new baby cell!
						a_cells[y][x] = 1
					} else {
						// No parents, no new cells!
						a_cells[y][x] = 0
					}
				}
			}
		}
		// The itteration is done, update the before cell space before itterating again
		for j := 0; j < numRows; j++ {
			for k := 0; k < numCols; k++ {
				b_cells[j][k] = a_cells[j][k]
				a_cells[j][k] = 0
			}
		}
		//printBoard("Itteration " + strconv.Itoa(i+1), b_cells)
	}
	copy(cells, b_cells)

	return cells
}

func getChallenge() AutomataChallengeResponse {
	domain := "https://api.noopschallenge.com"

	// get and parse challenge from the api
	res, err := http.Get(domain + "/automatabot/challenges/new")
	if err != nil {
		panic(err.Error())
	}

	return parseApiResponse(res)
}

func sendSolution(sol string, resp string) bool {
	domain := "https://api.noopschallenge.com"
	print(resp + "\n")
	// get and parse challenge from the api
	res, err := http.Post(domain + resp, "application/json", bytes.NewBuffer([]byte(sol)))
	 // http.NewRequest("POST", , sol)
	if err != nil {
		panic(err.Error())
	}

	//var response AutomataSolutionResponse

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	print(string(body) + "\n")
	return string(body) ==  "{\"result\":\"correct\"}"
}


func parseApiResponse(res *http.Response) AutomataChallengeResponse {
	var response AutomataChallengeResponse
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err.Error())
	}

	return response
}

// print info about the ruleset we are using
func printChallenge(challenge AutomataChallenge) {
	fmt.Println()
	fmt.Println("Ruleset:", challenge.Rules.Name)
	fmt.Println("Birth:", challenge.Rules.Birth)
	fmt.Println("Survival:", challenge.Rules.Survival)
	fmt.Println("Generations:", challenge.Generations)
}

// print a board for human consumption, with borders
func printBoard(title string, cells [][]int) {
	fmt.Println()
	fmt.Println(title)
	for x := 0; x < len(cells[0])+2; x++ {
		fmt.Print("=")
	}
	fmt.Println()
	for y := range cells {
		fmt.Print("|")
		for x := range cells[y] {
			if cells[y][x] == 1 {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("|")
		fmt.Println()
	}
	for x := 0; x < len(cells[0])+2; x++ {
		fmt.Print("=")
	}
	fmt.Println()
}
