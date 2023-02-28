package etcd

import (
	"context"
	v3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

const Prefix = "/myetcd"

func NewClient(url, username, password string) *v3.Client {

	client, err := v3.New(v3.Config{
		Endpoints:            strings.Split(url, ","),
		AutoSyncInterval:     0,
		DialTimeout:          time.Second * 5,
		DialKeepAliveTime:    0,
		DialKeepAliveTimeout: 0,
		MaxCallSendMsgSize:   0,
		MaxCallRecvMsgSize:   0,
		TLS:                  nil,
		Username:             username,
		Password:             password,
		RejectOldCluster:     false,
		DialOptions:          nil,
		Context:              nil,
		Logger:               nil,
		LogConfig:            nil,
		PermitWithoutStream:  false,
	})
	if err != nil {
		panic(err)
	}

	return client
}

func Register(client *v3.Client, service, address string) context.CancelFunc {

	leaseGrantResponse, err := client.Grant(context.TODO(), 2)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	leaseKeepAliveResponse, err := client.KeepAlive(ctx, leaseGrantResponse.ID)
	if err != nil {
		panic(err)
	}

	go func() {
	B:
		for {
			select {
			case <-leaseKeepAliveResponse:

			case <-ctx.Done():
				break B
			}
		}
	}()

	_, err = client.Put(context.TODO(), Prefix+"/"+service, address, v3.WithLease(leaseGrantResponse.ID))
	if err != nil {
		panic(err)
	}

	return cancel
}
