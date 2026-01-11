package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type config struct {
	name   string `json:"name"`
	APIKey string `json:apiKey`
}

func main() {

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Width(44)

}

// This function will be used to show the message and take the inptut from the terminal

func PromptLine(reader *bufio.Reader, style lipgloss.Style, label string)(string, error){
	fmt.Println(style.Render(label))
	s, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF){
		return "", nil
	}

	s = strings.TrimSpace(s)
	if s == "" {
		return "", errors.New("Value cant be empty")
	}
	return s, nil
}
``
func config