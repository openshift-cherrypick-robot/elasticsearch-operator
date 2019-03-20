package k8shandler

import (
	"fmt"

	v1alpha1 "github.com/openshift/elasticsearch-operator/pkg/apis/elasticsearch/v1alpha1"
)

const (
	modeUnique    = "unique"
	modeSharedOps = "shared_ops"
	defaultMode   = modeSharedOps

	defaultMasterCPULimit     = "100m"
	defaultMasterCPURequest   = "100m"
	defaultCPULimit           = "4000m"
	defaultCPURequest         = "100m"
	defaultMemoryLimit        = "4Gi"
	defaultMemoryRequest      = "1Gi"
	elasticsearchDefaultImage = "quay.io/openshift/origin-logging-elasticsearch5"

	maxMasterCount = 3

	elasticsearchCertsPath  = "/etc/openshift/elasticsearch/secret"
	elasticsearchConfigPath = "/usr/share/java/elasticsearch/config"
	heapDumpLocation        = "/elasticsearch/persistent/heapdump.hprof"
)

func kibanaIndexMode(mode string) (string, error) {
	if mode == "" {
		return defaultMode, nil
	}
	if mode == modeUnique || mode == modeSharedOps {
		return mode, nil
	}
	return "", fmt.Errorf("invalid kibana index mode provided [%s]", mode)
}

func esUnicastHost(clusterName, namespace string) string {
	return fmt.Sprintf("%v-cluster.%v.svc", clusterName, namespace)
}

func rootLogger() string {
	return "rolling"
}

func calculateReplicaCount(dpl *v1alpha1.Elasticsearch) int {
	dataNodeCount := int((getDataCount(dpl)))
	repType := dpl.Spec.RedundancyPolicy
	switch repType {
	case v1alpha1.FullRedundancy:
		return dataNodeCount - 1
	case v1alpha1.MultipleRedundancy:
		return (dataNodeCount - 1) / 2
	case v1alpha1.SingleRedundancy:
		return 1
	case v1alpha1.ZeroRedundancy:
		return 0
	default:
		return 1
	}
}
