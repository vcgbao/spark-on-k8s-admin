package config

type K8sConfig struct {
	ServiceHost string
	ServicePort int
	Token       string
	InCluster   bool
}
