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


local function ProcessControlledUniqueId(input, inputTable)

    local function randomString(length, upper, seedValue)
        if seedValue ~= nil then
            math.randomseed(seedValue)
        end

        local chars = upper and 'ABCDEFGHIJKLMNOPQRSTUVWXYZ' or 'abcdefghijklmnopqrstuvwxyz'
        local result = ''
        for i = 1, length do
            local randIndex = math.random(#chars)
            result = result .. chars:sub(randIndex, randIndex)
        end
        return result
    end

    local function replaceDateTime(pattern, format)
        input = input:gsub(pattern, os.date(format))
    end

    local function replaceRandomNumber(pattern, seedValue)
        
        math.randomseed(seedValue)

        input = input:gsub(pattern, function(n) 
            return string.format("%0"..#n.."d", math.random(10^#n-1))
        end)
    end

    local function replaceRandomString(pattern, upper, seedValue)

        math.randomseed(seedValue)

        input = input:gsub(pattern, function(n, seed) 
            return randomString(tonumber(n), upper, tonumber(seed))
        end)
    end

    -- Extract from input table
    local arrayPositionTable  = inputTable[1]
    local seedValue = inputTable[2]

    -- Replace date patterns
    replaceDateTime("%%YYYY%-MM%-DD%%", "%Y-%m-%d")
    replaceDateTime("%%YYYYMMDD%%", "%Y%m%d")
    replaceDateTime("%%YYMMDD%%", "%y%m%d")

    -- Replace time patterns
    replaceDateTime("%%hh:mm:ss%%", "%H:%M:%S")
    replaceDateTime("%%hh%.mm%.ss%%", "%H.%M.%S")
    replaceDateTime("%%hhmmss%%", "%H%M%S")
    replaceDateTime("%%hhmm%%", "%H%M")


    -- Replace time with milliseconds
    input = input:gsub("%%hh:mm:ss%%", function()
        return os.date("%H:%M:%S")
    end)
    input = input:gsub("%%hh.mm.ss%%", function()
        return os.date("%H.%M.%S")
    end)

    -- Create seed value using arrayPositionTable and entropy
    local seedValue = arrayPositionTable[1] + seedValue


    -- Replace random number pattern
    replaceRandomNumber("%%(n+)%%", seedValue)

    -- Replace random string patterns with seeding
    replaceRandomString("%%a%((%d+);%s*(%d+)%)%%", false, seedValue)
    replaceRandomString("%%A%((%d+);%s*(%d+)%)%%", true, seedValue)

    return input
end


function Fenix_ControlledUniqueId(inputTable)
    
    local responseTable = {
        success = true,
        value = "",
        errorMessage = ""
    }

    -- ExtractInput
    local arrayPositionTable  = inputTable[2]
    local textToProcess = inputTable[3][1]
    local entropyTable = inputTable[4]

    -- Secure that 'textToProcess' is of string type 
    if type(textToProcess) ~= "string" then

        responseTable.success = false
        responseTable.errorMessage = "textToProcess must be a string, got " .. type(textToProcess)

        return responseTable
    end

    -- Secure that arrayPositionTable is not empty, if so, the use '1'
    if #arrayPositionTable == 0 then
        arrayPositionTable = {1}
    end

    -- More then one position is not allowed in arrayPositionTable
    if #arrayPositionTable > 1 then

        local tableAsString = tableToString(arrayPositionTable, ",")
        local error_message = "Error - there cant be more than 1 value in 'arrayPositionTable'. '" .. tableAsString .. "'"

        responseTable.success = false
        responseTable.errorMessage = tableAsString

        return responseTable
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

    -- Create a new Input array
    local newInputTable = {arrayPositionTable, entropyValue}

    local result = ProcessControlledUniqueId(textToProcess, newInputTable)

    responseTable.value = result

    return responseTable

end





-- Example usage

local inputString = "Date: %YYYY-MM-DD%"
local inputTable = {"Fenix_ControlledUniqueId", {}, {inputString}, {true, 0}}
local result = Fenix_ControlledUniqueId(inputTable)
print(inputString)
print(result)

local inputString = "Date: %YYYY-MM-DD%, Date: %YYYYMMDD%, Date: %YYMMDD%, Time: %hh:mm:ss%, Time: %hhmmss%, Time: %hhmm%, Random Number: %nnnnn%, Random String: %a(5; 11)%, Random String Uppercase: %A(5; 10)%, Time: %hh:mm:ss%, Time: %hh.mm.ss% "
local inputTable = {"ControlledUniqueId", {0}, {inputString}, {0}}
local result = Fenix_ControlledUniqueId(inputTable)
print(inputString)
print(result)
