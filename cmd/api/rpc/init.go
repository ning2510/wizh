package rpc

import "wizh/pkg/viper"

func init() {
	commentConfig := viper.InitConf("comment")
	InitComment(&commentConfig)

	favoriteConfig := viper.InitConf("favorite")
	InitFavorite(&favoriteConfig)

	messageConfig := viper.InitConf("message")
	InitMessage(&messageConfig)

	relationConfig := viper.InitConf("relation")
	InitRelation(&relationConfig)

	userConfig := viper.InitConf("user")
	InitUser(&userConfig)

	videoConfig := viper.InitConf("video")
	InitVideo(&videoConfig)
}
