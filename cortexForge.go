package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type UsernameData struct {
	Username string `json:"username"`
}

type UsernameResponse struct {
	Data UsernameData `json:"data"`
}

type PuzzleData struct {
	Puzzle string `json:"puzzle"`
	Username string `json:"username,omitempty"`
	Achievement string `json:"achievement,omitempty"`
	Completed bool `json:"completed"`
}

type PuzzleResponse struct {
	Data PuzzleData `json:"data"`
}

type AnswerRequest struct {
	Username string `json:"username"`
	Answer   string `json:"answer"`
}

const (
	apiToken        = "redacted"
	usernameAPI     = "https://pointlesscode.dev/redacted"
	questionAPI     = "https://pointlesscode.dev/redacted"
	submitAnswerAPI = "https://pointlesscode.dev/redacted"
)

var (
	envUsername = os.Getenv("username")
	username    string
	puzzle      string
	achievement string
	completed   = false
)

func main() {
	initializeUsername()
	getPuzzle(username)

	for {
		displayGameStatus()
		if completed == true {
			endGame()
		} else {
			fmt.Printf("Puzzle: %s\n\n", puzzle)
			answer := getUserInput("Your Answer: ")
			submitAnswer(username, answer)
		}
	}
}

func initializeUsername() {
	if envUsername != "" {
		username = envUsername
		newUsername, err := getPuzzle(username)
		if err != nil {
			username = newUsername
		}
	} else {
		username = getUsername()
	}
}

func displayGameStatus() {
	clearTerminal()
	dashCounter := max(len("New Achievement Unlocked: "+achievement), 44)

	fmt.Println(".less X CortexForge")
	fmt.Println(strings.Repeat("-", dashCounter) + "\n")
	fmt.Println("CTRL + C at any time to quit.\n")
	fmt.Printf("Your username: %s\n\n", username)

	if achievement != "" {
		fmt.Println(strings.Repeat("-", dashCounter))
		fmt.Println("New Achievement Unlocked: " + achievement)
		fmt.Println(strings.Repeat("-", dashCounter) + "\n")
	} else {
		fmt.Println(strings.Repeat("-", dashCounter) + "\n")
	}
}

func endGame() {
	fmt.Println("Well done sherlock!\n")
	fmt.Println("Go and check your achievements at:")
	fmt.Printf("https://pointlesscode.dev/cortex-forge/achievements/%s\n\n", username)
	waitForExit()
}

func waitForExit() {
	fmt.Println("Press Enter to exit...")
	getUserInput("")
	os.Exit(0)
}

func getUserInput(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return input
}

func clearTerminal() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println("Unsupported platform. Cannot clear terminal.")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getUsername() string {
	response := &UsernameResponse{}
	makeHTTPRequest("GET", usernameAPI, nil, response)
	return response.Data.Username
}

func makeHTTPRequest(method, url string, body []byte, response interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %v", err)
	}

	if err := json.Unmarshal(bodyBytes, response); err != nil {
		return fmt.Errorf("Error unmarshalling response: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("resource not found (404)")
	}

	return nil
}

func getPuzzle(username string) (string, error) {
	reqBody, _ := json.Marshal(map[string]string{"username": username})

	response := &PuzzleResponse{}
	err := makeHTTPRequest("POST", questionAPI, reqBody, response)

	if err != nil {
		if response.Data.Username != "" {
			username = response.Data.Username
			return username, fmt.Errorf("404")
		}
		return "", fmt.Errorf("Error making request: %v", err)
	}

	puzzle = response.Data.Puzzle
	completed = response.Data.Completed
	return "completed", nil
}

func submitAnswer(username, answer string) error {
	reqBody, err := json.Marshal(AnswerRequest{Username: username, Answer: answer})
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	response := &PuzzleResponse{}
	err = makeHTTPRequest("POST", submitAnswerAPI, reqBody, response)
	if err != nil {
		return fmt.Errorf("error submitting answer: %v", err)
	}

	puzzle = response.Data.Puzzle
	completed = response.Data.Completed
	if response.Data.Achievement != "" {
		achievement = response.Data.Achievement
	}

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
