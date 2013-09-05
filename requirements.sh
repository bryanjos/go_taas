export GOPATH=$(pwd)
export PATH=$PATH:$GOPATH/bin
go get github.com/emicklei/go-restful
go get github.com/emicklei/go-restful/swagger
go get labix.org/v2/mgo
go get github.com/op/go-logging