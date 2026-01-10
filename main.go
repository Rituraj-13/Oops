package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

func main() {

	// var ApiKey string

	var style = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Width(44)

	// fmt.Println(style.Render("Hello, kitty"))
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(style.Render("Please enter your Gemini API key: "))
	line, _ := reader.ReadString('\n')

	fmt.Printf("Your API Key is : %s", line)
	
}
