package connections

import (
	"github.com/abr-ooo/hodhod/configs"
	"github.com/hibiken/asynq"
)

var RedisClient *asynq.Client

func ConnectRedis() {

	RedisClient = asynq.NewClient(asynq.RedisClientOpt{Addr: configs.Cfg.Redis.Addr})

}
