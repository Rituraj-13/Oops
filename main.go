package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
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
		Width(30)

	// Trying to get update data using Flags (fieldSelector + value)
	setField := flag.String("set", "", `field to update "name" or "apiKey"`)
	setValue := flag.String("value", "", "new value for the field provided by --set")
	flag.Parse()

	if strings.TrimSpace(*setField) != "" {
		if err := UpdateConfig(*setField, *setValue, style); err != nil {
			fmt.Fprintln(os.Stderr, "Update Failed:", err)
			os.Exit(1)
		}
	}

	cfg, cfgPath, err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load config: ", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	// For 1st run
	if cfg == nil || strings.TrimSpace(cfg.APIKey) == "" || strings.TrimSpace(cfg.Name) == "" {
		fmt.Println(style.Render("First time setup"))

		name, err := PromptLine(reader, style, "Enter your name: ")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read name:", err)
			os.Exit(1)
		}

		apiKey, err := PromptLine(reader, style, "Enter your Gemini API Key: ")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read API Key:", err)
			os.Exit(1)
		}

		cfg = &Config{Name: name, APIKey: apiKey}
		if err := SaveConfig(cfgPath, cfg); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to save config:", err)
			os.Exit(1)
		}

		fmt.Println(style.Render("Saved config to: " + cfgPath))
		return
	}
	fmt.Println(style.Render("Loaded config for: " + cfg.Name))

}

// update the name/apikey of the user
func UpdateConfig(setField string, setValue string, style lipgloss.Style) error {

	field := strings.ToLower(strings.TrimSpace(setField))
	value := strings.TrimSpace(setValue)

	if value == "" {
		flag.Usage()
		return errors.New("--value cannot be empty")
	}

	cfg, cfgPath, err := LoadConfig()

	if err != nil {
		fmt.Println(style.Render("Couldn't load config !"))
		return err
	}

	if cfg == nil {
		cfg = &Config{}
	}

	switch field {
	case "name":
		cfg.Name = value
	case "apikey":
		cfg.APIKey = value
	default:
		flag.Usage()
		return fmt.Errorf(`unknown field %q (use "name" or "apikey")`, setField)
	}

	if err := SaveConfig(cfgPath, cfg); err != nil {
		return err
	}
	fmt.Println(style.Render("Config Updated! Saved to: " + cfgPath))
	return nil
}

// Ask the user for one line of input in the terminal and return it (trimmed).
func PromptLine(reader *bufio.Reader, style lipgloss.Style, label string) (string, error) {
	fmt.Printf((style.Render(label)))
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
	if err != nil {
		return "", err
	}

	dir = filepath.Join(dir, "oops")                // ? C:\Users\Rituraj\AppData\Roaming\oops
	if err := os.MkdirAll(dir, 0o700); err != nil { // ? This will create a new directory named - oops
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil // ? => C:\Users\Rituraj\AppData\Roaming\oops\config.json
}

// Read the config file (if it exists) and decode it into a Config struct.
func LoadConfig() (*Config, string, error) {
	path, err := configPath() // ? will return the path of the config.json file
	if err != nil {
		return nil, "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, path, nil
		}
		return nil, path, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, path, err
	}
	return &cfg, path, nil
}

// Save your config safely (avoids half-written config files).
func SaveConfig(path string, cfg *Config) error {
	dir := filepath.Dir(path)

	tmp, err := os.CreateTemp(dir, "config-*.json")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}()

	enc := json.NewEncoder(tmp)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	_ = os.Chmod(tmpName, 0o600)

	return os.Rename(tmpName, path)
}
