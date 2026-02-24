RunLuaUnitTests:
	cd luaEngine && luatest -v

UnitTests_ScriptEngine:
	go test -v ./scriptEngine

UnitTests_ScriptEngineCoverage:
	go test ./scriptEngine -coverprofile=coverage.out && go tool cover -func=coverage.out && go tool cover -html=coverage.out

