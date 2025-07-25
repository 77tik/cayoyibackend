package cluster

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"my_backend/weedfilesys/glog"
	"my_backend/weedfilesys/pb"
	"my_backend/weedfilesys/pb/filer_pb"
	"my_backend/weedfilesys/util"
	"time"
)

// 主要用于跨多个 Filer 或服务实例之间，协调资源访问:
// 比如多个 Filer 都想执行某个任务（比如分片压缩、目录迁移、主目录挂载等），需要先抢占锁，谁抢到谁执行，其它等待或退出。

// 锁客户端
type LockClient struct {
	grpcDialOption  grpc.DialOption
	maxLockDuration time.Duration // 最长锁定时间
	sleepDuration   time.Duration // 重试间隔
	seedFiler       pb.ServerAddress
}

func NewLockClient(grpcDialOption grpc.DialOption, seedFiler pb.ServerAddress) *LockClient {
	return &LockClient{
		grpcDialOption:  grpcDialOption,
		maxLockDuration: 5 * time.Second,
		sleepDuration:   2000 * time.Millisecond,
		seedFiler:       seedFiler,
	}
}

type LiveLock struct {
	key            string           // 锁的唯一键名
	renewToken     string           // 用于续约的token
	hostFiler      pb.ServerAddress // 当前的锁定的filer
	cancelCH       chan struct{}    // 停止控制通道
	isLocked       bool
	lc             *LockClient
	owner          string // 当前锁拥有者
	expireAtNs     int64
	grpcDialOption grpc.DialOption
	self           string
}

// 创建一个5s时间的短锁，里面包含重试，直到获取或者超时
func (lc *LockClient) NewShortLiveLock(key string, owner string) (lock *LiveLock) {
	lock = &LiveLock{
		key:            key,
		hostFiler:      lc.seedFiler,
		cancelCH:       make(chan struct{}),
		expireAtNs:     time.Now().Add(5 * time.Second).UnixNano(),
		grpcDialOption: lc.grpcDialOption,
		self:           owner,
		lc:             lc,
	}
	lock.retryUntilLocked(5 * time.Second)
	return lock
}

func (lock *LiveLock) retryUntilLocked(lockDuration time.Duration) {
	util.RetryUntil("create lock:"+lock.key, func() error {
		return lock.AttemptToLock(lockDuration)
	}, func(err error) (shouldContinue bool) {
		if err != nil {
			glog.Warningf("create lock %s: %s", lock.key, err)
		}
		return lock.renewToken == ""
	})
}

func (lock *LiveLock) AttemptToLock(lockDuration time.Duration) error {
	errorMessage, err := lock.doLock(lockDuration)
	if err != nil {
		time.Sleep(time.Second)
		return err
	}
	if errorMessage != "" {
		time.Sleep(time.Second)
		return fmt.Errorf("%v", errorMessage)
	}
	lock.isLocked = true
	return nil
}

func (lock *LiveLock) doLock(lockDuration time.Duration) (errorMessage string, err error) {
	err = pb.WithFilerClient(false, 0, lock.hostFiler, lock.grpcDialOption, func(client filer_pb.WeedfilesysFilerClient) error {
		resp, err := client.DistributedLock(context.Background(), &filer_pb.LockRequest{
			Name:          lock.key,
			SecondsToLock: int64(lockDuration.Seconds()),
			RenewToken:    lock.renewToken,
			IsMoved:       false,
			Owner:         lock.self,
		})
		if err == nil && resp != nil {
			lock.renewToken = resp.RenewToken
		} else {
			lock.renewToken = " "
		}
		if resp != nil {
			errorMessage = resp.Error
			if resp.LockHostMovedTo != "" {
				lock.hostFiler = pb.ServerAddress(resp.LockHostMovedTo)
				lock.lc.seedFiler = lock.hostFiler
			}
			if resp.LockOwner != "" {
				lock.owner = resp.LockOwner
				// fmt.Printf("lock %s owner: %s\n", lock.key, lock.owner)
			} else {
				// fmt.Printf("lock %s has no owner\n", lock.key)
				lock.owner = ""
			}
		}
		return err
	})
	return
}
