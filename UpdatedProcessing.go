package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func OverwriteThisFile(thisstring string) {
	newfile, err := os.Open(thisstring)
	if err != nil {
		panic("file not found")
	}
	r := csv.NewReader(newfile)
	r.Comma = ' '
	//r.FieldsPerRecord = 3
	r.FieldsPerRecord = -1
	data, err := r.ReadAll()
	fmt.Println("\n")
	if err == csv.ErrFieldCount {
		r.FieldsPerRecord = -1
		fmt.Printf("WARNING! FIELDS PER RECORD AINT 3")
	} else if err != nil {
		fmt.Println("ERROR reading CSV" + err.Error())
	}
	os.Remove("newcsv.csv")
	file, err := os.Create("newcsv.csv")
	if err != nil {
		panic(err)
	}
	mywriter := csv.NewWriter(file)
	mywriter.Comma = ','

	for _, row := range data {
		checker := true
		for _, col := range row {
			if len(col) == 3 {
				//fmt.Println("False found")
				checker = false
				fmt.Printf(col)

			}
		}

		if checker {
			mywriter.Write(row)
			mywriter.Flush()
		} else if !checker {
			if row[2] == "nan" {
			} else {
				fmt.Printf("Bad news. We've found a value that should maybe be read? \n")
				mywriter.Write(row)
				mywriter.Flush()
				//for g, _ := range row {
				//	fmt.Printf(row[g])
				//}
			}
		}
		os.Remove(newfile.Name())
		os.Create(thisstring)

		input, err := os.ReadFile("newcsv.csv")
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.WriteFile(thisstring, input, 0644)
		if err != nil {
			fmt.Println("Error creating", thisstring)
			fmt.Println(err)
			return
		}
	}
}

// SelectRandomRowsToSize takes a source CSV file, a target file size in bytes,
// and creates a new CSV file by randomly selecting rows until the target size is reached.
func SelectRandomRowsToSize(sourceFile string, targetSize int64, outputFile string) error {
	// Open the source CSV file
	file, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer file.Close()

	// Read the CSV file into memory
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %v", err)
	}

	// Open the output file
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer output.Close()

	writer := csv.NewWriter(output)
	defer writer.Flush()

	// Write the header row if present
	if len(rows) > 0 {
		err = writer.Write(rows[0])
		if err != nil {
			return fmt.Errorf("failed to write header: %v", err)
		}
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Keep selecting random rows until the file reaches the target size
	var currentSize int64
	for currentSize < targetSize {
		// Randomly select a row (skip the header row if it exists)
		row := rows[rand.Intn(len(rows)-1)+1]

		// Write the row to the output file
		err = writer.Write(row)
		if err != nil {
			return fmt.Errorf("failed to write row: %v", err)
		}

		// Flush the writer to disk to get the current file size
		writer.Flush()
		fi, err := output.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file size: %v", err)
		}
		currentSize = fi.Size()
	}

	fmt.Printf("Successfully created a CSV file with approximately %d bytes.\n", currentSize)
	return nil
}
