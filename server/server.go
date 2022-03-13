package main

import (
	"os"

	"github.com/sanumala/go-api-pub-sub/services"
	"github.com/vmware/transport-go/plank/pkg/server"
	"github.com/vmware/transport-go/plank/utils"
)

func main() {
	serverConfig, err := server.CreateServerConfig()
	if err != nil {
		utils.Log.Fatalln(err)
		return
	}

	platformServer := server.NewPlatformServer(serverConfig)

	if err = platformServer.RegisterService(services.NewUselessFactService(), services.UselessFactServiceChannel); err != nil {
		utils.Log.Fatalln(err)
		return
	}

	syschan := make(chan os.Signal, 1)

	platformServer.StartServer(syschan)
}
