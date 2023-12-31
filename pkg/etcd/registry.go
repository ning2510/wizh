package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"wizh/pkg/zap"

	"github.com/cloudwego/kitex/pkg/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ttlKey     = "KITEX_ETCD_REGISTRY_LEASE_TTL"
	defaultTTL = 60
)

type registerMeta struct {
	leaseId clientv3.LeaseID
	ctx     context.Context
	cancel  context.CancelFunc
}

type etcdRegistry struct {
	etcdClient *clientv3.Client
	leaseTTL   int64
	meta       *registerMeta
}

func NewEtcdRegistry(endpoints []string) (registry.Registry, error) {
	return NewEtcdRegistryWithAuth(endpoints, "", "")
}

func NewEtcdRegistryWithAuth(endpoints []string, username, password string) (registry.Registry, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
		Username:  username,
		Password:  password,
	})

	if err != nil {
		return nil, err
	}

	return &etcdRegistry{
		etcdClient: etcdClient,
		leaseTTL:   getTTL(),
	}, nil
}

func getTTL() int64 {
	var ttl int64 = defaultTTL
	if str, ok := os.LookupEnv(ttlKey); ok {
		if t, err := strconv.ParseInt(str, 10, 64); err == nil {
			ttl = t
		}
	}
	return ttl
}

func validateRegistryInfo(info *registry.Info) error {
	if info.ServiceName == "" {
		return fmt.Errorf("missing service name in Register")
	}

	if info.Addr == nil {
		return fmt.Errorf("missing addr in Register")
	}
	return nil
}

// register a server with given registry info
func (e *etcdRegistry) Register(info *registry.Info) error {
	if err := validateRegistryInfo(info); err != nil {
		return err
	}

	leaseId, err := e.grantLease()
	if err != nil {
		return err
	}

	if err := e.register(info, leaseId); err != nil {
		return err
	}

	meta := registerMeta{
		leaseId: leaseId,
	}
	meta.ctx, meta.cancel = context.WithCancel(context.Background())
	if err := e.keepalive(&meta); err != nil {
		return err
	}
	e.meta = &meta
	return nil
}

func (e *etcdRegistry) register(info *registry.Info, leaseId clientv3.LeaseID) error {
	val, err := json.Marshal(&instanceInfo{
		NetWork: info.Addr.Network(),
		Address: info.Addr.String(),
		Weight:  info.Weight,
		Tags:    info.Tags,
	})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err = e.etcdClient.Put(ctx, serviceKey(info.ServiceName, info.Addr.String()), string(val), clientv3.WithLease(leaseId))
	return err
}

// deregister a server with given registry info
func (e *etcdRegistry) Deregister(info *registry.Info) error {
	if info.ServiceName == "" {
		return fmt.Errorf("missing service name in Deregister")
	}

	if err := e.deregister(info); err != nil {
		return err
	}
	e.meta.cancel()
	return nil
}

func (e *etcdRegistry) deregister(info *registry.Info) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := e.etcdClient.Delete(ctx, serviceKey(info.ServiceName, info.Addr.String()))
	return err
}

func (e *etcdRegistry) grantLease() (clientv3.LeaseID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	res, err := e.etcdClient.Grant(ctx, e.leaseTTL)
	if err != nil {
		return clientv3.NoLease, err
	}
	return res.ID, nil
}

func (e *etcdRegistry) keepalive(meta *registerMeta) error {
	logger := zap.InitLogger()
	keetAlive, err := e.etcdClient.KeepAlive(meta.ctx, meta.leaseId)
	if err != nil {
		return err
	}

	go func() {
		logger.Infof("Start keepalive lease %d for etcd registry", meta.leaseId)
		for range keetAlive {
			select {
			case <-meta.ctx.Done():
				break
			default:
			}
		}
		logger.Infof("Stop keepalive lease %d for etcd registry", meta.leaseId)
	}()
	return nil
}
