package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
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

type CharacterCountProcessor struct {
	ValueGetter
}

func (processor *CharacterCountProcessor) process(file *os.File) error {
	err := resetFile(file)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	for {
		_, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		processor.value++

	}

	return nil
}

func (processor *WordCountProcessor) process(file *os.File) error {
	err := resetFile(file)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		processor.value += len(words)
	}

	if err := scanner.Err(); err != nil {
		return errors.New("unable to parse file")
	}
	return nil

}

func (processor *LineCountProcessor) process(file *os.File) error {
	err := resetFile(file)
	if err != nil {
		return err
	}
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

func resetFile(file *os.File) error {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return errors.New("Error rewinding temporary file:")
	}
	return nil
}

func main() {
	var countBytesFlag bool
	var countLinesFlag bool
	var countWordsFlag bool
	var countCharactersFlag bool
	var fileName string
	var file *os.File
	var err error
	defer file.Close()

	flag.BoolVar(&countBytesFlag, "c", false, "Count Bytes Flag")
	flag.BoolVar(&countLinesFlag, "l", false, "Count Lines Flag")
	flag.BoolVar(&countWordsFlag, "w", false, "Count Words Flag")
	flag.BoolVar(&countCharactersFlag, "m", false, "Count Characters Flag")
	flag.Parse()
	if len(flag.Args()) == 0 {
		file, err = os.CreateTemp("", "stdin-")
		if err != nil {
			fmt.Println("Error creating temporary file:", err)
			return
		}
		_, err = io.Copy(file, os.Stdin)
		if err != nil {
			fmt.Println("Error copying data to temporary file:", err)
			return
		}
	} else {
		fileName = flag.Args()[0]
		file, err = os.Open(fileName)
		if err != nil {
			fmt.Print("File not found")
			return
		}

	}
	var processors []Processor
	if countBytesFlag {
		processors = append(processors, &ByteCountProcessor{})
	}
	if countLinesFlag {
		processors = append(processors, &LineCountProcessor{})
	}
	if countWordsFlag {
		processors = append(processors, &WordCountProcessor{})
	}
	if countCharactersFlag {
		processors = append(processors, &CharacterCountProcessor{})
	}

	if len(processors) == 0 {
		processors = append(processors, &LineCountProcessor{})
		processors = append(processors, &WordCountProcessor{})
		processors = append(processors, &ByteCountProcessor{})
	}

	var output string
	for _, processor := range processors {
		err = processor.process(file)
		if err != nil {
			fmt.Print("Unable to print file")
			return
		}
		output = output + fmt.Sprintf("%d", processor.getValue()) + "\t"
	}

	fmt.Printf("%s %s \n", output, fileName)

}
