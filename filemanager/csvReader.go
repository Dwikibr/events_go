package filemanager

import (
	"encoding/csv"
	"os"
)

type resList []map[string]interface{}

func ReadCsvFile(filePath string) resList {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	record, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var columnMap = make(map[string]int)
	header := record[0]
	for idx, col := range header {
		columnMap[col] = idx
	}

	var result resList
	for _, row := range record[1:] {
		rowMap := make(map[string]interface{})
		for col, idx := range columnMap {
			rowMap[col] = row[idx]
		}
		result = append(result, rowMap)
	}

	return result
}
