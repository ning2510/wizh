package etcd

const etcdPrefix = "kitex/registry-etcd"

func serviceKeyPrefix(serviceName string) string {
	return etcdPrefix + "/" + serviceName
}

func serviceKey(serviceName, addr string) string {
	return serviceKeyPrefix(serviceName) + "/" + addr
}

type instanceInfo struct {
	NetWork string            `json:"network"`
	Address string            `json:"address"`
	Weight  int               `json:"weight"`
	Tags    map[string]string `json:"tags"`
}
