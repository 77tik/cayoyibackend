package dash

import (
	"cayoyibackend/weedfilesys/wdclient"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

type AdminServer struct {
	masterClient    *wdclient.MasterClient
	templateFS      http.FileSystem
	dataDir         string
	grpcDialOption  grpc.DialOption
	cacheExpiration time.Duration
	lastCacheUpdate time.Time
	cachedTopology  *ClusterTopology

	// Filer discovery and caching
	cachedFilers         []string
	lastFilerUpdate      time.Time
	filerCacheExpiration time.Duration

	// Credential management
	credentialManager *credential.CredentialManager

	// Configuration persistence
	configPersistence *ConfigPersistence

	// Maintenance system
	maintenanceManager *maintenance.MaintenanceManager

	// Topic retention purger
	topicRetentionPurger *TopicRetentionPurger

	// Worker gRPC server
	workerGrpcServer *WorkerGrpcServer
}
