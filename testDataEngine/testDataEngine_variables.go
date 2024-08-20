package testDataEngine

// Error/warning texts for control of Group Name
const (
	GroupNameIsUnique                  string = "Group name is OK!"
	GroupNameIsNotUnique               string = "Group name already exists!"
	GroupNameIsEmpty                   string = "Group name can't be empty"
	GroupNameCanNotStartOrEndWithSpace string = "Group name can not start or end with a space"
	SelectedListIsEmpty                string = "There is nothing to save"
)

// TestDataModel
// The Full TestDataModel
var TestDataModel TestDataModelStruct

// TestDataFromSimpleTestDataAreaStruct
// Struct for holding TestData for one TestDataArea within a TestDataDomain
type TestDataFromSimpleTestDataAreaStruct struct {
	TestDataDomainUuid         string
	TestDataDomainName         string
	TestDataDomainTemplateName string
	TestDataAreaUuid           string
	TestDataAreaName           string
	Headers                    []struct {
		ShouldHeaderActAsFilter bool
		HeaderName              string
		HeaderUiName            string
	}
	TestDataRows           [][]string
	TestDataFileSha256Hash string
}

// TestDataForGroupObjectStruct
// Is used to keep variables per instance of 'testDataForGroupObject'
// Used for keep track of Group-TestData per TestCase or per TestSuite
type TestDataForGroupObjectStruct struct {

	// The selected TestData for the TestCase
	SelectedTestDataGroup            string
	SelectedTestDataPoint            string
	SelectedTestDataPointRowSummary  string
	SelectedTestDataPointRowUniqueId string

	// When all of SelectedTestDataGroup, SelectedTestDataPoint and SelectedTestDataPointRowSummary
	// are populated then TestDataColumnDataNameToValueMap will be created with values
	TestDataColumnDataNameToValueMap map[string]string // map[TestDataColumnDataNameType]TestDataValueType

	// The slices for Groups ans TestDataPoints for a Group
	TestDataPointGroups     []TestDataPointGroupNameType // Define TestDataPointGroups
	TestDataPointsForAGroup []TestDataValueNameType      // Define TestDataPointGroups

	// The map holding all data for Groups and TestDataPoints for a group
	ChosenTestDataPointsPerGroupMap map[TestDataPointGroupNameType]*TestDataPointNameMapType

	// Variable to be used when closing window to inform calling window if the data was updated or not
	ShouldUpdateMainWindow ResponseChannelStruct
}

// TestDataModelStruct
// The structure holding all TestData
type TestDataModelStruct struct {
	TestDataDomainAndAreaNameToUuidMap *map[TestDataDomainOrAreaNameType]TestDataDomainOrAreaUuidType
	TestDataModelMap                   *map[TestDataDomainUuidType]*TestDataDomainModelStruct
}

// TestDataDomainModelStruct
// DataData for one Domain
type TestDataDomainModelStruct struct {
	TestDataDomainUuid TestDataDomainUuidType
	TestDataDomainName TestDataDomainNameType
	TestDataAreasMap   *map[TestDataAreaUuidType]*TestDataAreaStruct
}

// TestDataAreaStruct
// TestData for one Area within one Domain
type TestDataAreaStruct struct {
	TestDataDomainUuid                   TestDataDomainUuidType
	TestDataDomainName                   TestDataDomainNameType
	TestDataAreaUuid                     TestDataAreaUuidType
	TestDataAreaName                     TestDataAreaNameType
	TestDataValuesForRowMap              *map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct
	TestDataValuesForRowNameMap          *map[TestDataValueNameType]*[]TestDataPointRowUuidType
	TestDataValuesForColumnMap           *map[TestDataColumnUuidType]*[]*TestDataPointValueStruct
	TestDataValuesForColumnAndRowUuidMap *map[TestDataColumnAndRowUuidType]*TestDataPointValueStruct
	TestDataColumnsMetaDataMap           *map[TestDataColumnUuidType]*TestDataColumnMetaDataStruct
	UniqueTestDataValuesForColumnMap     *map[TestDataColumnUuidType]*map[TestDataValueType][]TestDataPointRowUuidType
}

// TestDataPointValueStruct
// Holding all information about one TestDataValue
type TestDataPointValueStruct struct {
	TestDataDomainUuid           TestDataDomainUuidType
	TestDataDomainName           TestDataDomainNameType
	TestDataAreaUuid             TestDataAreaUuidType
	TestDataAreaName             TestDataAreaNameType
	TestDataColumnUuid           TestDataColumnUuidType
	TestDataColumnDataName       TestDataColumnDataNameType
	TestDataColumnUIName         TestDataColumnUINameType
	TestDataPointRowUuid         TestDataPointRowUuidType
	TestDataColumnAndRowUuid     TestDataColumnAndRowUuidType
	TestDataValue                TestDataValueType
	TestDataValueNameDescription TestDataValueNameDescriptionType
	TestDataValueName            TestDataValueNameType
}

// TestDataColumnMetaDataStruct
// Holds the information about a column regarding if data should be visible when selecting for data, and some other stuff
type TestDataColumnMetaDataStruct struct {
	TestDataColumnUuid     TestDataColumnUuidType
	TestDataColumnDataName TestDataColumnDataNameType
	TestDataColumnUIName   TestDataColumnUINameType
	//TestDataPointRowsUuid                   []TestDataPointRowUuidType
	ShouldColumnBeUsedForFindingTestData    bool
	ShouldColumnBeUsedWithinTestDataSetName bool
}

// TestDataDomainUuidType
// The UUID of the Domain that owns the TestData
type TestDataDomainUuidType string

// TestDataDomainNameType
// The Name of the Domain that owns the TestData
type TestDataDomainNameType string

// TestDataDomainTemplateNameType
// The Name, of the Domain, used in templates
type TestDataDomainTemplateNameType string

// TestDataAreaUuidType
// The UUID for a specific TestData-area within the Domain
type TestDataAreaUuidType string

// TestDataAreaNameType
// The Name for a specific TestData-area within the Domain
type TestDataAreaNameType string

// TestDataDomainOrAreaUuidType
// The UUID for a specific TestData-domain or TestData-area
type TestDataDomainOrAreaUuidType string

// TestDataDomainOrAreaNameType
// The Name for a specific TestData-domain or TestData-area
type TestDataDomainOrAreaNameType string

// TestDataPointRowUuidType
// Identifies one TestDataPoint/TestDataRow
type TestDataPointRowUuidType string

// TestDataColumnUuidType
// Identifies a Column that goes across TestDataPoints/TestDataRows
type TestDataColumnUuidType string

// TestDataColumnDataNameType
// The DataName used for Placeholder and other data manipulations connected to a specific column in the TestDataSet
type TestDataColumnDataNameType string

// TestDataColumnUINameType
// The Name used in UI connected to a specific column in the TestDataSet
type TestDataColumnUINameType string

// TestDataColumnAndRowUuidType
// The hash ColumnUuid 'concat' RowUuid: SHA256(TestDataColumnUuidType 'concat' TestDataPointRowUuidType)
type TestDataColumnAndRowUuidType string

// TestDataValueType
// The value for specific TestDataPoint-value
type TestDataValueType string

// TestDataValueNameDescriptionType
// The Description of how the name for a TestDataPoint is constructed. Used to show/describe the TestDataPoint in the UI
type TestDataValueNameDescriptionType string

// TestDataValueNameType
// The name for a TestDataPoint. Used to show/describe the TestDataPoint in the UI
type TestDataValueNameType string

// TestDataPointRowValuesSummaryType
// All values for a TestDataRow is concatenated into single value as a summary
type TestDataPointRowValuesSummaryType string

// The type for the map that holds the connection between TestDataÅpintName and all the RowUUids connectoed to it
type TestDataPointNameMapType map[TestDataValueNameType]*[]*DataPointTypeForGroupsStruct

// Types used for data structures for a specific rowValue
/*type testDataPointRowStruct struct {
	testDataDomainUuid           TestDataDomainUuidType
	testDataDomainName           TestDataDomainNameType
	testDataAreaUuid             TestDataAreaUuidType
	testDataAreaName             TestDataAreaNameType
	testDataPointUuid            TestDataPointRowUuidType
	testDataPointName            TestDataValueNameType
	testDataPointNameDescription TestDataValueNameDescriptionType
	testDatapointValue           TestDataValueType
} */

// The type for Group name
type TestDataPointGroupNameType string

// Message sent back when a Group is Created or is Edited
type ResponseChannelStruct struct {
	ShouldBeUpdated        bool
	TestDataPointGroupName TestDataPointGroupNameType
}

type TestDataPointRowUuidStruct struct {
	TestDataPointRowUuid          TestDataPointRowUuidType
	TestDataPointRowValuesSummary TestDataPointRowValuesSummaryType
}

// Structure for the Group
type DataPointTypeForGroupsStruct struct {
	TestDataDomainUuid            TestDataDomainUuidType
	TestDataDomainName            TestDataDomainNameType
	TestDataDomainTemplateName    TestDataDomainTemplateNameType
	TestDataAreaUuid              TestDataAreaUuidType
	TestDataAreaName              TestDataAreaNameType
	TestDataPointName             TestDataValueNameType
	SearchResultDataPointUuidMap  map[TestDataPointRowUuidType]TestDataPointRowUuidStruct
	AvailableTestDataPointUuidMap map[TestDataPointRowUuidType]TestDataPointRowUuidStruct
	SelectedTestDataPointUuidMap  map[TestDataPointRowUuidType]TestDataPointRowUuidStruct
}
