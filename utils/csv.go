package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func CsvLoop(filename string, process func([]string) error) error {
	// 打开CSV文件
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// 创建一个CSV读取器
	csvReader := csv.NewReader(f)

	// 逐行读取数据
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("csv read: %w", err)
		}
		err = process(record)
		if err != nil {
			return fmt.Errorf("csv process: %w", err)
		}
	}
	return nil
}
