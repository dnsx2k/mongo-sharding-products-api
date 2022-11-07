package productshttphandler

import (
	"net/http"

	"github.com/dnsx2k/mongo-sharding-products-api/pkg/lookupclient"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type HTTPHandlerContext struct {
	lookupClient       *lookupclient.ClientCtx
	mongoClientPrimary *mongo.Client
	mongoClientHot     *mongo.Client
}

func New(mongoPrimary, mongoHot *mongo.Client, lookupClient *lookupclient.ClientCtx) *HTTPHandlerContext {
	return &HTTPHandlerContext{
		lookupClient:       lookupClient,
		mongoClientPrimary: mongoPrimary,
		mongoClientHot:     mongoHot,
	}
}

// Setup - setup for HTTP gin handler
func (sc *HTTPHandlerContext) Setup(route gin.IRouter) {
	route.GET("products/:id", sc.handleGet)
}

func (sc *HTTPHandlerContext) handleGet(gCtx *gin.Context) {
	key := gCtx.Param("id")
	location, err := sc.lookupClient.GetLookup(key)
	if err != nil {
		gCtx.JSON(http.StatusBadRequest, err)
		return
	}

	dbLocation := sc.mongoClientPrimary
	switch location {
	case "hot":
		dbLocation = sc.mongoClientHot
	}

	// TODO: Query mongo
}
