package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	WordwiseDictionaryPath = "resources/wordwise-dict.%s.csv"
	LemmaDictionaryPath    = "resources/lemmatization-en.csv"
)

var wordwiseDict *map[string]DictRow
var lemmaDict *map[string]string

type DictRow struct {
	Word             string
	Phoneme          string
	Full_Def         string
	Short_Def        string
	ExampleSentences string
	HintLvl          int
}

func (ws *DictRow) meaning(including_phoneme bool) string {
	definition := ""
	if including_phoneme && len(ws.Phoneme) > 0 {
		definition += ws.Phoneme
	}

	if defLenth == 1 {
		definition += " " + ws.Short_Def
	} else if defLenth == 2 {
		definition += " " + ws.Full_Def
	}

	return strings.TrimSpace(definition)
}

// Load Dict from CSV
func loadWordwiseDict() {
	wordwise_dict_path := fmt.Sprintf(WordwiseDictionaryPath, wLang)

	file, err := os.Open(wordwise_dict_path)
	if err != nil {
		logFatalln("Error when open ", wordwise_dict_path, "->", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	reader := csv.NewReader(file)

	dict := make(map[string]DictRow)
	row := DictRow{}

	// Read each record from csv
	// skip header
	record, err := reader.Read()
	if err == io.EOF {
		logFatalln("Empty csv file")
	}

	count := 0

	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logFatalln("Error when scan word ", count, "->", err)
		}

		if len(record) < 5 {
			log.Println("Invalid word: ", record)
			continue
		}

		hintLv, err := strconv.Atoi(record[6])
		if err != nil {
			log.Println("Can't get hint_level: ", record, "->", err)
			continue
		}

		row = DictRow{
			Word:             record[1],
			Phoneme:          record[2],
			Full_Def:         record[3],
			Short_Def:        record[4],
			ExampleSentences: record[5],
			HintLvl:          hintLv,
		}

		dict[row.Word] = row
		count++
	}

	log.Println("--> Csv words:", count)
	wordwiseDict = &dict
}

// Load Dict from CSV
func loadLemmatizerDict() {

	file, err := os.Open(LemmaDictionaryPath)
	if err != nil {
		logFatalln("Error when open ", LemmaDictionaryPath, "->", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	reader := csv.NewReader(file)

	dict := make(map[string]string)

	var record []string
	count := 0
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logFatalln("Error when scan word ", count, "->", err)
		}

		if len(record) < 2 {
			log.Println("Invalid word: ", record)
			continue
		}

		dict[record[1]] = record[0]
		count++
	}

	log.Println("--> Lemma pairs:", count)
	lemmaDict = &dict
}
