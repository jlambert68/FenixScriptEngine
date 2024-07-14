package testDataEngine

import (
	"sort"
)

// ListTestDataGroups
// List the current TestDataGroups that the User has
func (testDataForGroupObject *TestDataForGroupObjectStruct) ListTestDataGroups() (testDataPointGroupsAsStringSlice []string) {

	// Loop all 'TestDataPointGroups'
	for _, tempTestDataPointGroup := range testDataForGroupObject.TestDataPointGroups {
		testDataPointGroupsAsStringSlice = append(testDataPointGroupsAsStringSlice, string(tempTestDataPointGroup))
	}

	// Sort the Groups
	sort.Strings(testDataPointGroupsAsStringSlice)

	return testDataPointGroupsAsStringSlice

}

// ListTestDataGroupPointsForAGroup
// List the current TestDataGroupPoints for a specific TestDataGroup
func (testDataForGroupObject *TestDataForGroupObjectStruct) ListTestDataGroupPointsForAGroup(testDataGroup string) (testDataPointGroupsAsStringSlice []string) {

	// Extract the map with the TestDataPoints
	var tempTestDataPointNameMap TestDataPointNameMapType
	tempTestDataPointNameMap = *testDataForGroupObject.ChosenTestDataPointsPerGroupMap[TestDataPointGroupNameType(testDataGroup)]

	// Refill the slice with all TestDataPoints
	for testDataPoint, _ := range tempTestDataPointNameMap {
		testDataPointGroupsAsStringSlice = append(testDataPointGroupsAsStringSlice, string(testDataPoint))

	}

	// Sort the GroupPoints
	sort.Strings(testDataPointGroupsAsStringSlice)

	return testDataPointGroupsAsStringSlice

}

// ListTestDataRowsForAGroupPoint
// List the current TestDataRow for a specific TestDataGroupPoint
func (testDataForGroupObject *TestDataForGroupObjectStruct) ListTestDataRowsForAGroupPoint(testDataGroup string, testDataGroupPoint string) (
	testDataGroupPointRowsSummaryValueAsStringSlice []string) {

	//fixa denna

	// Extract the map with the TestDataPoints
	var tempTestDataPointNameMap TestDataPointNameMapType
	var dataPointRowsSlicePtr *[]*DataPointTypeForGroupsStruct
	var dataPointRowsSlice []*DataPointTypeForGroupsStruct

	// Extract DataPoints from for Group
	tempTestDataPointNameMap = *testDataForGroupObject.ChosenTestDataPointsPerGroupMap[TestDataPointGroupNameType(testDataGroup)]

	// Extract Rows for DataPoint
	dataPointRowsSlicePtr = tempTestDataPointNameMap[TestDataValueNameType(testDataGroupPoint)]
	dataPointRowsSlice = *dataPointRowsSlicePtr

	// Refill the slice with all TestDataPoints
	for _, testDataPointRowUuiObject := range dataPointRowsSlice[0].SelectedTestDataPointUuidMap {
		testDataGroupPointRowsSummaryValueAsStringSlice = append(testDataGroupPointRowsSummaryValueAsStringSlice,
			string(testDataPointRowUuiObject.TestDataPointRowValuesSummary))

	}

	// Sort the GroupPoints
	sort.Strings(testDataGroupPointRowsSummaryValueAsStringSlice)

	return testDataGroupPointRowsSummaryValueAsStringSlice

}

// GetTestDataPointValuesMapBasedOnGroupPointNameAndSummaryValue
// Generate a map with 'TestDataColumnDataName' as key and 'TestDataValue' as value
func (testDataForGroupObject *TestDataForGroupObjectStruct) GetTestDataPointValuesMapBasedOnGroupPointNameAndSummaryValue(
	testDataGroup string,
	testDataPointName string,
	testDataPointRowSummaryValue string) (
	testDataColumnDataNameMap map[string]string) { // map[TestDataColumnDataNameType]TestDataValueType

	// Initiate response-map
	testDataColumnDataNameMap = make(map[string]string)

	if testDataPointName == "" || testDataPointRowSummaryValue == "" {
		return testDataColumnDataNameMap
	}

	// Get 'TestDataPointRowUuid' base on 'TestDataPointRowSummaryValue'
	var testDataPointRowUuid string

	// Extract DataPoints from for Group
	tempTestDataPointNameMap := *testDataForGroupObject.ChosenTestDataPointsPerGroupMap[TestDataPointGroupNameType(testDataGroup)]

	// Extract Rows for DataPoint
	dataPointRowsSlicePtr := tempTestDataPointNameMap[TestDataValueNameType(testDataPointName)]
	dataPointRowsSlice := *dataPointRowsSlicePtr

	// Refill the slice with all TestDataPoints
	for _, testDataPointRowUuiObject := range dataPointRowsSlice[0].SelectedTestDataPointUuidMap {

		if string(testDataPointRowUuiObject.TestDataPointRowValuesSummary) == testDataPointRowSummaryValue {
			testDataPointRowUuid = string(testDataPointRowUuiObject.TestDataPointRowUuid)
			break
		}
	}

	// Create the data table for all matching 'testDataPointRowUuid'
	var tableData [][]string
	tableData = BuildPopUpTableDataFromTestDataPointName(testDataPointName, &TestDataModel)

	var headerSlice []string

	// Loop alla rows
	for rowIndex, rowData := range tableData {

		// Loop alla values for row
		for columnIndex, columnValue := range rowData {

			// Create a header slice
			if rowIndex == 0 {
				// Header row
				headerSlice = append(headerSlice, columnValue)

			} else {
				// Only process if this is the correct row
				if rowData[len(rowData)-1] == testDataPointRowUuid {

					testDataColumnDataNameMap[headerSlice[columnIndex]] = columnValue

				}
			}

		}

	}

	return testDataColumnDataNameMap

}

// GetTestDataModelPtr
// Returns a pointer to the full TestDataModel
func GetTestDataModelPtr() *TestDataModelStruct {
	return &TestDataModel
}
