SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ./wechat_server main.go
pause
