package main

import (
	"tencent"

	"github.com/hashicorp/packer/packer/plugin"
	//  "log"
	//  "os"
	//  "reflect"
	//  "fmt"
)

func main() {
	tencent.CloudAPIDebug = true

	/* Comment out the following to run the server */

	/*

	   x := struct{Foo string; Bar int }{"foo", 2}

	   v := reflect.ValueOf(x)

	   values := make([]interface{}, v.NumField())

	   for i := 0; i < v.NumField(); i++ {
	       values[i] = v.Field(i).Interface()
	       fmt.Println(values[i])
	   }

	   fmt.Println(values)

	   log.Println(tencent.SignatureString("action", "SecretId", make(map[string]string)))
	   log.Println(tencent.SignatureString("action", "SecretId", nil))
	   return // comment this out to run the plugin normally


	   // Doesn't work
	   log.Println("In Main app")
	   log.Println("PACKER_LOG")
	   log.Println(os.Getenv("PACKER_LOG"))
	   if tencent.DebugEnabled() == true {
	     log.Println("Debug enabled")
	   }
	*/

	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(tencent.Builder))
	server.Serve()

}
