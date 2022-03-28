package tencent

import "github.com/hashicorp/packer/packer/plugin"

func main() {

	// testCreateVM()

	//   log.Println("In tencent main")
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(Builder))
	server.Serve()
}
