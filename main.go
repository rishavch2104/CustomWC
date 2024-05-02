package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Processor interface {
	process(file *os.File) error
	getValue() int
}

type ValueGetter struct {
	value int
}

func (v ValueGetter) getValue() int {
	return v.value
}

type ByteCountProcessor struct {
	ValueGetter
}

type LineCountProcessor struct {
	ValueGetter
}

type WordCountProcessor struct {
	ValueGetter
}

func (processor *WordCountProcessor) process(file *os.File) error {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		processor.value += len(words)
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return errors.New("unable to parse file")
	}
	return nil

}

func (processor *LineCountProcessor) process(file *os.File) error {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		processor.value++
	}

	if err := scanner.Err(); err != nil {
		return errors.New("unable to parse file")
	}
	return nil
}

func (processor *ByteCountProcessor) process(file *os.File) error {
	fileInfo, err := file.Stat()
	if err != nil {
		return errors.New("unable to parse file")

	}
	fileSize := fileInfo.Size()
	processor.value = int(fileSize)
	return nil
}

func main() {
	var countBytesFlag bool
	var countLinesFlag bool
	var countWordsFlag bool
	var fileName string
	var file *os.File
	defer file.Close()

	flag.BoolVar(&countBytesFlag, "c", false, "Count Bytes Flag")
	flag.BoolVar(&countLinesFlag, "l", false, "Count Lines Flag")
	flag.BoolVar(&countWordsFlag, "w", false, "Count Words Flag")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Print("Filename missing!")
		return
	}
	fileName = flag.Args()[0]
	file, err := os.Open(fileName)

	if err != nil {
		fmt.Print("File not found")
		return
	}

	var processor Processor
	if countBytesFlag {
		processor = &ByteCountProcessor{}
	}
	if countLinesFlag {
		processor = &LineCountProcessor{}
	}
	if countWordsFlag {
		processor = &WordCountProcessor{}
	}

	err = processor.process(file)

	if err != nil {
		fmt.Print("Unable to print file")
		return
	}
	fmt.Printf("%d %s \n", processor.getValue(), fileName)

}
