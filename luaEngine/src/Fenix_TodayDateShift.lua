local date = require('date')

-- Function to shift the current date by a given number of days
function Fenix_TodayShiftDay(inputTable)

    local d = date()

    --inputTable = {
    --    functionName = "Fenix_TodayShiftDay",
    --    functionArrayPositions = {},
    --    functionValuesArray = {1},
    --    randomSeed = {true,0}
    --}

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

   -- Secure that no ArraysIndexArray is emtpty
   if (#arraysIndexTable > 0) then

        -- Convert array to string
        local tableAsString = TableToString (arraysIndexTable, ",")

        local error_message = "Error - array index is not supported. 'arraysIndexArray: " .. tableAsString .. "'"

        responseTable.success = false
        responseTable.errorMessage = error_message

        return responseTable

   end

   -- Extract FunctionArgumentsArray
   local functionArgumentsArray = inputTable[3]


   -- Handle different number of function arguments
    local shift_days = 0
   if (#functionArgumentsArray == 0) then
       shift_days = 0

   elseif (#functionArgumentsArray == 1) then
        local argumentValue
        argumentValue = functionArgumentsArray[1]

        -- secure that the value is an integer
        shift_days = convertToInteger(argumentValue)
        if shift_days == nil then

             local error_message = "Error - function argument is not an Integer: ''" .. tostring(argumentValue) .. "'"

            responseTable.success = false
            responseTable.errorMessage = error_message

           return responseTable
       end



    else
        local tableAsString = TableToString (functionArgumentsArray, ",")

        local error_message = "Error - more than 1 parameter argument. 'functionArgumentsArray: " .. tableAsString .. "'"

        responseTable.success = false
        responseTable.errorMessage = error_message

       return responseTable

    end


    -- Current date and time
    local now = date()

    -- Add days
    local futureDate = now:adddays(shift_days)

    responseTable.value = futureDate:fmt("%Y-%m-%d")

    return responseTable


end

function convertToInteger(value)
    -- First check if the value is a number or a string that contains only an integer
    if type(value) == "number" then
        value = tostring(value)  -- Convert to string for consistent pattern matching
    end

    -- Check if the string is strictly an integer (handling optional sign)
    if type(value) == "string" and value:match("^-?%d+$") then
        local num = tonumber(value)
        if num then
            print("Conversion successful:", num)
            return num
        end
    end

    -- If not a valid integer string, handle the error
    print("Failed to convert to integer.")
    return nil  -- or handle the error as needed
end

-- Examples of usage:
-- convertToInteger("-1")      -- Successfully converts
-- convertToInteger("123")     -- Successfully converts
-- convertToInteger("123.456") -- Fails to convert
-- convertToInteger("hello")   -- Fails to convert
-- convertToInteger(123.456)   -- Fails to convert
-- convertToInteger(123)       -- Successfully converts




function TableToString(tbl, sep)

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

-- Example invocation
--local result = Fenix_TodayShiftDay{"Fenix_TodayShiftDay", {}, {"-1"}, {true,0}}
--print(result.value)

