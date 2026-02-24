---------------------------------------------------------------------------------------
-- Module for replacing a Fenix Inception 'Placeholder'
--
-- Version 1.0

-- Placeholder usage
-- 'Fenix.RandomPositiveDecimalValue(IntegerSize, FractionSize)'
-- 'Fenix.RandomPositiveDecimalValue(IntegerSize, FractionSize, IntergerSpace, FractionSpace)'
-- 'Fenix.RandomPositiveDecimalValue[ArraysIndex](IntegerSize, FractionSize)'
-- 'Fenix.RandomPositiveDecimalValue[ArraysIndex](IntegerSize, FractionSize)(UseTestCaseExecutionUuidentropy)'
-- 'Fenix.RandomPositiveDecimalValue[ArraysIndex](IntegerSize, FractionSize)(UseTestCaseExecutionUuidentropy, ExtraentropyNumber)'
-- 'Fenix.RandomPositiveDecimalValue['Integer']('Integer', 'Integer')('Boolean', 'Integer')'
--
-- Resesponse is a lua table;  {value, success, errorMessage} with the following types {'strings', 'boolean', 'string'}
--
-- Usage exmaples
-- 'Fenix.RandomPositiveDecimalValue(2, 3)' is same as 'Fenix.RandomPositiveDecimalValue[1](2, 3)' which is the same as 'Fenix.RandomPositiveDecimalValue[1](2, 3)(true)'
-- They all could produce i.e. '35.693' and with same input the placeholder will allways have the same output.
-- UseTestCaseExecutionUuidentropy(true/false) is based on the 'TestCaseExecutionUuid' and within one TestCaseExecution values with the same input will have the same output
--
-- 'Fenix.RandomPositiveDecimalValue[2](2, 3)(true)' clould have the output of '27.568'
-- 'Fenix.RandomPositiveDecimalValue[2](2, 3)(false)' will allways produce the same output, independently of 'TestCaseExecutionUuid'
-- 'Fenix.RandomPositiveDecimalValue[2](2, 3)(true, 1)' will add extra entropy to seed, by adding 1 to the value based on 'TestCaseExecutionUuid'.

-- 'IntergerSpace' and 'FractionSpace' define the spaces for 'Intergerpart' and 'FractionPart'. Zeros will be added before 'Intergerpart' and after 'FractionPart'
-- If 'IntergerSpace' is less than 'IntegerSize' then it will be ignored
-- If 'FractionSpace' is less than 'FractionSize' then it will be ignored
-- 'Fenix.RandomPositiveDecimalValue(1, 2, 3, 4)' will produce "004.5700"
-- 'Fenix.RandomPositiveDecimalValue(3, 2, 2, 1)' will produce "344.54"
-- 'Fenix.RandomPositiveDecimalValue(0, 0, 2, 2)' will produce "00.00"
-- 
-- Examples of different parameter values 'IntegerSize' and 'FractionSize' and a possible output
-- 'Fenix.RandomPositiveDecimalValue(1, 2)' = "5.48"
-- 'Fenix.RandomPositiveDecimalValue(2, 2)' = "28.87"
-- 'Fenix.RandomPositiveDecimalValue(0, 1)' = "0.37"
-- 'Fenix.RandomPositiveDecimalValue(0, 0)' = "0"
-- 'Fenix.RandomPositiveDecimalValue(0, 0, 2, 3)' = "00.000"
-- 'Fenix.RandomPositiveDecimalValue(0, 0)' = "0.0"
-- 'Fenix.RandomPositiveDecimalValue(3, 0)' = "293"
-- 'Fenix.RandomPositiveDecimalValue(3, 0, 3, 2)' = "293.00"



local function tableToString(tbl, sep)

    -- More then one parameter in table
    if #tbl > 1 then
        sep = sep or ", "
        local result = "["
        for _, v in ipairs(tbl) do
            result = result .. tostring(v) .. sep
        end

        result = result:sub(1, -#sep - 1)

        result = result .. "]"

        return result
    end

    -- Only one parameter in table
    local result = tostring(tbl)

    return result
end

-- ***********************************************************************************
-- round
--
-- Function to round a decimal to a specific number of decimal places

local function round(x, places)

    local shift = 10 ^ places

    return math.floor(x * shift + 0.5) / shift
end

-- ***********************************************************************************



-- ***********************************************************************************
-- formatDecimal
--
-- Function to format a decimal with a specific number of decimals
local function formatDecimal(number, numberOfDecimals)

    -- Convert the number to a string
    local str = tostring(number)

    -- Find the position of the decimal point
    local dotIndex = nil
    for i = 1, #str do
        if str:sub(i, i) == '.' then
            dotIndex = i
            break
        end
    end


    -- Add a decimal point if it doesn't exist
    if numberOfDecimals > 0 and dotIndex == nil then
        dotIndex = #str + 1
        str = str .. "."
    end

        -- Handle case when there should be decimals 
    if numberOfDecimals > 0 then

        -- Calculate the number of decimal places currently in the string
        local currentDecimals = dotIndex and #str - dotIndex or 0

        -- Add zeros to reach the desired number of decimal places
        while currentDecimals < numberOfDecimals do
            str = str .. "0"
            currentDecimals = currentDecimals + 1
        end

        return str
    end

    -- No decimals so remove decimal point and any following zero, i.e. 13.0
    if dotIndex ~= nil then
    local integerWithOutDecimals = string.sub(str, 1, dotIndex -1)

    return integerWithOutDecimals

    else
        return str

    end

end

-- ***********************************************************************************

-- ***********************************************************************************
-- padValueWithZeros
--
-- Function to pad the integer part and the fraction part with correct number of zeros

local function padValueWithZeros(valueAsString, integerSpace, fractionSpace)

    -- Split the value into integer and fractional parts
    local integerPart, fractionPart = valueAsString:match("^(%d+)%.(%d+)$")
    local noFractions = false
    
    -- Check if value only has intgerpart
    if not integerPart or not fractionPart then
        integerPart = valueAsString
        fractionPart = ""
        noFractions = true
    end

    -- Pad the integer part with zeros if needed
    if #integerPart < integerSpace then
        integerPart = string.rep("0", integerSpace - #integerPart) .. integerPart
    end

    -- Pad the fractional part with zeros if needed
    if #fractionPart < fractionSpace and noFractions == false then
        fractionPart = fractionPart .. string.rep("0", fractionSpace - #fractionPart)
    end

    -- Combine the padded parts if there are two parts
    local zeroPaddedValue = ""

    if noFractions == false then
        zeroPaddedValue = integerPart .. "." .. fractionPart

    else
        zeroPaddedValue = integerPart

    end

    return zeroPaddedValue
end

-- ***********************************************************************************



-- ***********************************************************************************
-- randomize

-- Function to generate random numbers


local function randomize(arrayIndex, maxIntegerPartSize, numberOfDecimals, baseentropyToUse)

    math.randomseed(baseentropyToUse+ arrayIndex)

    -- Generate Integer part of random number
    local randomIntegerPart = math.random()
    local integerPart = math.floor(10 ^ maxIntegerPartSize * randomIntegerPart)

    -- Generate Decimal part of random number
    local randomDecimalPart = math.random()
    local decimalPart = 0

    if numberOfDecimals > 0 then
        decimalPart = math.floor(10 ^ numberOfDecimals * randomDecimalPart)
    end

    -- Combine Integer and decimal part into one random number
    local randomNumber = integerPart + 10 ^ (-1 * numberOfDecimals) * decimalPart

    randomNumber = round(randomNumber, numberOfDecimals)

    return randomNumber
end


-- ***********************************************************************************


 -- ***********************************************************************************
-- stringToInteger
--
-- Converts a string to an integer in Lua
local function stringToInteger(str)
    local num = tonumber(str)
    if num and num == math.floor(num) then
        return num
    else
        return nil, "The provided string is not an integer."
    end

    --local num = math.tointeger(str)
    --if num then
    --    return num, ""
    --else
    --    return nil, "The provided string is not an integer."
    --end

end

 -- ***********************************************************************************

 -- ***********************************************************************************
-- stringToBoolean
--
-- Converts a string to a boolean in Lua

local function stringToBoolean(inputString)

    -- Secure that input is a string
    if type(inputString) ==  "string" then

        inputString = inputString:lower()  -- Convert the string to lower case to make the function case-insensitive
        if inputString == "true" then
            return true
        elseif inputString == "false" then
            return false
        else
            return nil, "Invalid input for boolean conversion: "
        end
    else
        return nil, "Invalid input for boolean conversion: "
    end
end

-- ***********************************************************************************



 -- ***********************************************************************************
-- Fenix_RandomDecimalValue_ArrayValue // Fenix.RandomDecimalValue[n](maxIntegerPartSize, numberOfDecimals)
--
-- Function to generate random value with a specif max number of integer and speciic number of decimals
-- inputArray := [arrayPosition, maxIntegerPartSize, numberOfDecimals, testCaseUuidentropy]


local function Fenix_RandomDecimalValue_ArrayValue(inputArray)
    local arrayPositionToUse = inputArray[1]
    local maxIntegerPartSize = inputArray[2][1]
    local numberOfDecimals = inputArray[2][2]

    local entropyToUse = inputArray[3]

    local tempValueAsDecimal = randomize(arrayPositionToUse, maxIntegerPartSize, numberOfDecimals, entropyToUse)

    local valueIsBaseFormated =  formatDecimal(tempValueAsDecimal, numberOfDecimals)

    -- No padding should be done
    if #inputArray[2] == 2 then
        return valueIsBaseFormated
    end

    -- extract padding sizes
    local integerSpace = inputArray[2][3]
    local fractionSpace = inputArray[2][4]

    local zeroPaddedValue = padValueWithZeros(valueIsBaseFormated, integerSpace, fractionSpace)

    return zeroPaddedValue

end

-- ***********************************************************************************


 -- ***********************************************************************************
-- Fenix_RandomDecimalValue_Sum_ArrayValue // Fenix.RandomDecimalValue[n](maxIntegerPartSize, numberOfDecimals)
--
-- For each arrayPositionValuem the Function  generate random value with a specif max number of integer and speciic number of decimals
-- Then the function sums the values
-- inputArray := [arrayPosition, maxIntegerPartSize, numberOfDecimals, testCaseUuidentropy]


local function Fenix_RandomDecimalValue_Sum_ArrayValue(inputArray)
    local arrayPositionArray = inputArray[1]
    local maxIntegerPartSize = inputArray[2][1]
    local numberOfDecimals = inputArray[2][2]

    local entropyToUse = inputArray[3]

    local sumOfvalues = 0

    local tempValueAsDecimal

    local valueIsBaseFormated
    local arrayPositionToUse

    -- Loop every arrayPositionArray
    for tableIndex, arrayPositionValue in ipairs(arrayPositionArray) do

        -- Extract value from array
        arrayPositionToUse = math.abs(arrayPositionValue)

        -- Generate value
        tempValueAsDecimal = randomize(arrayPositionToUse, maxIntegerPartSize, numberOfDecimals, entropyToUse)

        if arrayPositionValue >= 0 then
            sumOfvalues = sumOfvalues + tempValueAsDecimal
        else
            sumOfvalues = sumOfvalues - tempValueAsDecimal
        end

    end

    -- Format to correct number of decimals
    valueIsBaseFormated =  formatDecimal(sumOfvalues, numberOfDecimals)


    -- No padding should be done
    if #inputArray[2] == 2 then
        return valueIsBaseFormated
    end

    -- extract padding sizes
    local integerSpace = inputArray[2][3]
    local fractionSpace = inputArray[2][4]

    local zeroPaddedValue = padValueWithZeros(valueIsBaseFormated, integerSpace, fractionSpace)

    return zeroPaddedValue

end

-- ***********************************************************************************



-- ***********************************************************************************
-- Fenix_RandomPositiveDecimalValue
--
-- Function to generate random value with a specif max number of integer and speciic number of decimals
-- Always use array value 1, first array position from user perspective
--
-- inputArray := [[arrayindex], [maxIntegerPartSize, numberOfDecimals], testCaseUuidentropy]

function Fenix_RandomPositiveDecimalValue(inputTable)

    local responseTable = {
        success = true,
        value = "",
        errorMessage = ""
    }

    -- There must be 4 rows in the InputTable
    if #inputTable ~= 4 then

        local error_message = "Error - there should be exactly four rows in InputTable."

        responseTable.success = false
        responseTable.errorMessage = error_message

        return responseTable

    end


    -- Extract ArraysIndexArray
    local arraysIndexTable = inputTable[2]
    local arraysIndexToUse

    -- Secure that ArraysIndexArray is not emtpy or only have one value
    if #arraysIndexTable > 1 then
        -- Have more then 1 value

        -- Convert array to string
        local tableAsString = tableToString (arraysIndexTable, ",")
        local error_message = "Error - array index array can only have a maximum of one value. '" .. tableAsString .. "'"

        responseTable.success = false
        responseTable.errorMessage = error_message

        return responseTable

    elseif  #arraysIndexTable == 0 then
        -- zero array index, so use first index position

        arraysIndexToUse = 1

    else
        -- Array has one value so use that
        arraysIndexToUse = arraysIndexTable[1]

    end

    -- Secure that 'arraysIndexToUse' is an integer
    if  type(arraysIndexToUse) ~=  "number" then

       -- If not an number then try to convert into an Integer
        local tempInteger, err = stringToInteger(arraysIndexToUse)
        if tempInteger then
            arraysIndexToUse = tempInteger
        else

            local tableAsString = tableToString (arraysIndexTable, ",")
            local error_message = "Error - Array index must be of type 'Integer', '" .. tableAsString .. "'"

            responseTable.success = false
            responseTable.errorMessage = error_message

            return responseTable
        end
    end

    -- Extract FunctionArgumentsArray
    local functionArgumentsTable = inputTable[3]

    -- Handle if function arguments is not 2 arguments or 4 arguments
    if #functionArgumentsTable ~=  2 and #functionArgumentsTable ~=  4 then

        -- More then 2 arguments
        if #functionArgumentsTable > 2 then

            local tableAsString = tableToString(functionArgumentsTable, ",")
            local error_message = "Error - there must be exact 2 or 4 function parameter. '" .. tableAsString .. "'"

            responseTable.success = false
            responseTable.errorMessage = error_message

            return responseTable

            -- Exact one argument
        elseif #functionArgumentsTable == 1 then

                local result = "[" .. tostring(functionArgumentsTable[1]) .. "]"

                local error_message = "Error - there must be exact 2 or 4 function parameter. '" .. result .. "'"

                responseTable.success = false
                responseTable.errorMessage = error_message

                return responseTable

            -- Zero values    
        else
                local error_message = "Error - there must be exact 2 or 4 function parameter but it is empty."

                responseTable.success = false
                responseTable.errorMessage = error_message

                return responseTable

        end
    end

    -- verify that each function parameter is a number
    for tableIndex, v in ipairs(functionArgumentsTable) do

        -- Must be an integer 
        if type(v) ~=  "number" then

          -- If not an number then try to convert into an Integer
            local tempInteger, err = stringToInteger(v)
            if tempInteger then
                functionArgumentsTable[tableIndex] = tempInteger
            else

                local tableAsString = tableToString (functionArgumentsTable, ",")
                local error_message = "Error - functions parameters must be of type Integer. '" .. tableAsString .. "'"

                responseTable.success = false
                responseTable.errorMessage = error_message

                return responseTable

          end

        end
    end

    -- Extract entropy
    local entropyTable = inputTable[4]

    -- Verify that content in entropy is of type 'Table'
    if type(entropyTable) ~= "table" then

        local error_message = "Error - entropy is not of type 'Table', but is of type '" .. type(entropyTable)  .. "'."

        responseTable.success = false
        responseTable.errorMessage = error_message

        return responseTable
    end

    -- verify that first parameter is true|false and second paramter in entropy table is a number 
    local entropyValue = 0
    
     -- verify that first parameter is true|false, ie a boolean
    if type(entropyTable[1]) ~= "boolean" then

        -- Try to convert into a boolean
        local stringToBooleanResponse
        stringToBooleanResponse = stringToBoolean(entropyTable[1])

        if type(stringToBooleanResponse) == "boolean" then
            entropyTable[1] = stringToBooleanResponse

        else
            local tableAsString = tableToString (entropyTable, ",")
            local error_message = "Error - entropy parameter no. 1 must be of type 'Boolean' but is'" .. type(entropyTable[1]) .. "'', " .. tableAsString .. "'"

            responseTable.success = false
            responseTable.errorMessage = error_message

            return responseTable

        end
    end
    
    -- verify that second parameter is an Integer
    if type(entropyTable[2]) ~=  "number" then

        -- If not an number then try to convert into an Integer
        local tempInteger, err = stringToInteger(entropyTable[2])
        if tempInteger then
            entropyTable[2] = tempInteger
        else

            local tableAsString = tableToString (entropyTable, ",")
            local error_message = "Error - entropy parameters no. 2 must be of type 'Integer', '" .. tableAsString .. "'"

            responseTable.success = false
            responseTable.errorMessage = error_message

            return responseTable
            end
    end

    entropyValue = entropyTable[2]

    -- Make new Array to be send to the function that does stuff
    local inputTableForProcessing = {arraysIndexToUse, functionArgumentsTable, entropyValue}

    -- Call and process Random Decimal Value
    local response = Fenix_RandomDecimalValue_ArrayValue(inputTableForProcessing)

    responseTable.success = true
    responseTable.errorMessage = ""
    responseTable.value = response

    return responseTable

end


-- ***********************************************************************************
-- Fenix_RandomPositiveDecimalValue_Sum
--
-- Function to generate random value with a specif max number of integer and speciic number of decimals
-- Always use array value 1, first array position from user perspective
--
-- inputArray := [[integer, integer, ...], [maxIntegerPartSize, numberOfDecimals], testCaseUuidentropy]

function Fenix_RandomPositiveDecimalValue_Sum(inputTable)

    local responseTable = {
        success = true,
        value = "",
        errorMessage = ""
    }

    -- There must be 4 rows in the InputTable
    if #inputTable ~= 4 then

        local error_message = "Error - there should be exactly four rows in InputTable."

        responseTable.success = false
        responseTable.errorMessage = error_message

        return responseTable

    end


    -- Extract ArraysIndexArray
    local arraysIndexTable = inputTable[2]

    -- Secure that ArraysIndexArray is not emtpy
    if #arraysIndexTable == 0 then
        -- zero array index, so use first index position

        arraysIndexTable[1] = 1

    end

    -- Secure that 'arraysIndexTable' only has positive or negative integers
    local newArraysIndexTable
    for tableIndex, v in ipairs(arraysIndexTable) do
        if  type(v) ~=  "number" then

           -- If not an number then try to convert into an Integer
            local tempInteger, err = stringToInteger(v)
            if tempInteger then
                arraysIndexTable[tableIndex] = tempInteger
            else

                local tableAsString = tableToString (arraysIndexTable, ",")
                local error_message = "Error - indexArray can only be of type 'Integer', '" .. tableAsString .. "'"

                responseTable.success = false
                responseTable.errorMessage = error_message

                return responseTable
            end

        end
    end

    -- Extract FunctionArgumentsArray
    local functionArgumentsTable = inputTable[3]

    -- Handle if function arguments is not 2 arguments or 4 arguments
    if #functionArgumentsTable ~=  2 and #functionArgumentsTable ~=  4 then

        -- More then 2 arguments
        if #functionArgumentsTable > 2 then

            local tableAsString = tableToString(functionArgumentsTable, ",")
            local error_message = "Error - there must be exact 2 or 4 function parameter. '" .. tableAsString .. "'"

            responseTable.success = false
            responseTable.errorMessage = error_message

            return responseTable

            -- Exact one argument
        elseif #functionArgumentsTable == 1 then

                local result = "[" .. tostring(functionArgumentsTable[1]) .. "]"

                local error_message = "Error - there must be exact 2 or 4 function parameter. '" .. result .. "'"

                responseTable.success = false
                responseTable.errorMessage = error_message

                return responseTable

            -- Zero values
        else
                local error_message = "Error - there must be exact 2 or 4 function parameter but it is empty."

                responseTable.success = false
                responseTable.errorMessage = error_message

                return responseTable

        end
    end

    -- verify that each function parameter is a number
    for tableIndex, v in ipairs(functionArgumentsTable) do

        -- Must be an integer
        if type(v) ~=  "number" then

          -- If not an number then try to convert into an Integer
            local tempInteger, err = stringToInteger(v)
            if tempInteger then
                functionArgumentsTable[tableIndex] = tempInteger
            else

                local tableAsString = tableToString (functionArgumentsTable, ",")
                local error_message = "Error - functions parameters must be of type Integer. '" .. tableAsString .. "'"

                responseTable.success = false
                responseTable.errorMessage = error_message

                return responseTable

          end

        end
    end

    -- Extract entropy
    local entropyTable = inputTable[4]

    -- Verify that content in entropy is of type 'Table'
    if type(entropyTable) ~= "table" then

        local error_message = "Error - entropy is not of type 'Table', but is of type '" .. type(entropyTable)  .. "'."

        responseTable.success = false
        responseTable.errorMessage = error_message

        return responseTable
    end

    -- verify that first parameter is true|false and second paramter in entropy table is a number
    local entropyValue = 0

     -- verify that first parameter is true|false, ie a boolean
    if type(entropyTable[1]) ~= "boolean" then

        -- Try to convert into a boolean
        local stringToBooleanResponse
        stringToBooleanResponse = stringToBoolean(entropyTable[1])

        if type(stringToBooleanResponse) == "boolean" then
            entropyTable[1] = stringToBooleanResponse

        else
            local tableAsString = tableToString (entropyTable, ",")
            local error_message = "Error - entropy parameter no. 1 must be of type 'Boolean' but is'" .. type(entropyTable[1]) .. "'', " .. tableAsString .. "'"

            responseTable.success = false
            responseTable.errorMessage = error_message

            return responseTable

        end
    end

    -- verify that second parameter is an Integer
    if type(entropyTable[2]) ~=  "number" then

        -- If not an number then try to convert into an Integer
        local tempInteger, err = stringToInteger(entropyTable[2])
        if tempInteger then
            entropyTable[2] = tempInteger
        else

            local tableAsString = tableToString (entropyTable, ",")
            local error_message = "Error - entropy parameters no. 2 must be of type 'Integer', '" .. tableAsString .. "'"

            responseTable.success = false
            responseTable.errorMessage = error_message

            return responseTable
            end
    end

    entropyValue = entropyTable[2]

    -- Make new Array to be send to the function that does stuff
    local inputTableForProcessing = {arraysIndexTable, functionArgumentsTable, entropyValue}
    


    -- Call and process Random Decimal Value
    local response = Fenix_RandomDecimalValue_Sum_ArrayValue(inputTableForProcessing)

    responseTable.success = true
    responseTable.errorMessage = ""
    responseTable.value = response

    return responseTable

end





local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{"2", "3"}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{'2', '3'}, {'true', '0'}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{"2", "3"}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{'2', '3'}, {'true', '0'}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {3},{"2", "3"}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {3},{'2', '3'}, {'true', '0'}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {2},{"2", "3"}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {2},{'2', '3'}, {'true', '0'}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue_Sum", {1},{2, 3}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue_Sum(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{2, 3},  {'true', '0'}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue_Sum", {'-1', '2'},{2, 3}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue_Sum(inputArray)
print("{'Fenix_RandomPositiveDecimalValue_Sum', {'1', '2'},{2, 3},  {'true', '0'}}")
print("Fenix_RandomPositiveDecimalValue_Sum: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue_Sum", {1, -2},{2, 3}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue_Sum(inputArray)
print("{'Fenix_RandomPositiveDecimalValue_Sum', {1, -2},{2, 3},  {'true', '0'}}")
print("Fenix_RandomPositiveDecimalValue_Sum: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue_Sum", {-1, -2},{2, 3}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue_Sum(inputArray)
print("{'Fenix_RandomPositiveDecimalValue_Sum', {-1, -2},{2, 3},  {'true', '0'}}")
print("Fenix_RandomPositiveDecimalValue_Sum: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue_Sum", {1, 2, 3},{2, 3}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue_Sum(inputArray)
print("{'Fenix_RandomPositiveDecimalValue_Sum', {1, 2, 3},{2, 3},  {'true', '0'}}")
print("Fenix_RandomPositiveDecimalValue_Sum: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue_Sum", {1, 2, 3},{2, 3, 2, 3}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue_Sum(inputArray)
print("{'Fenix_RandomPositiveDecimalValue_Sum', {1, 2, 3},{2, 3, 2, 3},  {'true', '0'}}")
print("Fenix_RandomPositiveDecimalValue_Sum: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue_Sum", {1, 2, 3},{2, 3, 4, 4}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue_Sum(inputArray)
print("{'Fenix_RandomPositiveDecimalValue_Sum', {1, 2, 3},{2, 3, 4, 4},  {'true', '0'}}")
print("Fenix_RandomPositiveDecimalValue_Sum: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{2, 3}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{2, 3},  {'true', '0'}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '81.986'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{1, 2}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{1, 2}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '8.98'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {2},{1, 2}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {2},{1, 2}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '6.48'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{1, 1}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{1, 1}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '8.9'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{1, 1}, {"true", "1"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{1, 1}, 1}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '6.4'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{1, 1}, {"true", "1"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{1, 1}, 1}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '6.4'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{0, 1}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{0, 1}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '0.9'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{1, 0}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{1, {0}}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '8'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{0, 0}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{0, {0}}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '0'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{0, 0, 2, 3}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{0, 0, 2, 3}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '0'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{0, 2, 3, 4}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{0, 2, 3, 4}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '0'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{2, 2, 3, 4}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{2, 2, 3, 4}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '0'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{6, 6}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{6, 6}, {0}}")
print("Fenix_RandomPositiveDecimalValue: " .. response.value .. " :: Expected OK - i.e. '815587.986577'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {},{6, 10}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {},{6, 10}, {0}}")
print("Fenix_RandomPositiveDecimalValue Date: " .. response.value .. " :: Expected OK - i.e. '815587.9865775100'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{0}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{0}, {0}}")
print("Fenix_RandomPositiveDecimalValue Date: " .. response.errorMessage .. " :: Expected ERROR - there must be exact 2 function parameter. '[0]'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{}, {0}}")
print("Fenix_RandomPositiveDecimalValue Date: " .. response.errorMessage .. " :: Expected Error - there must be exact 2 function parameter but it is empty.")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{1, 2, 3}, {"true", "0"}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{1, 2, 3}, {0}}")
print("Fenix_RandomPositiveDecimalValue Date: " .. response.errorMessage .. " :: Expected Error - there must be exact 2 function parameter. '[1,2,3]'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1, 2},{2, 3}, {0}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1, 2},{2, 3}, {0}}")
print("Fenix_RandomPositiveDecimalValue Date: " .. response.errorMessage .. " :: Expected Error - array index array can only have a maximum of one value. '[1,2]'")
print("")

local inputArray = {"Fenix_RandomPositiveDecimalValue", {1},{2, 3}}
local response = Fenix_RandomPositiveDecimalValue(inputArray)
print("{'Fenix_RandomPositiveDecimalValue', {1},{2, 3}}")
print("Fenix_RandomPositiveDecimalValue Date: " .. response.errorMessage .. " :: Expected Error - there should be exactly four rows in InputTable.")
print("")





-- ***********************************************************************************
--]]