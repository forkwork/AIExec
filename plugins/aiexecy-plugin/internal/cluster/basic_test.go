package cluster

import "aiexec-plugin/internal/utils/cache"

func clearClusterState() {
	cache.Del(CLUSTER_STATUS_HASH_MAP_KEY)
	cache.Del(PREEMPTION_LOCK_KEY)
}
