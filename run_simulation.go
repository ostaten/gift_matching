package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ostaten/gift_matching/data"
	"github.com/ostaten/gift_matching/util"
)

var AllData data.Data

/*
 $ Gets all of the JSON files and puts them into the global data structure
*/
func getJSONFiles() error {
	all, groupNames, groups, err := retrieveGroups()

	AllData.AllPeople = all
	AllData.GroupNames = groupNames
	AllData.Groups = groups

	if err != nil {
		fmt.Println(err)
		return err
	}

	numAssignments, pastAssignments, rawHistory, err := retrieveHistory(all)

	AllData.NumAssignments = numAssignments
	AllData.PastAssignments = pastAssignments
	AllData.RawHistoryData = rawHistory

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/*
 $ Retrieves all information from the groups.json file and returns each of the three sections:
 $ - all: all members of the gift exchange
 $ - GroupNames: the last names and categorization of the gift exchange
 $ - SubGroups: the map of each surname containing an array of the members of said subgroup inside
*/
func retrieveGroups() ([]string, []string, map[string][]string, error) {
	jsonFile, err := os.Open("./resources/groups.json");

	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	
	var subGroups SubGroups

	json.Unmarshal([]byte(byteValue), &subGroups)

	return subGroups.All, subGroups.GroupNames, subGroups.SubGroups, nil
}

/*
 $ Retrieves all information from the history.json file
 */
func retrieveHistory(people []string) (int, map[string][]string, []map[string]string, error) {
	jsonFile, err := os.Open("./resources/history.json");

	if err != nil {
		fmt.Println(err)
		return 0, nil, nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	
	var result History

	json.Unmarshal([]byte(byteValue), &result)

	var iterations int

	var eachPersonsHistory map[string][]string = make(map[string][]string)

	for _, assignment := range result.Data {
		for _, person := range people {
			eachPersonsHistory[person] = append(eachPersonsHistory[person], assignment[person])
		}
    if err != nil {
			fmt.Println(err)
      return 0, nil, nil, err
    }
		iterations++
	}

	return iterations, eachPersonsHistory, result.Data, nil
}

/*
 $ Asks the user for how thorough of a search they desire -> this correlates to how many attempts we will do before "giving up".
 $ Stores valid result in global structure
 $ Once we receive that info, we go to the next user query.
 */
func getTriesAllowed(repeat bool) (bool, error) {
	fmt.Print("To start, please select how determined you are to ensure a gift giving match is found.  Enter 'QUICK', 'DEFAULT', or, 'THOROUGH': ")
	var triesAllowedMode string
	if repeat {
		fmt.Scanf(" %c", &triesAllowedMode)
	} 
	fmt.Scanf("%s\n", &triesAllowedMode)
	triesAllowedMode = strings.ToLower(triesAllowedMode)

	if (triesAllowedMode == "back") {
		return getTriesAllowed(true)
	}

	if (triesAllowedMode == "exit") {
		return true, nil
	}

	for (triesAllowedMode != "quick" && triesAllowedMode != "thorough" && triesAllowedMode != "default" && triesAllowedMode != data.BSpaceHold && triesAllowedMode != "") {
		triesAllowedMode = data.BSpaceHold
		fmt.Print("Invalid input! Please enter 'QUICK', 'DEFAULT', or, 'THOROUGH': ")
		fmt.Scanf("%s\n", &triesAllowedMode)
		triesAllowedMode = strings.ToLower(triesAllowedMode)
		if (triesAllowedMode == "back") {
			return getTriesAllowed(true)
		}
	
		if (triesAllowedMode == "exit") {
			return true, nil
		}
	}

	if (triesAllowedMode == "quick") {
		AllData.TriesAllowed = 100
	}

	if (triesAllowedMode == "thorough") {
		AllData.TriesAllowed = 100000000
	}

	// var nextMode string = triesAllowedMode

	if (triesAllowedMode == "default" || triesAllowedMode == data.BSpaceHold || triesAllowedMode == "") {
		triesAllowedMode = "default"
		AllData.TriesAllowed = 1000
	}

	fmt.Println("Got it, doing", strings.ToUpper(triesAllowedMode), "matching . . . ")

	isExit, err := getDepth(false)

	if isExit {
		return true, nil
	}

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return false, nil
}

/*
 $ After asking for tries, this function queries the amount of iterations back that we want to look (1, 2, etc). 
 $ Stores valid result into global structure 
 $ After doing this, we then go on to actually execute the matching
 */ 
func getDepth(needBuffer bool) (bool, error) {
	fmt.Println("Next, please select how many iterations of gift giving back you want to go from 1 ->", strconv.Itoa(AllData.NumAssignments), "or just put 'DEFAULT'.")
	fmt.Println("This number is the amount of iterations in the past you wish to avoid having anyone give to each other.")
	fmt.Print("The larger the number, the less recent each person will have given to someone else, but the less likely there will be a match: ")
	var input string
	if needBuffer {
		fmt.Scanf(" %c", &input)
	}
	fmt.Scanf("%s", &input)
	input = strings.ToLower(input)

	if (input == "back") {
		return getTriesAllowed(true)
	}

	if (input == "exit") {
		return true, nil
	}

	if (input == "default" || input == "") {
		input = "default"
		AllData.Depth = int(math.Min(float64(AllData.NumAssignments), 2))
	}

	numTries, _ := strconv.Atoi(input)

	for (numTries < 1 || numTries > AllData.NumAssignments) {
		if input == "default" || input == "" {
			input = "default"
			break
		}

		if (input == "back") {
			return getTriesAllowed(true)
		}
	
		if (input == "exit") {
			return true, nil
		}

		if (numTries == 0) {
			fmt.Print("Invalid input! Please enter an integer between 0 and ", AllData.NumAssignments + 1, " or just 'DEFAULT': ")
			fmt.Scanf(" %c", &input)
			fmt.Scanf("%s", &input)
			input = strings.ToLower(input)
			numTries, _ = strconv.Atoi(input)
			continue
		}

		if (numTries < 1) {
			fmt.Print("Invalid input! Please enter an integer greater than 0: ")
			fmt.Scanf(" %c", &input)
			fmt.Scanf("%s", &input)
			input = strings.ToLower(input)
			numTries, _ = strconv.Atoi(input)
			continue
		}

		if (numTries > AllData.NumAssignments) {
			fmt.Print("Invalid input! Please enter an integer less than ", AllData.NumAssignments + 1, ": ")
			fmt.Scanf(" %c", &input)
			fmt.Scanf("%s", &input)
			input = strings.ToLower(input)
			numTries, _ = strconv.Atoi(input)
			continue
		}
	}

	if (input == "default" || input == "") {
		AllData.Depth = int(math.Min(float64(AllData.NumAssignments), 2))
	} else {
		AllData.Depth = numTries
	}

	fmt.Println("Got it, reaching a depth of", AllData.Depth, ". . . ")

	return false, nil
}

/*
 $ Initiates the user queries
 */
func getCommandLineInitiators(startIndent bool) (bool, error) {
	fmt.Println("Hello there! We've successfully loaded up our previous data and are ready to initiate matching.")
	fmt.Println("-----------------------------------------------------------------------------------------------")
	fmt.Println("If at any point you want to abandon the match, type 'EXIT'.  If you would like to go back a step, type 'BACK'.")

	isExit, err := getTriesAllowed(startIndent)
	
	if isExit {
		return true, nil
	}

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return false, nil
}


/*
 $ This is the core of the program. 
 $ For each person, starts with everyone, then allows their options to be a) not themselves, b)
 $ not their family and c) not anyone in the past iterationsDeep iterations they've given to. 
 $ If we can't find someone, returns an error.
 */
func createPossibleMatches(iterationsDeep int) ([]data.Possibility, error) {
	var res []data.Possibility
	for _, person := range AllData.AllPeople {
		var newEntry data.Possibility
		newEntry.Name = person
		var newPossibilities []string
		var surname string = ""
		for _, lastName := range AllData.GroupNames {
			if util.Contains[string](AllData.Groups[lastName], person) {
				surname = lastName
				break
			}
		}
		if surname == "" {
			fmt.Println("Ugh, couldn't find family surname")
			return nil, errors.New("ugh, couldn't find family surname")
		}

		var unusablePeople []string = append(AllData.PastAssignments[person][0 : iterationsDeep], AllData.Groups[surname]...)

		for _, possibility := range AllData.AllPeople {
			if (possibility != person && !util.Contains[string](unusablePeople, possibility)) {
				newPossibilities = append(newPossibilities, possibility)
			}
		}
		newEntry.Possibilities = newPossibilities
		var added bool = false
		for i := 0; i < len(res); i++ {
			if (len(res[i].Possibilities) > len(newPossibilities)) {
				res = append(res, newEntry)
				copy(res[i + 1:], res[i:]) 
				res[i] = newEntry     
				added = true
				break  
			}
		}
		if !added {
			res = append(res, newEntry)
		}
	}
	return res, nil
}

/*
 $ Helper function for runSimulation() where we remove a certain person from everyone's pool of people (they've now been given to)
 */
func removePersonFromPool(allPools []data.Possibility, selectedPerson string) []data.Possibility {
	var revisedPos []data.Possibility
	for _, pool := range allPools {
		var newPoss data.Possibility
		newPoss.Name = pool.Name
		for _, person := range pool.Possibilities {
			if (person != selectedPerson) {
				newPoss.Possibilities = append(newPoss.Possibilities, person)
			}
		}
		revisedPos = append(revisedPos, newPoss)
	}
	return revisedPos
}

/*
 $ Helper function for runSimulation() that deep copies the possiblity data structure so we can repeat as many times as desired
*/
func copyPossibleMatches() ([]data.Possibility) {
	var trialPossibilities []data.Possibility

	for _, person := range AllData.PossibleMatches {
		var newEntry data.Possibility
		newEntry.Name = person.Name
		newEntry.Possibilities = append(newEntry.Possibilities, person.Possibilities...) 
		trialPossibilities = append(trialPossibilities, newEntry)
	}
	return trialPossibilities
}

/*
 $ Makes a set of assignments for gift giving for the upcoming iteration based on their previous gift giving and groupings
*/
func runSimulation() (map[string]string, bool, error) {
	var res map[string]string = make(map[string]string)
	var trialPossibilities []data.Possibility = copyPossibleMatches()
	for (len(trialPossibilities) > 0) {
		if (len(trialPossibilities[0].Possibilities) == 0) {
			return nil, false, errors.New("ran out of options")
		} 
		randomIndex := util.GetRandom(0, len(trialPossibilities[0].Possibilities))
		selectedPerson := trialPossibilities[0].Possibilities[randomIndex]
		res[trialPossibilities[0].Name] = selectedPerson
		trialPossibilities = append(trialPossibilities[:0], trialPossibilities[1:]...)
		trialPossibilities = removePersonFromPool(trialPossibilities, selectedPerson)
	}

	if len(res) < len(AllData.AllPeople) {
		return nil, false, errors.New("didn't have enough distribution")
	}
	return res, true, nil
} 

/*
 $ Pretty prints the result of a successful simulation run
*/
func printSimulation(assignments map[string]string, triesAttempted int, startWidth int) {
	fmt.Println("Had to try", triesAttempted, "time(s)")
	fmt.Println()
	for i := 0; i < startWidth; i++ {
		fmt.Print("*")
	}
	fmt.Println()
	for giver, whoTheyGiveTo := range assignments {
		fmt.Print("*")
		statement := fmt.Sprintf("%s will give to %s", giver, whoTheyGiveTo)
		for i := 0; i < (startWidth - 2 - len(statement)) / 2; i++ {
			fmt.Print(" ")
		}
		fmt.Print(statement)
		for i := 0; i < startWidth - 2 - (startWidth - 2 - len(statement)) / 2 - len(statement); i++ {
			fmt.Print(" ")
		}
		fmt.Print("*")
		fmt.Println()
	}
	for i := 0; i < startWidth; i++ {
		fmt.Print("*")
	}
	fmt.Println()
	fmt.Println()
}

/*
 $ Permanently records the current iteration's assignments.
 */
func saveAssignments(assignments map[string]string) (error) {
	AllData.NumAssignments++

	AllData.RawHistoryData = append([]map[string]string{assignments}, AllData.RawHistoryData...)
	var fullStructure History
	fullStructure.Data = AllData.RawHistoryData
	rawBoi, err := json.MarshalIndent(fullStructure, "", "  ")

	os.WriteFile("./resources/history.json", rawBoi, 0777)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Successfully saved matches in permanent data!")

	return nil
}

/*
 $ Asks the user what they'd like to do after the simulation has successfully been run. 
 $ They can either use it and permanently etch it into history.json, redo it using the parameters selected, or restart from the beginning.
 */
func confirmNextAction(assignments map[string]string) (bool, error) {
	fmt.Print("What would you like to do with your matching?  To use it and enter it into the database, enter 'KEEP'.  To do it again with the same parameters, enter 'REDO'.  To restart with different parameters, enter 'RESTART': ")
	var input string
	fmt.Scanf("%s\n", &input)
	input = strings.ToLower(input)

	if (input == "exit") {
		return true, nil
	}

	for (input != "" && input != "keep" && input != "redo" && input != "restart" && input != data.BSpaceHold) {
		input = data.BSpaceHold
		fmt.Print("Invalid input! Please enter 'KEEP', 'REDO', or 'RESTART': ")
		fmt.Scanf("%s\n", &input)
		input = strings.ToLower(input)

		if (input == "exit") {
			return true, nil
		}

	}

	if (input == "redo" || input == "" || input == data.BSpaceHold) {
		input = "redo"
		return doAllSimulations()
	}

	if (input == "keep") {
		err := saveAssignments(assignments)

		if err != nil {
			return false, err
		}

		return false, nil
	}

	if (input == "restart") {
		return startMainProcess(false)
	}


	return false, nil
}

/*
 $ Wrapper function to do all the simulations the amount of times the user specified, and is called by the "redo" option in confirmNextAction()
 */
func doAllSimulations() (bool, error) {
	var assignments map[string]string
	var didIt bool = false

	var err error

	var i = AllData.TriesAllowed
	for !didIt && i > 0 {
		assignments, didIt, err = runSimulation()
		i--
	}


	if i == 0 {
		fmt.Println("Could not complete assigment in", AllData.TriesAllowed, "tries")
		return false, nil
	}

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	printSimulation(assignments, AllData.TriesAllowed - i, 60)

	return confirmNextAction(assignments)
}

/*
 $ Wrapper function for getting the command line stuff, and then goes on to do everything else.  Used for the "restart" option in confirmNextAction()
 */
func startMainProcess(startIndent bool) (bool, error) {
	isExit, err := getCommandLineInitiators(startIndent)

	if isExit {
		return true, nil
	}

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	matches, err := createPossibleMatches(AllData.Depth)

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	AllData.PossibleMatches = matches

	isExit, err = doAllSimulations()

	if isExit {
		return true, nil
	}

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return false, nil
}


func main() {
	err := getJSONFiles()

	if err != nil {
		fmt.Println("Oops, we couldn't load previous data!")
		return
	}

	isExit, err := startMainProcess(false)


	if isExit {
		fmt.Println("Hope to see you again soon!")
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Hope to see you again soon!")
}


/*
 $ Defines structure of json found in the groups file
*/
type SubGroups struct {
	All        []string `json:"all"`
	GroupNames   []string `json:"groupNames"`
	SubGroups   map[string][]string   `json:"subGroups"`
}

/*
 $ Defines structure of json found in the history file
*/

type History struct {
	Data    []map[string]string `json:"data"`
}