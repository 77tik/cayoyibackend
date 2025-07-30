package s3_backend

import (
	"cayoyibackend/weedfilesys/util/version"
	"fmt"

	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

//  S3 后端适配器中的 AWS 会话管理逻辑，主要功能是为 S3 操作创建或复用 AWS 的 S3 客户端。下面我们逐段解释它的作用：

var (
	// 定义了一个全局缓存 s3Sessions，
	// 用于保存创建过的 S3 客户端实例（以 region + endpoint 作为 key）
	// 并通过 sessionsLock 来进行并发保护。
	s3Sessions   = make(map[string]s3iface.S3API)
	sessionsLock sync.RWMutex
)

func getSession(region string) (s3iface.S3API, bool) {
	sessionsLock.RLock()
	defer sessionsLock.RUnlock()

	sess, found := s3Sessions[region]
	return sess, found
}

func createSession(awsAccessKeyId, awsSecretAccessKey, region, endpoint string, forcePathStyle bool) (s3iface.S3API, error) {

	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	// 缓存检查
	cacheKey := fmt.Sprintf("%s|%s", region, endpoint)
	if t, found := s3Sessions[cacheKey]; found {
		return t, nil
	}

	// 配置AWS SDK Session
	config := &aws.Config{
		Region:                        aws.String(region),
		Endpoint:                      aws.String(endpoint),
		S3ForcePathStyle:              aws.Bool(forcePathStyle),
		S3DisableContentMD5Validation: aws.Bool(true),
	}

	// 设置静态AK和SK（如果提供的话）
	if awsAccessKeyId != "" && awsSecretAccessKey != "" {
		config.Credentials = credentials.NewStaticCredentials(awsAccessKeyId, awsSecretAccessKey, "")
	}

	// 创建Session，并附加自定义请求头
	// 这会在请求头中加入 User-Agent: SeaweedFS/3.5 这样的信息，方便 AWS 日志追踪请求来源。
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, fmt.Errorf("create aws session in region %s: %v", region, err)
	}
	sess.Handlers.Build.PushBack(func(r *request.Request) {
		r.HTTPRequest.Header.Set("User-Agent", "Weedfilesys/"+version.VERSION_NUMBER)
	})

	// 创建s3客户端并缓存
	t := s3.New(sess)

	s3Sessions[region] = t

	return t, nil

}

func deleteFromS3(sess s3iface.S3API, sourceBucket string, sourceKey string) (err error) {
	_, err = sess.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(sourceBucket),
		Key:    aws.String(sourceKey),
	})
	return err
}
