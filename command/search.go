package command

import (
	"unicode/utf8"
	"encoding/json"
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"flag"
	"strings"
)

var (
	NAME = "search"
	HOME = os.Getenv("HOME")
	CONF = ".godic"
	READER_SIZE = 4096
	MAX_WROD = 3
	DEFAULT_CASING_FLAG = ""
	DEFAULT_PROJECT_FLAG = -1
	URL = "https://api.codic.jp/v1/engine/translate.json"
)

const (
	ExitCodeOK = iota
	ExitCodeInternalError
)

type Result struct {
	Successful bool
	Text string
	TranslatedText string `json:"translated_text"`
	Words []Word
}

type Word struct {
	Successful bool
	Text string
	TranslatedText string `json:"translated_text"`
	Candidates []Candidate
}

type Candidate struct {
	Text string
}

type SearchCommand struct {
	Meta
}

func maxTextWidth(results []Result) int {
	width := 0
	for _, w := range results {
		count := utf8.RuneCountInString(w.Text)
		if count > width {
			width = count
		}
	}
	return width
}

func (c *SearchCommand) Run(args []string) int {
	// Parse subcommand flags
	var (
		projectFlag int
		casingFlag string
		helpFlag bool
		words []string
	)

	flags := flag.NewFlagSet(NAME, flag.ContinueOnError)
	flags.BoolVar(&helpFlag, "h", false, "")
	flags.BoolVar(&helpFlag, "help", false, "")
	flags.StringVar(&casingFlag, "c", "", "")
	flags.StringVar(&casingFlag, "casing", "", "")
	flags.IntVar(&projectFlag, "p", DEFAULT_PROJECT_FLAG, "")
	flags.IntVar(&projectFlag, "project", DEFAULT_PROJECT_FLAG, "")

	if err := flags.Parse(args); err != nil {
		return ExitCodeInternalError
	}
	words = flags.Args()

	// Get AccessToken
	fp, err := os.Open(path.Join(HOME, CONF))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not open configuration file.\n")
		return ExitCodeInternalError
	}
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, READER_SIZE)
	line, _, err := reader.ReadLine()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not read AccessToken.\n")
		return ExitCodeInternalError
	}
	token := string(line)

	// Setup search query
	if helpFlag || len(words) == 0 || len(words) > MAX_WROD {
		fmt.Println(c.Help())
		return ExitCodeInternalError
	}
	
	values := url.Values{}
	values.Add("text", strings.Join(words, "\n"))

	switch casingFlag {
	case DEFAULT_CASING_FLAG:
	case "camel", "pascal", "hyphen":
		values.Add("casing", casingFlag)
	case "lower_underscore", "upper_underscore":
		values.Add("casing", strings.Replace(casingFlag, "_", " ", -1))
	default:
		fmt.Fprintf(os.Stderr, "Casing type is invalid.\n")
		return ExitCodeInternalError
	}

	if projectFlag != DEFAULT_PROJECT_FLAG {
		values.Add("project_id", fmt.Sprint(projectFlag))
	}

	// Throw query
	req, err := http.NewRequest("GET", URL + "?" + values.Encode(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute: %s\n", err.Error())
		return ExitCodeInternalError
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute: %s\n", err.Error())
		return ExitCodeInternalError
	}
	defer res.Body.Close()

	// Parse result
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute: %s\n", err.Error())
		return ExitCodeInternalError
	}

	results := make([]Result, 0)
	err = json.Unmarshal(body, &results)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute: %s\n", err.Error())
		return ExitCodeInternalError
	}

	// Show results
	wordWidth := maxTextWidth(results)
	space := wordWidth - 2
	if space < 0 {
		space = 0
	}
	fmt.Printf("TEXT%s TRANSLATE\n", strings.Repeat("　", space))

	var wordMap map[string]bool = map[string]bool{}
	for _, result := range results {
		translate := result.TranslatedText
		if !result.Successful {
			translate = "FAILED"
		}
		wordMap[result.Text] = true
		fmt.Printf("%s%s %s\n",
			result.Text,
			strings.Repeat("　", wordWidth - utf8.RuneCountInString(result.Text)),
			translate)
	}

	fmt.Printf("\nWORD%s TRANSLATE\n", strings.Repeat("　", space))
	for _, result := range results {
		for _, word := range result.Words {
			if !word.Successful || wordMap[word.Text] {
				continue
			}

			wordMap[word.Text] = true
			fmt.Printf("%s%s ",
				word.Text,
				strings.Repeat("　", wordWidth - utf8.RuneCountInString(word.Text)))
			
			for i, candidate := range word.Candidates {
				if i != 0 {
					fmt.Print(", ")
				}
				fmt.Print(candidate.Text)
			}
			fmt.Print("\n")
		}
	}
	
	return ExitCodeOK
}

func (c *SearchCommand) Synopsis() string {
	return "search function name"
}

func (c *SearchCommand) Help() string {
	helpText := `
Usage:
  godic search [OPTIONS] <word..>

  Japanes-English translation.

Example:
  godic search --casing pascal ユーザ削除 登録

Word:
  Set string(Japanes) to be trancelated.
  Specify up to three words separated by spaces.

Options:
  -p --project ID    Specifies the project ID to be used in the translation.
  -c --casing TYPE   Valid TYPE is camel, pascal, lower_underscore, upper_underscore and hyphen.
  -h --help          Show this.
        `
	return strings.TrimSpace(helpText)
}
