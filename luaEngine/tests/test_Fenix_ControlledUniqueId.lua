local assert = require('luassert')
local mymodule = require('src/Fenix_ControlledUniqueId')

local tests = {}

-- Mock input table structure
local function createInputTable(text, seed)
    return {"Fenix_ControlledUniqueId", {0}, {text}, {seed}}
end

-- Test for date in YYYY-MM-DD format
function tests.test_dateFormatYYYYMMDD()
    local inputString = "Today's date is %YYYY-MM-DD%"
    local inputTable = createInputTable(inputString, 0)
    local result = Fenix_ControlledUniqueId(inputTable)

    local expectedDate = os.date("%Y-%m-%d")
    local expectedString = "Today's date is " .. expectedDate
    assert.is_equal(result, expectedString)
end

-- Test for time in hh:mm:ss format
function tests.test_timeFormathhmmss()
    local inputString = "Current time is %hh:mm:ss%"
    local inputTable = createInputTable(inputString, 0)
    local result = Fenix_ControlledUniqueId(inputTable)

    local expectedTime = os.date("%H:%M:%S")
    local expectedString = "Current time is " .. expectedTime
    assert.is_equal(result, expectedString)
end

-- Test for random number of length n
function tests.test_randomNumberLength()
    local inputString = "Random number: %nnnnn%"
    local inputTable = createInputTable(inputString, 12345) -- Fixed seed for predictability
    local result = Fenix_ControlledUniqueId(inputTable)

    assert.is_equal(#result:match("%d+"), 5)  -- Check if the length of the number is 5
end

function tests.test_randomSmallLetterLength()
    local inputString = "Random small letters: %a(5, 12345)%"
    local inputTable = createInputTable(inputString, 12345)  -- Fixed seed
    local result = Fenix_ControlledUniqueId(inputTable)

    local match = result:match("Random small letters: (%a+)")
    assert.is_equal(#match, 5)  -- Check if the length of the string is 5
    assert.matches("^[a-z]+$", match)  -- Check if all characters are small letters
end


function tests.test_randomCapitalLetterLength()
    local inputString = "Random capital letters: %A(5, 12345)%"
    local inputTable = createInputTable(inputString, 12345)  -- Fixed seed
    local result = Fenix_ControlledUniqueId(inputTable)

    local match = result:match("Random capital letters: (%a+)")
    assert.is_equal(#match, 5)  -- Check if the length of the string is 5
    assert.matches("^[A-Z]+$", match)  -- Check if all characters are capital letters
end



-- More tests can be added for each format type...

return tests

-- Running the tests
-- os.exit(lu.LuaUnit.run())
