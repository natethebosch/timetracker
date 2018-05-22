package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	configFile := resolveFileName("~/.timetrack.json")
	cfg, err := loadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	lastTime := cfg.GetLastEntryTime()
	fmt.Printf("%s since your last entry\n", time.Now().Sub(lastTime).String())

	var desc string
	for len(desc) == 0 {
		desc = getDesc()
	}

	cfg.Entries = append(cfg.Entries, &Entry{
		Finished:    time.Now(),
		Description: desc,
	})

	err = writeConfig(configFile, cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Saved %s\n", configFile)
}

func resolveFileName(fileName string) string {
	if strings.HasPrefix(fileName, "~/") {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		fileName = filepath.Join(usr.HomeDir, fileName[2:])
	}

	return fileName
}

type Config struct {
	Entries []*Entry
}

func (c *Config) GetLastEntryTime() time.Time {

	maxTime := time.Time{}

	for _, entry := range c.Entries {
		if entry.Finished.After(maxTime) {
			maxTime = entry.Finished
		}
	}

	if maxTime.IsZero() {
		return time.Now()
	}

	return maxTime
}

type Entry struct {
	Finished    time.Time
	Description string
}

func loadConfig(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {

		if os.IsNotExist(err) {
			fmt.Printf("Creating config file %s\n", fileName)

			if err = createConfig(fileName); err != nil {
				return nil, err
			}

			return &Config{}, nil
		}

		return nil, errors.New("Load Config: " + err.Error())
	}

	cfg := &Config{}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, errors.New("Load Config: " + err.Error())
	}

	return cfg, nil
}

func createConfig(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return errors.New("Create Config: " + err.Error())
	}

	defer f.Close()
	return nil
}

func writeConfig(fileName string, cfg *Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return errors.New("Write Config: " + err.Error())
	}

	err = ioutil.WriteFile(fileName, data, os.ModePerm)
	if err != nil {
		return errors.New("Write Config: " + err.Error())
	}

	return nil
}

func getUserString() string {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
	}

	return strings.Trim(text, " \n\t")
}

func getDesc() string {
	fmt.Println("Enter a task description:")
	return getUserString()
}
