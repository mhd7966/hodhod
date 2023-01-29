package connections

import (
	"github.com/hibiken/asynq"
	"github.com/mhd7966/hodhod/configs"
)

var RedisClient *asynq.Client

func ConnectRedis() {

	RedisClient = asynq.NewClient(asynq.RedisClientOpt{Addr: configs.Cfg.Redis.Addr})

}
