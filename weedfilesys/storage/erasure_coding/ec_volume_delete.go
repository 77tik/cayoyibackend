package erasure_coding

import (
	"cayoyibackend/weedfilesys/storage/types"
	"fmt"
	"io"
	"os"
)

// Erasure Coding（EC）卷中删除 needle 的逻辑，其目标是：
// 不直接从 .ecXX 分片中物理删除 needle 内容，而是在 .ecx 索引文件中将其标记为删除（逻辑删除），
// 并记录到 .ecj journal 文件中，以便后续批量重建 .ecx 文件。

// Delete needleId → 修改 .ecx 文件中 size 字段为 Tombstone
//               → 在 .ecj 中记录该 needleId
//
//RebuildEcxFile:
//    遍历 .ecj 文件中所有 needleId
//    → 用 SearchNeedleFromSortedIndex 查找
//    → 再次标记删除
//    → 最后删除 .ecj 文件

var (
	MarkNeedleDeleted = func(file *os.File, offset int64) error {
		b := make([]byte, types.SizeSize)
		types.SizeToBytes(b, types.TombstoneFileSize)
		n, err := file.WriteAt(b, offset+types.NeedleIdSize+types.OffsetSize)
		if err != nil {
			return fmt.Errorf("sorted needle write error: %w", err)
		}
		if n != types.SizeSize {
			return fmt.Errorf("sorted needle written %d bytes, expecting %d", n, types.SizeSize)
		}
		return nil
	}
)

func (ev *EcVolume) DeleteNeedleFromEcx(needleId types.NeedleId) (err error) {

	_, _, err = SearchNeedleFromSortedIndex(ev.ecxFile, ev.ecxFileSize, needleId, MarkNeedleDeleted)

	if err != nil {
		if err == NotFoundError {
			return nil
		}
		return err
	}

	b := make([]byte, types.NeedleIdSize)
	types.NeedleIdToBytes(b, needleId)

	ev.ecjFileAccessLock.Lock()

	ev.ecjFile.Seek(0, io.SeekEnd)
	ev.ecjFile.Write(b)

	ev.ecjFileAccessLock.Unlock()

	return
}

func RebuildEcxFile(baseFileName string) error {

	if !util.FileExists(baseFileName + ".ecj") {
		return nil
	}

	ecxFile, err := os.OpenFile(baseFileName+".ecx", os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("rebuild: failed to open ecx file: %w", err)
	}
	defer ecxFile.Close()

	fstat, err := ecxFile.Stat()
	if err != nil {
		return err
	}

	ecxFileSize := fstat.Size()

	ecjFile, err := os.OpenFile(baseFileName+".ecj", os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("rebuild: failed to open ecj file: %w", err)
	}

	buf := make([]byte, types.NeedleIdSize)
	for {
		n, _ := ecjFile.Read(buf)
		if n != types.NeedleIdSize {
			break
		}

		needleId := types.BytesToNeedleId(buf)

		_, _, err = SearchNeedleFromSortedIndex(ecxFile, ecxFileSize, needleId, MarkNeedleDeleted)

		if err != nil && err != NotFoundError {
			ecxFile.Close()
			return err
		}

	}

	ecxFile.Close()

	os.Remove(baseFileName + ".ecj")

	return nil
}
