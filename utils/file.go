package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func FileToLines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

/**
 * Insert sting to n-th line of file.
 * If you want to insert a line, append newline '\n' to the end of the string.
 */
func InsertStringToFile(path, str string, index int) error {
	lines, err := FileToLines(path)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\n"
	}

	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}

func IsFileExisted(filePath string) (bool, error) {
	_, err := os.Stat(filePath)

	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("file does not exist")
		return false, nil
	} else if err == nil {
		fmt.Println("file exists")
		return true, nil
	}
	return false, err
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func CreateFile(filePath string) error {
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func WriteFile(filePath string, content string) error {
	//write the slice of bytes to the given filename with the 0644 permissions
	err := ioutil.WriteFile(filePath, []byte(content), 0644)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

/**
 * Read lines from the end of the file
 */
func ReadFileByLinesFromEnd(filePath string, maxLines int) ([]string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// Buffer to store lines in reverse order
	var lines []string
	var currentLine strings.Builder

	// Start from the end of the file
	pos := fileSize - 1

	// Buffer for reading chunks
	buffer := make([]byte, 1)

	for pos >= 0 {
		// Seek to position
		_, err := file.Seek(pos, 0)
		if err != nil {
			return nil, err
		}

		// Read one byte
		_, err = file.Read(buffer)
		if err != nil {
			return nil, err
		}

		// Handle newline
		if buffer[0] == '\n' {
			if currentLine.Len() > 0 {
				lines = append(lines, ReverseString(currentLine.String()))
				currentLine.Reset()
				if maxLines > 0 && len(lines) >= maxLines {
					break
				}
			}
		} else if buffer[0] != '\r' { // Skip carriage returns
			currentLine.WriteByte(buffer[0])
		}

		pos--
	}

	// Add the first line if there's content
	if currentLine.Len() > 0 {
		lines = append(lines, ReverseString(currentLine.String()))
	}

	return lines, nil
}
