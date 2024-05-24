package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var outputFormats = []string{"azw3", "epub", "docx", "fb2", "htmlz", "oeb", "lit", "lrf", "mobi", "pdb", "pmlz", "rb", "pdf", "rtf", "snb", "tcr", "txt", "txtz", "zip"}
var hintLevel int = 5
var formatType string
var inputPath string
var wLang string = "en"
var including_phoneme bool = false
var defLenth int = 1
var isDirectRun bool = false

func main() {
	readInputParams(os.Args)

	log.Println("[+] Input path:", inputPath)
	log.Println(fmt.Sprintf("[+] Hint level: %d, Output format type: %s, Language: %s, Phoneme: %t, Definetion lenth: %d", hintLevel, formatType, wLang, including_phoneme, defLenth))

	log.Println("[+] Load wordwise dict")
	loadWordwiseDict()

	log.Println("[+] Load lemma dict")
	loadLemmatizerDict()

	// clean old temp
	log.Println("[+] Cleaning old temp files")
	cleanTempData()

	// get ebook convert cmd
	ebookConvertCmd = getEbookConvertCmd()

	// convert book to html
	createTempFolder()
	convertBookToHtml(inputPath)

	// process book
	log.Println("[+] Process book with wordwise")
	processHtmlBookData()
	modifyCalibreTitle()

	// create wordwise book
	createBookWithWordwised(inputPath)

	cleanTempData()
	if isDirectRun {
		pauseConsole()
	}
}

func readInputParams(args []string) {
	if len(args) < 2 {
		isDirectRun = true
		readInputFromConsole()
	} else {
		assignInputPath(args[1])

		if len(args) > 2 {
			assignHintLevel(args[2])
		}

		if len(args) > 3 {
			assignOutputFormat(args[3])
		} else {
			assignOutputFormat("")
		}

		if len(args) > 4 {
			assignLanguage(args[4])
		}

		if len(args) > 5 {
			assignPhoneme(args[5])
		}

		if len(args) > 6 {
			assignDefinitionLength(args[6])
		}
	}
}

func readInputFromConsole() {
	checkThenChangeWorkingDir()

	userInput := bufio.NewReader(os.Stdin)
	log.Println("Enter the book's path OR drag n drop your book here:")
	fmt.Print("                    ")
	scanValue, _ := userInput.ReadString('\n')
	scanValue = strings.TrimSpace(scanValue)
	assignInputPath(scanValue)

	log.Println("Enter hint level (1-5): ")
	fmt.Print("                    ")
	scanValue, _ = userInput.ReadString('\n')
	scanValue = strings.TrimSpace(scanValue)
	assignHintLevel(scanValue)

	log.Println("Enter output format (not support \"mobi\"): ")
	fmt.Print("                    ")
	scanValue, _ = userInput.ReadString('\n')
	scanValue = strings.TrimSpace(scanValue)
	assignOutputFormat(scanValue)

	log.Println("Enter language (\"en\", \"vi\", \"cn\", \"ar\", \"de\", \"es\", \"fr\", \"hi\", \"jp\", \"ko\", \"pt\", \"ru\", \"th\", \"ua\"): ")
	fmt.Print("                    ")
	scanValue, _ = userInput.ReadString('\n')
	scanValue = strings.TrimSpace(scanValue)
	assignLanguage(scanValue)

	log.Println("Including Phoneme? (y/n): ")
	fmt.Print("                    ")
	scanValue, _ = userInput.ReadString('\n')
	scanValue = strings.TrimSpace(scanValue)
	assignPhoneme(scanValue)

	log.Println("Enter the length of wordwise(0: disable the definition, 1: short description, 2: long description)? (0,1,2): ")
	fmt.Print("                    ")
	scanValue, _ = userInput.ReadString('\n')
	scanValue = strings.TrimSpace(scanValue)
	assignDefinitionLength(scanValue)
}

func assignInputPath(scanValue string) {
	inputPath = strings.ReplaceAll(strings.Trim(scanValue, "\""), "\\ ", " ")
	if _, err := os.Stat(inputPath); err != nil {
		logFatalln(fmt.Sprintf("File at %s is not found!", inputPath))
	}
}

func assignHintLevel(scanValue string) {
	parseNum, err := strconv.Atoi(scanValue)
	if err == nil {
		hintLevel = parseNum
	}
}

func assignOutputFormat(scanValue string) {
	if contains(outputFormats, scanValue) {
		formatType = scanValue
	} else {
		extension := strings.Trim(filepath.Ext(inputPath), ".")
		if contains(outputFormats, extension) {
			formatType = extension
		} else {
			formatType = "epub"
		}
	}
}

func assignLanguage(scanValue string) {
	wLang = scanValue
}

func assignPhoneme(scanValue string) {
	including_phoneme = false

	if scanValue == "y" || scanValue == "yes" {
		including_phoneme = true
	} else if scanValue == "n" || scanValue == "no" {
		including_phoneme = false
	}
}

func assignDefinitionLength(scanValue string) {
	parseNum, err := strconv.Atoi(scanValue)
	if err == nil {
		if defLenth >= 0 && defLenth <= 2 {
			defLenth = parseNum
		} else {
			defLenth = 1
		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func logFatalln(v ...any) {
	log.Println(v...)
	if isDirectRun {
		pauseConsole()
	}
	os.Exit(1)
}

func pauseConsole() {
	log.Println("Press the Enter Key to exit!")
	fmt.Scanln()
}

func checkThenChangeWorkingDir() {
	isFoundResources := true
	if _, err := os.Stat(LemmaDictionaryPath); err != nil {
		isFoundResources = false
	}
	if _, err := os.Stat(WordwiseDictionaryPath); err != nil {
		isFoundResources = false
	}
	if !isFoundResources {
		execPath, _ := os.Executable()
		workingDir := filepath.Dir(execPath)
		os.Chdir(workingDir)
	}
}
