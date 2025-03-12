package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func dataprocessingfunction() {
	//

	fmt.Println("Beginning process")
	basefile, err := excelize.OpenFile("final_curvature.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := basefile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetnum := "600 N"
	//findsheetName(sheetnum)
	er2 := basefile.SetSheetVisible(sheetnum, true)
	if er2 != nil {
		fmt.Println(er2)
		return
	}
	fmt.Println("Before processing logic")

	//do it here, move into func later
	//1. process the script into 1/3rd of the normal size
	//downsample by 33%
	//	newCSV, _ := os.Create("exampleDownsample.csv")

}
func findsheetName(sheetnum int, thisname string) (sheetname string) {

	nametoint, casterr := strconv.Atoi(thisname)
	if casterr != nil {
		panic(casterr)
	}

	if sheetnum > 2400 {
		nametoint += 200
		finalString := strconv.Itoa(nametoint)
		return finalString + " " + "N"
	}
	return thisname
}
