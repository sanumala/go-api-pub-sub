package services

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/vmware/transport-go/model"
	"github.com/vmware/transport-go/service"
)

const (
	UselessFactServiceChannel = "useless-fact-service"
)

type Fact struct {
	Id        string `json:"id"`
	Text      string `json:"text"`
	SourceUrl string `json:"source_url"`
}

type UselessFactService struct{}

func NewUselessFactService() *UselessFactService {
	return &UselessFactService{}
}

func (ufs *UselessFactService) Init(core service.FabricServiceCore) error {
	core.SetDefaultJSONHeaders()
	return nil
}

func (ufs *UselessFactService) HandleServiceRequest(request *model.Request, core service.FabricServiceCore) {
	switch request.Request {
	case "get-useless-fact":
		ufs.getUselessFact(request, core)
	default:
		core.HandleUnknownRequest(request)
	}
}

func (ufs *UselessFactService) getUselessFact(request *model.Request, core service.FabricServiceCore) {

	core.RestServiceRequest(&service.RestServiceRequest{
		Uri:    "https://uselessfacts.jsph.pl/random.json?language=en",
		Method: "GET",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		ResponseType: reflect.TypeOf(&Fact{}),
	}, func(response *model.Response) {

		core.SendResponse(request, response.Payload.(*Fact))

	}, func(response *model.Response) {

		fabricError := service.GetFabricError("Get Useless Fact API Call Failed", response.ErrorCode, response.ErrorMessage)
		core.SendErrorResponseWithPayload(request, response.ErrorCode, response.ErrorMessage, fabricError)
	})
}

func (ufs *UselessFactService) GetRESTBridgeConfig() []*service.RESTBridgeConfig {
	return []*service.RESTBridgeConfig{
		{
			ServiceChannel: UselessFactServiceChannel,
			Uri:            "/rest/uselessfact",
			Method:         http.MethodGet,
			AllowHead:      true,
			AllowOptions:   true,
			FabricRequestBuilder: func(w http.ResponseWriter, r *http.Request) model.Request {

				return model.Request{
					Id:                &uuid.UUID{},
					Request:           "get-useless-fact",
					BrokerDestination: nil,
				}
			},
		},
	}
}

func (ufs *UselessFactService) OnServerShutdown() {}

func (ufs *UselessFactService) OnServiceReady() chan bool {

	readyChan := make(chan bool, 1)
	readyChan <- true
	return readyChan
}
