package channel

import (
	"atlas-login/rest"
	"fmt"
	"github.com/Chronicle20/atlas-rest/requests"
)

const (
	WorldsResource = "worlds/"
	WorldsById     = WorldsResource + "%d"
	Resource       = WorldsById + "/channels"
)

func getBaseRequest() string {
	return requests.RootUrl("CHANNELS")
}

func requestChannelsForWorld(worldId byte) requests.Request[[]RestModel] {
	return rest.MakeGetRequest[[]RestModel](fmt.Sprintf(getBaseRequest()+Resource, worldId))
}
