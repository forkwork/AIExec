package server

import (
	"aiexec-plugin/internal/cluster"
	"aiexec-plugin/internal/core/plugin_daemon/backwards_invocation/transaction"
)

type App struct {
	// cluster instance of this node
	// schedule all the tasks related to the cluster, like request direct
	cluster *cluster.Cluster

	// endpoint handler
	// customize behavior of endpoint
	endpointHandler EndpointHandler

	// aws transaction handler
	// accept aws transaction request and forward to the plugin daemon
	awsTransactionHandler *transaction.AWSTransactionHandler
}
