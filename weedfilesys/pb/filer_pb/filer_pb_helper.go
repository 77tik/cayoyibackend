package filer_pb

import (
	"context"
	"errors"
	"fmt"
	"my_backend/weedfilesys/glog"
	"strings"
)

var ErrNotFound = errors.New("filer: no entry is found in filer store")

func LookupEntry(ctx context.Context, client WeedfilesysFilerClient, request *LookupDirectoryEntryRequest) (*LookupDirectoryEntryResponse, error) {
	resp, err := client.LookupDirectoryEntry(ctx, request)
	if err != nil {
		if err == ErrNotFound || strings.Contains(err.Error(), ErrNotFound.Error()) {
			return nil, ErrNotFound
		}
		glog.V(3).InfofCtx(ctx, "read %s/%v: %v", request.Directory, request.Name, err)
		return nil, fmt.Errorf("LookupEntry1:%w", err)
	}
	if resp.Entry == nil {
		return nil, ErrNotFound
	}
	return resp, nil
}

func UpdateEntry(ctx context.Context, client WeedfilesysFilerClient, request *UpdateEntryRequest) error {
	_, err := client.UpdateEntry(ctx, request)
	if err != nil {
		glog.V(1).InfofCtx(ctx, "update entry %s/%s :%v", request.Directory, request.Entry.Name, err)
		return fmt.Errorf("UpdateEntry: %w", err)
	}
	return nil
}

func CreateEntry(ctx context.Context, client WeedfilesysFilerClient, request *CreateEntryRequest) error {
	resp, err := client.CreateEntry(ctx, request)
	if err != nil {
		glog.V(1).InfofCtx(ctx, "create entry %s/%s %v: %v", request.Directory, request.Entry.Name, request.OExcl, err)
		return fmt.Errorf("CreateEntry: %w", err)
	}
	if resp.Error != "" {
		glog.V(1).InfofCtx(ctx, "create entry %s/%s %v: %v", request.Directory, request.Entry.Name, request.OExcl, resp.Error)
		return fmt.Errorf("CreateEntry : %v", resp.Error)
	}
	return nil
}
