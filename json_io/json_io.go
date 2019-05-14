package json_io

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadJSON(path string) ([]map[string]interface{}, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	byteStream, _ := ioutil.ReadFile(path)

	var bowOut []map[string]interface{}

	errJSON := json.Unmarshal(byteStream, &bowOut)

	if errJSON != nil {
		return nil, errJSON
	}

	return bowOut, nil
}

func ReadJsonRaw(path string) ([]string, error) {
	lines := make([]string, 1)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
