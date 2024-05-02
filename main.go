package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
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
	var fileName string
	var file *os.File
	defer file.Close()

	flag.BoolVar(&countBytesFlag, "c", false, "Count Bytes Flag")
	flag.BoolVar(&countLinesFlag, "l", false, "Count Lines Flag")
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
	err = processor.process(file)

	if err != nil {
		fmt.Print("Unable to print file")
		return
	}
	fmt.Printf("%d %s \n", processor.getValue(), fileName)

}
