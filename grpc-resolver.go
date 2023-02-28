package etcd

import (
	"context"
	v3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

const Scheme = "etcd"

func init() {
	resolver.Register(&builder{})
}

type builder struct {
	client *v3.Client
}

func (b *builder) Watch(key string, cc resolver.ClientConn) {

	_ = cc.UpdateState(resolver.State{Addresses: b.getAddress(key)})

	channel := b.client.Watch(context.TODO(), key, v3.WithPrefix())
	go func() {
		for {
			select {
			case watchResponse := <-channel:
				if watchResponse.Err() != nil {
					break
				}

				_ = cc.UpdateState(resolver.State{Addresses: b.getAddress(key)})
			}
		}
	}()
}

func NewBuilder(client *v3.Client) resolver.Builder {
	return &builder{client: client}
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	key := Prefix + "/" + target.URL.Host

	defer b.Watch(key, cc)
	return &myResolver{}, nil
}

func (b *builder) Scheme() string {
	return Scheme
}

func (b *builder) getAddress(key string) (address []resolver.Address) {

	getResponse, err := b.client.Get(context.TODO(), key)
	if err != nil {
		return
	}

	for _, v := range getResponse.Kvs {
		address = append(address, resolver.Address{
			Addr: string(v.Value),
		})
	}

	return
}

type myResolver struct {
}

func (r *myResolver) ResolveNow(options resolver.ResolveNowOptions) {

}

func (r *myResolver) Close() {

}

func NewTarget(service string) string {
	return Scheme + "://" + service
}
