package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	Name   string `json:"name"`
	APIKey string `json:"apiKey"`
}

func main() {

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Width(44)

}

// Ask the user for one line of input in the terminal and return it (trimmed).
func PromptLine(reader *bufio.Reader, style lipgloss.Style, label string) (string, error) {
	fmt.Println(style.Render(label))
	s, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	s = strings.TrimSpace(s)
	if s == "" {
		return "", errors.New("Value cant be empty")
	}
	return s, nil
}

// Compute where your app’s config file should live on the user’s machine, and ensure the folder exists.
func configPath() (string, error) {
	dir, err := os.UserConfigDir() // ? This will return the rootDir => C:\Users\Rituraj\AppData\Roaming
	if err != nil{
		return "", err
	}

	dir = filepath.Join(dir, "oops") // ? C:\Users\Rituraj\AppData\Roaming\oops
	if err := os.MkdirAll(dir, 0o700); err != nil{ // ? This will create a new directory named - oops
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil // ? => C:\Users\Rituraj\AppData\Roaming\oops\config.json
}

// Read the config file (if it exists) and decode it into a Config struct.
func LoadConfig() (*Config, string, error){
	path, err := configPath() // ? will return the path of the config.json file
	if err != nil{
		return nil, "", err
	}

	data, err := os.ReadFile(path);
	if err != nil{
		if os.IsNotExist(err){
			return nil, path, err
		}
		return nil, path, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil{
		return nil, path, err
	}
	return &cfg, path, nil
}

// Save your config safely (avoids half-written config files).
func SaveConfig(path string, cfg *Config) error {
	dir := filepath.Dir(path)

	tmp, err := os.CreateTemp(dir, "config-*.json")
	if err != nil{
		return err
	}
	tmpName := tmp.Name()
	defer func(){
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}()

	enc := json.NewEncoder(tmp)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil{
		return err
	}

	if err := tmp.Close(); err != nil{
		return err
	}

	_ = os.Chmod(tmpName, 0o600)

	return os.Rename(tmpName, path)
}