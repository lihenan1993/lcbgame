set GOOS=linux
set GOARCH=amd64

for /F %%i in ('git rev-parse HEAD') do ( set commitid=%%i)
echo commitid=%commitid%

go build -ldflags "-X main.BuildTime=`` -X main.CommitID=%commitid% -X main.LogLevel="info" -s -w" -o ./dist/slots_server main.go
