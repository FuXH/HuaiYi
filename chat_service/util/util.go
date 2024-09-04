package util

import (
	"bufio"
	"io"
	"os"
)

func WriteFile(fileName string, data []string) error {
	file, err := os.Create(fileName)
	//file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	for _, val := range data {
		_, _ = file.WriteString(val)
	}
	return nil
}

func ReadFile(fileName string) (string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ReadFileByLine(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	res := make([]string, 0)
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		res = append(res, string(line))
	}
	return res
}
