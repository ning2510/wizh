package rpc

import (
	"context"
	"fmt"
	"time"
	"wizh/kitex/kitex_gen/favorite"
	"wizh/kitex/kitex_gen/favorite/favoriteservice"
	"wizh/pkg/etcd"
	"wizh/pkg/middleware"
	"wizh/pkg/viper"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

var favoriteClient favoriteservice.Client

func InitFavorite(config *viper.Config) {
	etcdAddr := fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	serviceName := config.Viper.GetString("server.name")

	// 服务发现
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}

	// 初始化 etcd
	favoriteClient, err = favoriteservice.NewClient(
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

func FavoriteVideoAction(ctx context.Context, req *favorite.FavoriteVideoActionRequest) (*favorite.FavoriteVideoActionResponse, error) {
	return favoriteClient.FavoriteVideoAction(ctx, req)
}

func FavoriteVideoList(ctx context.Context, req *favorite.FavoriteVideoListRequest) (*favorite.FavoriteVideoListResponse, error) {
	return favoriteClient.FavoriteVideoList(ctx, req)
}

func FavoriteCommentAction(ctx context.Context, req *favorite.FavoriteCommentActionRequest) (*favorite.FavoriteCommentActionResponse, error) {
	return favoriteClient.FavoriteCommentAction(ctx, req)
}
