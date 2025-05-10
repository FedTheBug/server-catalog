package main

import "github.com/server-catalog/cmd"

//@title           Server Catalog API
//@version         1.0
//@description     A server catalog service API documentation
//@termsOfService  http://swagger.io/terms/
//
//@contact.name   API Support
//@contact.url    http://www.swagger.io/support
//@contact.email  support@swagger.io
//
//@license.name  Apache 2.0
//@license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
//@host      localhost:8080
//@BasePath  /api/v1
//
// @securityDefinitions.apikey AppKeyAuth
// @in header
// @name App-key

func main() {
	cmd.Execute()
}
