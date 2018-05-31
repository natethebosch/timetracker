package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

const sinceDateFormat = "Jan-02"

func main() {

	yesterday := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -1).Format(sinceDateFormat)

	print := flag.Bool("print", false, "Prints out records in the database")
	since := flag.String("since", yesterday, "Limits by date the records to print [only applies when -print is enabled]")

	flag.Parse()

	configFile := resolveFileName("~/.timetrack.json")
	cfg, err := loadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	if *print {
		printRecords(since, cfg)
		return
	}

	track(configFile, cfg)
}

func printRecords(since *string, cfg *Config) {

	tm, err := time.Parse(sinceDateFormat, *since)
	if err != nil {
		log.Printf("Date should be in format %s\n", sinceDateFormat)
		log.Fatal(err)
	}

	tm = tm.Truncate(24 * time.Hour)

	// add year component b/c it wasn't parsed
	tm = tm.AddDate(time.Now().Year(), 0, 0)

	fmt.Printf("Showing results since %s\n", tm.Format("2006-Jan-02 15:04"))
	fmt.Println()

	sort.Sort(byFinishedTime(cfg.Entries))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Time", "Duration", "Description"})

	for i, entry := range cfg.Entries {
		if !entry.Finished.After(tm) {
			continue
		}

		var lastTime time.Time
		if i > 0 {
			lastTime = cfg.Entries[i-1].Finished
		}

		timeSpent := entry.Finished.Sub(lastTime).Round(time.Minute)
		timeSpentString := timeSpent.String()
		if timeSpent > 8*time.Hour {
			timeSpentString = "-"
		}

		finished := entry.Finished.Format("Jan-02 15:04")
		table.Append([]string{finished, timeSpentString, entry.Description})
	}

	table.Render()
}

func track(configFile string, cfg *Config) {

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

	err := writeConfig(configFile, cfg)
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

type byFinishedTime []*Entry

func (a byFinishedTime) Len() int           { return len(a) }
func (a byFinishedTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFinishedTime) Less(i, j int) bool { return a[i].Finished.Before(a[j].Finished) }

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

	if len(data) == 0 {
		return cfg, nil
	}

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
