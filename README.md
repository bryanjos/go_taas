Tamagotchi as a service
========================

Written in GO and uses MongoDB as a database. 

The requirements are as follows. You can run the script, requirements.sh to set the gopath to the current directory and get the required libraries:

* github.com/emicklei/go-restful
* github.com/emicklei/go-restful/swagger
* labix.org/v2/mgo
* github.com/op/go-logging



To run it, you can use the script, start.sh. It's hard coded to use port 5000, but you can change to port in main.go