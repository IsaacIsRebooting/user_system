package utils

import (
	"golang.org/x/net/context"
	"my_user_system/conf"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()
	conf.InitConfig()
	res, err := GetRedisCli().Set(ctx, "2222", "user2", 60*time.Second).Result()
	if err != nil {
		t.Errorf("redis set err:%v", err)
	}
	t.Logf("res=%s", res)

	val := GetRedisCli().Get(ctx, "2222").Val()
	if err != nil {
		t.Errorf("redis get err:%v", err)
	}
	t.Logf("val=%s", val)

}

func TestGetSession(t *testing.T) {
	ctx := context.Background()
	conf.InitConfig()
	val, err := GetRedisCli().Get(ctx, "session_abfc434dfbd8fd817449ce58438f8413").Result()
	if err != nil {
		t.Errorf("err=%v", err)
	}
	t.Logf("val=%s", val)
}
