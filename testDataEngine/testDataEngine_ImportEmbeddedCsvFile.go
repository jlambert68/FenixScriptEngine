package testDataEngine

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
)

// ImportEmbeddedSimpleCsvTestDataFile
// Imports an embedded csv-file with relative path and name in 'fileNameAndRelativePath'
// and having a data divider of type 'divider'
// The first row must consist of column headers
func ImportEmbeddedSimpleCsvTestDataFile(
	embeddedFileAsByteArray []byte,
	//fileNameAndRelativePath string,
	divider rune) (
	testDataFromTestDataArea TestDataFromSimpleTestDataAreaStruct) {

	var err error

	// This is the structure of the "simple" csv-version of the testdata
	// i.e. the file should follow this structure
	// Domain and Area info consists of only one value per row
	var testDataHeadersInCsv []string
	var testDataDomainUuid []string
	var testDataDomainName []string
	var testDataAreaUuid []string
	var testDataAreaName []string
	var testDataHeadersUsedInFiltersInCsv []string
	var testDataRows [][]string

	var testDataHeadersUsedInFiltersInCsvMap map[string]bool
	var testDataHeaders []struct {
		ShouldHeaderActAsFilter bool
		HeaderName              string
		HeaderUiName            string
	}

	// Read the embedded file
	//data, err := embeddedFilePtr.ReadFile(fileNameAndRelativePath)
	//if err != nil {
	//	log.Fatalf("Error reading the embedded file: %v", err)
	//}

	// Parse the CSV file
	r := csv.NewReader(bytes.NewReader(embeddedFileAsByteArray))
	r.Comma = divider

	// Read the headers as 1st row
	testDataHeadersInCsv, err = r.Read()
	if err != nil {
		log.Fatalf("Error reading headers as 1st row: %v", err)
	}

	// Read the TestDataDomainUuid as 2nd row
	testDataDomainUuid, err = r.Read()
	if err != nil && err.Error() != "record on line 2: wrong number of fields" {
		log.Fatalf("Error reading TestDataDomainUuid row as 2nd row: %v", err)
	}

	// Read the TestDataDomainName as 3rd row
	testDataDomainName, err = r.Read()
	if err != nil && err.Error() != "record on line 3: wrong number of fields" {
		log.Fatalf("Error reading TestDataDomainName row as 3rd row: %v", err)
	}

	// Read the TestDataAreaUuid as 4th row
	testDataAreaUuid, err = r.Read()
	if err != nil && err.Error() != "record on line 4: wrong number of fields" {
		log.Fatalf("Error reading TestDataAreaUuid row as 4th row: %v", err)
	}

	// Read the TestDataAreaName as 5th row
	testDataAreaName, err = r.Read()
	if err != nil && err.Error() != "record on line 5: wrong number of fields" {
		log.Fatalf("Error reading TestDataAreaName row as 5th row: %v", err)
	}

	// Read the header filters as 6th row
	testDataHeadersUsedInFiltersInCsv, err = r.Read()
	if err != nil && err.Error() != "record on line 6: wrong number of fields" {
		log.Fatalf("Error reading Headerfilter row as 6th row: %v", err)
	}

	// Iterate through the records and extract rows
	for {
		rowRecord, errOrEOF := r.Read()

		// Check if we reach end of file
		if errOrEOF == io.EOF {
			break
		}

		// Check for error
		if errOrEOF != nil {
			log.Fatal(err)
		}

		// Loop all records in  row and extract them
		var testDataRow []string
		for _, recordItem := range rowRecord {
			testDataRow = append(testDataRow, recordItem)
		}

		// Add row to TestDataRows
		testDataRows = append(testDataRows, testDataRow)

	}

	// Create a Map with the headers that should be part of filter when searching TestData
	testDataHeadersUsedInFiltersInCsvMap = make(map[string]bool)
	for _, testDataHeaderUsedInFilter := range testDataHeadersUsedInFiltersInCsv {
		testDataHeadersUsedInFiltersInCsvMap[testDataHeaderUsedInFilter] = true
	}

	// Convert Headers from CSV into TestData struct structure
	for _, testDataHeader := range testDataHeadersInCsv {
		var tempTestDataHeader struct {
			ShouldHeaderActAsFilter bool
			HeaderName              string
			HeaderUiName            string
		}

		tempTestDataHeader.HeaderName = testDataHeader
		tempTestDataHeader.HeaderUiName = testDataHeader
		tempTestDataHeader.ShouldHeaderActAsFilter = testDataHeadersUsedInFiltersInCsvMap[testDataHeader]

		testDataHeaders = append(testDataHeaders, tempTestDataHeader)
	}

	// Create full TestDataFromTestDataArea-object
	testDataFromTestDataArea = TestDataFromSimpleTestDataAreaStruct{
		TestDataDomainUuid: testDataDomainUuid[0],
		TestDataDomainName: testDataDomainName[0],
		TestDataAreaUuid:   testDataAreaUuid[0],
		TestDataAreaName:   testDataAreaName[0],
		Headers:            testDataHeaders,
		TestDataRows:       testDataRows,
	}

	return testDataFromTestDataArea
}
