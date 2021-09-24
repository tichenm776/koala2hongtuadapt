SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm
go build -o ./wechat_server main.go
pause
