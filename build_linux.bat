@ECHO off

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

ECHO "build dispatcher"
go build github.com/xiaonanln/goworld/components/dispatcher
ECHO "build gate"
go build github.com/xiaonanln/goworld/components/gate
ECHO "build game"
go build -o game github.com/yidaowan/tadpole-srv/tadpole

PAUSE
