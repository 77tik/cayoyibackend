package cluster

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"my_backend/weedfilesys/cluster/lock_manager"
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

func (lc *LockClient) StartLongLivedLock(key string, owner string, onLockOwnerChange func(newLockOwner string)) (lock *LiveLock) {
	lock = &LiveLock{
		key:       key,
		hostFiler: lc.seedFiler,
		cancelCH:  make(chan struct{}),
		// TODO: 长时间应该提取出变量
		expireAtNs:     time.Now().Add(lock_manager.LiveLockTTL).UnixNano(),
		grpcDialOption: lc.grpcDialOption,
		self:           owner,
		lc:             lc,
	}
	go func() {
		isLocked := false
		lockOwner := ""
		for {
			// 如果当前已经持有锁，就执行“续约”（本质上也是重新申请）：
			//
			//如果失败，说明锁丢失了 → isLocked = false
			//
			//如果没持有锁，就尝试“获取锁”：
			//
			//如果成功，则将 isLocked = true
			if isLocked {
				if err := lock.AttemptToLock(lock_manager.LiveLockTTL); err != nil {
					glog.V(0).Infof("Lost lock %s: %v", key, err)
					isLocked = false
				}
			} else {
				if err := lock.AttemptToLock(lock_manager.LiveLockTTL); err == nil {
					isLocked = true
				}
			}

			// lock.LockOwner() 获取当前锁的实际持有者（从中心节点获得的最新状态）。
			//
			//如果发现和本地记录 lockOwner 不一致，就说明锁被他人抢占或释放被他人获取了。
			//
			//此时触发 onLockOwnerChange 回调。
			if lockOwner != lock.LockOwner() && lock.LockOwner() != "" {
				glog.V(0).Infof("Lock owner changed from %s to %s", lockOwner, lock.LockOwner())
				onLockOwnerChange(lock.LockOwner())
				lockOwner = lock.LockOwner()
			}

			// 如果外部调用了 lock.Stop() 关闭 cancelCH，协程会退出。
			//
			//否则每轮循环休眠一段时间（RenewInterval，如 3 秒），继续下一轮尝试。
			select {
			case <-lock.cancelCH:
				return
			default:
				time.Sleep(lock_manager.RenewInterval)
			}
		}
	}()
	return
}

func (lock *LiveLock) LockOwner() string {
	return lock.owner
}
