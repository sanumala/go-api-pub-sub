package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sanumala/go-api-pub-sub/services"
	"github.com/vmware/transport-go/bridge"
	"github.com/vmware/transport-go/bus"
	"github.com/vmware/transport-go/model"
	"github.com/vmware/transport-go/plank/utils"
)

func main() {

	config := &bridge.BrokerConnectorConfig{
		Username:   "guest",
		Password:   "guest",
		ServerAddr: "localhost:30080",
		UseWS:      true,
		WebSocketConfig: &bridge.WebSocketConfig{
			WSPath: "/ws",
			UseTLS: false,
		}}

	b := bus.GetBus()

	cm := b.GetChannelManager()

	c, err := b.ConnectBroker(config)
	if err != nil {
		utils.Log.Fatalf("unable to connect to %s, error: %v", config.ServerAddr, err.Error())
	}

	factSubChan := "uselessfacts"
	cm.CreateChannel(factSubChan)

	factSubHandler, _ := b.ListenOnce(factSubChan)

	cm.MarkChannelAsGalactic(factSubChan, "/queue/useless-fact-service", c)

	var wg sync.WaitGroup
	wg.Add(1)

	factSubHandler.Handle(
		func(msg *model.Message) {

			var fact services.Fact
			if err := msg.CastPayloadToType(&fact); err != nil {
				fmt.Printf("failed to cast payload: %s\n", err.Error())
			} else {
				fmt.Printf("Useless but Fact is ::::: %s", fact.Text)
			}
			wg.Done()
		},
		func(err error) {
			utils.Log.Errorf("error received on channel: %e", err)
			wg.Done()
		})

	req := &model.Request{Request: "get-useless-fact"}
	reqBytes, _ := json.Marshal(req)

	c.SendJSONMessage("/pub/queue/useless-fact-service", reqBytes)

	wg.Wait()

	// mark channels as local (unsubscribe)
	cm.MarkChannelAsLocal(factSubChan)

	// disconnect from our broker.
	c.Disconnect()
}
