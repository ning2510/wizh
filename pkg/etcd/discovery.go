package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const defaultWeight = 10

type etcdResolver struct {
	etcdClient *clientv3.Client
}

func NewEtcdResolver(endpoints []string) (discovery.Resolver, error) {
	return NewEtcdResolverWithAuth(endpoints, "", "")
}

func NewEtcdResolverWithAuth(endpoints []string, username, password string) (discovery.Resolver, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
		Username:  username,
		Password:  password,
	})

	if err != nil {
		return nil, err
	}
	return &etcdResolver{
		etcdClient: etcdClient,
	}, nil
}

// Target implements the Resolver interface
func (e *etcdResolver) Target(ctx context.Context, target rpcinfo.EndpointInfo) (description string) {
	return target.ServiceName()
}

// Resolve implements the Resolver interface
func (e *etcdResolver) Resolve(ctx context.Context, desc string) (discovery.Result, error) {
	prefix := serviceKeyPrefix(desc)
	res, err := e.etcdClient.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return discovery.Result{}, err
	}

	var info instanceInfo
	var eps []discovery.Instance
	for _, kv := range res.Kvs {
		err := json.Unmarshal(kv.Value, &info)
		if err != nil {
			continue
		}

		weight := info.Weight
		if weight <= 0 {
			weight = defaultWeight
		}
		eps = append(eps, discovery.NewInstance(info.NetWork, info.Address, weight, info.Tags))
	}

	if len(eps) == 0 {
		return discovery.Result{}, fmt.Errorf("no instance remains for %v", desc)
	}
	return discovery.Result{
		Cacheable: true,
		CacheKey:  desc,
		Instances: eps,
	}, nil
}

// Diff implements the Resolver interface
func (e *etcdResolver) Diff(cacheKey string, prev, next discovery.Result) (discovery.Change, bool) {
	return discovery.DefaultDiff(cacheKey, prev, next)
}

// Name implements the Resolver interface
func (e *etcdResolver) Name() string {
	return "etcd"
}
