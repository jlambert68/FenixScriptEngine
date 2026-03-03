-- HappyLuaTime placeholder
-- Usage: {{HappyLuaTime()}}
-- Returns: "My name is Lua and the time is <HH:MM:SS>"

function HappyLuaTime(inputTable)
    local responseTable = {
        success = true,
        value = "",
        errorMessage = ""
    }

    if type(inputTable) ~= "table" then
        responseTable.success = false
        responseTable.errorMessage = "Error - input must be a table."
        return responseTable
    end

    if #inputTable ~= 4 then
        responseTable.success = false
        responseTable.errorMessage = "Error - there should be exactly four rows in InputTable."
        return responseTable
    end

    local arrayIndexes = inputTable[2]
    local functionArgs = inputTable[3]

    if type(arrayIndexes) ~= "table" or #arrayIndexes > 0 then
        responseTable.success = false
        responseTable.errorMessage = "Error - array index is not supported."
        return responseTable
    end

    if type(functionArgs) ~= "table" or #functionArgs ~= 0 then
        responseTable.success = false
        responseTable.errorMessage = "Error - HappyLuaTime() takes no function arguments."
        return responseTable
    end

    local currentTime = os.date("%H:%M:%S")
    responseTable.value = "My name is Lua and the time is " .. currentTime

    return responseTable
end

