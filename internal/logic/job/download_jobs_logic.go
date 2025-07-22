package job

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"io/fs"
	"my_backend/internal/svc"
	"my_backend/internal/types"
	"os"
	"path/filepath"
)

type DownloadJobsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 作业文件下载
func NewDownloadJobsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadJobsLogic {
	return &DownloadJobsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadJobsLogic) DownloadJobs(req *types.DownloadJobsReq) (resp []byte, err error) {
	baseDir := "./work" // 可改为配置路径

	// 将 jobIds 转为 map，方便匹配
	targetSet := make(map[string]bool)
	for _, id := range req.JobNumbers {
		targetSet[id] = true
	}

	// 找出目标路径（保留分类名）
	var targets []struct {
		absPath string
		zipRoot string
	}

	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}
		name := filepath.Base(path)
		parent := filepath.Base(filepath.Dir(path))
		if targetSet[name] {
			zipPath := filepath.Join(parent, name)
			targets = append(targets, struct {
				absPath string
				zipRoot string
			}{absPath: path, zipRoot: zipPath})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("未找到任何匹配的作业目录")
	}

	// 在内存中创建 zip 文件
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, t := range targets {
		fsys := os.DirFS(t.absPath)
		err := addFsToZip(fsys, zipWriter, t.zipRoot)
		//err := addDirToZip(zipWriter, t.absPath, t.zipRoot)
		if err != nil {
			return nil, err
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 将指定目录压缩到 zip 中，保留路径结构
func addDirToZip(zipWriter *zip.Writer, basePath, zipRoot string) error {
	return filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		zipPath := filepath.ToSlash(filepath.Join(zipRoot, relPath))

		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(writer, f)
		return err
	})
}

// 使用 fs.FS 接口实现的压缩
func addFsToZip(fsys fs.FS, zipWriter *zip.Writer, zipRoot string) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		zipPath := filepath.ToSlash(filepath.Join(zipRoot, path))
		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		file, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}
