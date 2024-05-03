package main

import (
	"bufio"
	"errors"
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
		return errors.New("error rewinding temporary file")
	}
	return nil
}
