
local assert = require "luassert"
local mymodule = require('src/Fenix_TodayDateShift')
local date = require('date')

local tests = {}

-- Check that Todays date is produced
function tests.is_today_no_parameter()
    local today = date()
    local expectedDateResponse = today:fmt("%Y-%m-%d")
    local inputArray = {"Fenix_TodayShiftDay", {}, {}, {true, 0}}
    local response = Fenix_TodayShiftDay(inputArray)

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)


end

-- Check that Todays date is produced
function tests.is_today()
    local today = date()
    local expectedDateResponse = today:fmt("%Y-%m-%d")
    local inputArray = {"Fenix_TodayShiftDay", {}, {0}, {true, 0}}
    local response = Fenix_TodayShiftDay(inputArray)

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- Check that Yesterdays date is produced
function tests.is_yesterday()
    local today = date()
    local shoftedDate = today:adddays(-1)
    local expectedDateResponse = today:fmt("%Y-%m-%d")
    local inputArray = {"Fenix_TodayShiftDay", {}, {-1}, {true, 0}}
    local response = Fenix_TodayShiftDay(inputArray)

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- Check that Tomorrows date is produced
function tests.is_tomorrow()
    local today = date()
    local shoftedDate = today:adddays(1)
    local expectedDateResponse = today:fmt("%Y-%m-%d")
    local inputArray = {"Fenix_TodayShiftDay", {}, {1}, {true, 0}}
    local response = Fenix_TodayShiftDay(inputArray)

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- Check that Array Index is not supported
function tests.array_value_not_allowed_v1()

    local expectedDateResponse = ""
    local inputArray = {"Fenix_TodayShiftDay", {1}, {0}, {true, 0}}
    local response = Fenix_TodayShiftDay(inputArray)

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(false, response.success)

end

-- Check that Array Index is not supported
function tests.array_value_not_allowed_v2()

    local expectedDateResponse = ""
    local inputArray = {"Fenix_TodayShiftDay", {1,2}, {0}, {true, 0}}
    local response = Fenix_TodayShiftDay(inputArray)

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(false, response.success)

end

-- Check that there are not more then one parameter
function tests.array_value_not_allowed_v3()

    local expectedDateResponse = ""
    local inputArray = {"Fenix_TodayShiftDay", {}, {1, 2}, {true, 0}}
    local response = Fenix_TodayShiftDay(inputArray)

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(false, response.success)

end

-- Check that there are exact four parameters in input Table
function tests.exact_four_parameters_in_input_table()

    local expectedDateResponse = ""
    local inputArray = {"Fenix_TodayShiftDay", {}, {1, 2}}
    local response = Fenix_TodayShiftDay(inputArray)

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(false, response.success)

end

--tests.test_my_func()

return tests






