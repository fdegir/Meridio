package client

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/cls"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/kernel"

	"github.com/sirupsen/logrus"
)

type NetworkServiceClient struct {
	Id                         string
	NetworkServiceName         string
	NetworkServiceEndpointName string
	Labels                     map[string]string
	Connection                 *networkservice.Connection
	nsmgrClient                NSMgrClient
	InterfaceName              string
}

type NSMgrClient interface {
	Request(*networkservice.NetworkServiceRequest) (*networkservice.Connection, error)
	Close(*networkservice.Connection) (*empty.Empty, error)
}

// Request
func (nsc *NetworkServiceClient) Request() {
	request := nsc.prepareRequest()
	for true {
		var err error
		nsc.Connection, err = nsc.nsmgrClient.Request(request)
		if err != nil {
			time.Sleep(15 * time.Second)
			logrus.Errorf("Network Service Client: Request err: %v", err)
			continue
		}
		// TODO interface adverstise
		break
	}
}

// Close -
func (nsc *NetworkServiceClient) Close() {
	var err error
	_, err = nsc.nsmgrClient.Close(nsc.Connection)
	if err != nil {
		logrus.Errorf("Network Service Client: Close err: %v", err)
	}
}

func (nsc *NetworkServiceClient) prepareRequest() *networkservice.NetworkServiceRequest {
	// TODO
	outgoingMechanism := &networkservice.Mechanism{
		Cls:  cls.LOCAL,
		Type: kernel.MECHANISM,
		Parameters: map[string]string{
			kernel.NetNSURL:         (&url.URL{Scheme: "file", Path: "/proc/thread-self/ns/net"}).String(),
			kernel.InterfaceNameKey: nsc.InterfaceName,
		},
	}

	request := &networkservice.NetworkServiceRequest{
		Connection: &networkservice.Connection{
			Id:                         nsc.Id,
			NetworkService:             nsc.NetworkServiceName,
			Labels:                     nsc.Labels,
			NetworkServiceEndpointName: nsc.NetworkServiceEndpointName,
		},
		MechanismPreferences: []*networkservice.Mechanism{
			outgoingMechanism,
		},
	}
	return request
}

// NewnetworkServiceClient
func NewNetworkServiceClient(networkServiceName string, nsmgrClient NSMgrClient) *NetworkServiceClient {
	identifier := rand.Intn(100)
	id := fmt.Sprintf("%d", identifier)

	// TODO
	randomID := rand.Intn(1000)
	interfaceName := "nsc" + strconv.Itoa(randomID)

	networkServiceClient := &NetworkServiceClient{
		Id:                 id,
		NetworkServiceName: networkServiceName,
		nsmgrClient:        nsmgrClient,
		InterfaceName:      interfaceName,
	}

	return networkServiceClient
}
