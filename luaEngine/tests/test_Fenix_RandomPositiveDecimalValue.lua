
local assert = require "luassert"
local mymodule = require('src/Fenix_RandomPositiveDecimalValue')

local tests = {}

-- OK - {"Fenix_RandomPositiveDecimalValue", {},{2, 3}, {0}}
function tests.ok_array__parameters_2_3_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {},{2, 3}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "84.394"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- OK - {"Fenix_RandomPositiveDecimalValue", {1},{2, 3}, {0}}
function tests.ok_array_1_parameters_2_3_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {1},{2, 3}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "84.394"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- OK - {"Fenix_RandomPositiveDecimalValue", {},{1, 2}, {0}}
function tests.ok_array__parameters_1_2_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {},{1, 2}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "8.39"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- OK - {"Fenix_RandomPositiveDecimalValue", {2},{1, 2}, {0}}
function tests.ok_array_2_parameters_1_2_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {2},{1, 2}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "7.80"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end


-- OK - {"Fenix_RandomPositiveDecimalValue", {},{1, 1}, {0}}
function tests.ok_array__parameters_1_1_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {},{1, 1}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "8.3"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- OK - {"Fenix_RandomPositiveDecimalValue", {},{1, 1}, 1}
function tests.ok_array__parameters_1_1_entropi_1()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {},{1, 1}, {true, 1}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "7.8"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- OK - {"Fenix_RandomPositiveDecimalValue", {1},{1, {0}}, {0}}
function tests.ok_array_1_parameters_1_0_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {1},{1, 0}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "8"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- OK - {"Fenix_RandomPositiveDecimalValue", {1},{0, {0}}, {0}}
function tests.ok_array_1_parameters_0_0_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {1},{0, 0}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "0"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- OK - {"Fenix_RandomPositiveDecimalValue", {1},{0, 0, 2, 3}, {0}}
function tests.ok_array_1_parameters_0_0_2_3_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {1},{0, 0, 2, 3}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "00"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- OK - {"Fenix_RandomPositiveDecimalValue", {0},{6, 6}, {0}}
function tests.ok_array__parameters_6_6_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {0},{6, 6}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "840187.394382"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- OK - {"Fenix_RandomPositiveDecimalValue", {0},{6, 10}, {0}}
function tests.ok_array__parameters_6_10_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {0},{6, 10}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "840187.3943829300"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- OK - {"Fenix_RandomPositiveDecimalValue", {0},{0, 0, 2, 3}, {0}}
function tests.ok_array__parameters_0_0_2_3_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {0},{0, 0, 2, 3}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "00"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- OK - {"Fenix_RandomPositiveDecimalValue", {0},{0, 2, 3, 4}, {0}}
function tests.ok_array__parameters_0_2_3_4_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {0},{0, 2, 3, 4}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "000.3900"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- OK - {"Fenix_RandomPositiveDecimalValue", {0},{2, 2, 3, 4}, {0}}
function tests.ok_array__parameters_2_2_3_4_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {0},{2, 2, 3, 4}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = "084.3900"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(true, response.success)
    assert.is_equal("", response.errorMessage)

end

-- ERROR - {"Fenix_RandomPositiveDecimalValue", {1},{0}, {0}}
function tests.error_array_1_parameters_1__entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {1},{1}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = ""
    local expectedErrorMessage = "Error - there must be exact 2 or 4 function parameter. '[1]'"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(false, response.success)
    assert.is_equal(expectedErrorMessage, response.errorMessage)

end

-- ERROR - {"Fenix_RandomPositiveDecimalValue", {1},{}, {0}}
function tests.error_array_1_parameters___entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {1},{}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = ""
    local expectedErrorMessage = "Error - there must be exact 2 or 4 function parameter but it is empty."

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(false, response.success)
    assert.is_equal(expectedErrorMessage, response.errorMessage)

end

-- ERROR - {"Fenix_RandomPositiveDecimalValue", {1},{1, 2, 3}, {0}}
function tests.error_array_1_parameters_1_2_3_entropi_0()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {1},{1, 2, 3}, {true, 0}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = ""
    local expectedErrorMessage = "Error - there must be exact 2 or 4 function parameter. '[1,2,3]'"

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(false, response.success)
    assert.is_equal(expectedErrorMessage, response.errorMessage)

end

-- ERROR - {"Fenix_RandomPositiveDecimalValue", {1},{2, 3}}
function tests.error_array_1_parameters_2_3_entropi_()


    local inputArray =  {"Fenix_RandomPositiveDecimalValue", {1},{2, 3}}
    local response = Fenix_RandomPositiveDecimalValue(inputArray)
    local expectedDateResponse = ""
    local expectedErrorMessage = "Error - there should be exactly four rows in InputTable."

    assert.is_equal(expectedDateResponse, response.value)
    assert.is_equal(false, response.success)
    assert.is_equal(expectedErrorMessage, response.errorMessage)

end

return tests

--[[



--]]