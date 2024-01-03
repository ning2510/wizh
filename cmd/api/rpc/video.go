package rpc

import (
	"context"
	"fmt"
	"time"
	"wizh/kitex/kitex_gen/video"
	"wizh/kitex/kitex_gen/video/videoservice"
	"wizh/pkg/etcd"
	"wizh/pkg/middleware"
	"wizh/pkg/viper"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

var videoClient videoservice.Client

func InitVideo(config *viper.Config) {
	etcdAddr := fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	serviceName := config.Viper.GetString("server.name")

	// 服务发现
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}

	// 初始化 etcd
	videoClient, err = videoservice.NewClient(
		serviceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithMiddleware(middleware.ServerMiddleware),
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(180*time.Second),            // rpc timeout
		client.WithConnectTimeout(30000*time.Millisecond), // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)
	if err != nil {
		panic(err)
	}
}

func Feed(ctx context.Context, req *video.FeedRequest) (*video.FeedResponse, error) {
	return videoClient.Feed(ctx, req)
}

func PublishAction(ctx context.Context, req *video.PublishActionRequest) (*video.PublishActionResponse, error) {
	return videoClient.PublishAction(ctx, req)
}

func PublishList(ctx context.Context, req *video.PublishListRequest) (*video.PublishListResponse, error) {
	return videoClient.PublishList(ctx, req)
}

func PublishInfo(ctx context.Context, req *video.PublishInfoRequest) (*video.PublishInfoResponse, error) {
	return videoClient.PublishInfo(ctx, req)
}

func PublishDelete(ctx context.Context, req *video.PublishDeleteRequest) (*video.PublishDeleteResponse, error) {
	return videoClient.PublishDelete(ctx, req)
}
