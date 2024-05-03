package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

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
