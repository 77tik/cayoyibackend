package needle

import (
	"cayoyibackend/weedfilesys/glog"
	"cayoyibackend/weedfilesys/storage/backend"
	. "cayoyibackend/weedfilesys/storage/types"
	"cayoyibackend/weedfilesys/util"
	"cayoyibackend/weedfilesys/util/buffer_pool"
	"fmt"
)

// 将needle追加到磁盘中
// 1. 获取当前文件尾部偏移 end = w.GetStat()
// 2. 如果超出 MaxPossibleVolumeSize 且有 Data，报错返回
// 3. 从 buffer_pool 拿一个 bytes.Buffer（避免频繁分配内存）
// 4. 调用 writeNeedleByVersion(...) 将 needle 编码到 buffer 中
// 5. 写入 buffer 到 Volume 文件 offset 位置
// 6. 出错时 truncate 回原 offset，确保原子性
func (n *Needle) Append(w backend.BackendStorageFile, version Version) (offset uint64, size Size, actualSize int64, err error) {
	end, _, e := w.GetStat()
	if e != nil {
		err = fmt.Errorf("Cannot Read Current Volume Position: %w", e)
		return
	}
	offset = uint64(end)
	if offset >= MaxPossibleVolumeSize && len(n.Data) != 0 {
		err = fmt.Errorf("Volume Size %d Exceeded %d", offset, MaxPossibleVolumeSize)
		return
	}
	bytesBuffer := buffer_pool.SyncPoolGetBuffer()
	defer func() {
		if err != nil {
			if te := w.Truncate(end); te != nil {
				// handle error or log
			}
		}
		buffer_pool.SyncPoolPutBuffer(bytesBuffer)
	}()

	size, actualSize, err = writeNeedleByVersion(version, n, offset, bytesBuffer)
	if err != nil {
		return
	}

	_, err = w.WriteAt(bytesBuffer.Bytes(), int64(offset))
	if err != nil {
		err = fmt.Errorf("failed to write %d bytes to %s at offset %d: %w", actualSize, w.Name(), offset, err)
	}

	return offset, size, actualSize, err
}

// 1. 获取文件当前末尾 offset
// 2. 如果是 Version3，写入 appendAtNs 时间戳（尾部 ts 字段）
// 3. 调用 WriteAt(dataSlice, offset) 写入文件
// 4. 写失败则 truncate 回 offset
func WriteNeedleBlob(w backend.BackendStorageFile, dataSlice []byte, size Size, appendAtNs uint64, version Version) (offset uint64, err error) {

	if end, _, e := w.GetStat(); e == nil {
		defer func(w backend.BackendStorageFile, off int64) {
			if err != nil {
				if te := w.Truncate(end); te != nil {
					glog.V(0).Infof("Failed to truncate %s back to %d with error: %v", w.Name(), end, te)
				}
			}
		}(w, end)
		offset = uint64(end)
	} else {
		err = fmt.Errorf("Cannot Read Current Volume Position: %v", e)
		return
	}

	if version == Version3 {
		tsOffset := NeedleHeaderSize + size + NeedleChecksumSize
		util.Uint64toBytes(dataSlice[tsOffset:tsOffset+TimestampSize], appendAtNs)
	}

	if err == nil {
		_, err = w.WriteAt(dataSlice, int64(offset))
	}

	return

}
