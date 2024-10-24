package chess

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	EXE_FILE_PATH_ENV     = "PATH_TO_EXECUTABLE"
	SCRIPTS_FILE_PATH_ENV = "SCRIPTS_FILE_PATH"
	DEFAULT_PATH          = "./stockfish/stockfish-ubuntu-x86-64-avx2"
	DEFAULT_SCRIPTS_PATH  = "./scripts"
	MOVE_SCRIPT           = "move.sh"
	UPDATE_FEN_SCRIPT     = "update_fen.sh"
)

type TableState string

type Move struct {
	Move  string
	Table TableState
}

type Driver struct {
	exePath     string
	scriptsPath string
}

const (
	BASE_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

func New() *Driver {
	exepath := os.Getenv(EXE_FILE_PATH_ENV)
	if exepath == "" {
		exepath = DEFAULT_PATH
	}

	scriptsPath := os.Getenv(SCRIPTS_FILE_PATH_ENV)
	if scriptsPath == "" {
		scriptsPath = DEFAULT_SCRIPTS_PATH
	}

	return &Driver{exePath: exepath, scriptsPath: scriptsPath}
}

func (d *Driver) Move(skillLevel uint16, state TableState) (*Move, error) {

	if !state.IsValid() {
		return nil, errors.New("stockfish: invalid fen state")
	}

	buf := bytes.NewBuffer([]byte{})

	cmd := exec.Command(getScriptFile(d.scriptsPath, MOVE_SCRIPT), d.exePath, fmt.Sprint(skillLevel), string(state))
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, errors.New("stockfish: error occured when running stockfish command " + err.Error())
	}

	output := buf.String()
	moveTxt := parseOutput(output)
	if moveTxt == "" {
		return nil, errors.New("stockfish: couldn't parse stockfish output - " + output)
	}

	cmd = exec.Command(getScriptFile(d.scriptsPath, UPDATE_FEN_SCRIPT), d.exePath, string(state), moveTxt)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, errors.New("stockfish: error occured when running stockfish command " + err.Error())
	}

	output = buf.String()
	fenstr := parseFEN(output)

	return &Move{Move: moveTxt, Table: fenstr}, nil
}

func parseOutput(output string) string {
	output = strings.Replace(output, "\n", " ", -1)
	words := strings.Split(output, " ")
	next := false
	for _, word := range words {
		if next {
			return word
		}
		if word == "bestmove" {
			next = true
		}
	}
	return ""
}

func parseFEN(fenStr string) TableState {
	fs := strings.Index(fenStr, "Fen: ")
	fe := strings.Index(fenStr, "Key: ")
	return TableState(fenStr[fs+5 : fe-1])
}

func getScriptFile(folder string, name string) string {
	return folder + "/" + name
}

func (s TableState) IsValid() bool {
	// Define the regex pattern
	pattern := `^([rnbqkpRNBQKP1-8]+\/){7}[rnbqkpRNBQKP1-8]+ [wb] [KQkq-]{1,4} ([a-h][36]|-) \d+ \d+$`

	// Compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		// Handle error
		return false
	}

	// Check if the fen matches the pattern
	return re.MatchString(string(s))
}

func (d *Driver) EvaluateWinProbability(skillLevel uint16, state TableState) (string, error) {
	if !state.IsValid() {
		return "", errors.New("stockfish: invalid fen state")
	}

	buf := bytes.NewBuffer([]byte{})

	cmd := exec.Command(getScriptFile(d.scriptsPath, "evaluate_win_probability.sh"), d.exePath, fmt.Sprint(skillLevel), string(state))
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", errors.New("stockfish: error occurred when running stockfish command " + err.Error())
	}

	output := buf.String()
	// Предполагаем, что скрипт возвращает вероятность в виде строки
	return output, nil
}
