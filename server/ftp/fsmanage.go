package ftp

import (
	"context"
	"fmt"
	"github.com/AlliotTech/openalist/internal/errs"
	"github.com/AlliotTech/openalist/internal/fs"
	"github.com/AlliotTech/openalist/internal/model"
	"github.com/AlliotTech/openalist/internal/op"
	"github.com/AlliotTech/openalist/server/common"
	"github.com/pkg/errors"
	stdpath "path"
)

func Mkdir(ctx context.Context, path string) error {
	user := ctx.Value("user").(*model.User)
	reqPath, err := user.JoinPath(path)
	if err != nil {
		return err
	}
	if !user.CanWrite() || !user.CanFTPManage() {
		meta, err := op.GetNearestMeta(stdpath.Dir(reqPath))
		if err != nil {
			if !errors.Is(errors.Cause(err), errs.MetaNotFound) {
				return err
			}
		}
		if !common.CanWrite(meta, reqPath) {
			return errs.PermissionDenied
		}
	}
	return fs.MakeDir(ctx, reqPath)
}

func Remove(ctx context.Context, path string) error {
	user := ctx.Value("user").(*model.User)
	if !user.CanRemove() || !user.CanFTPManage() {
		return errs.PermissionDenied
	}
	reqPath, err := user.JoinPath(path)
	if err != nil {
		return err
	}
	return fs.Remove(ctx, reqPath)
}

func Rename(ctx context.Context, oldPath, newPath string) error {
	user := ctx.Value("user").(*model.User)
	srcPath, err := user.JoinPath(oldPath)
	if err != nil {
		return err
	}
	dstPath, err := user.JoinPath(newPath)
	if err != nil {
		return err
	}
	srcDir, srcBase := stdpath.Split(srcPath)
	dstDir, dstBase := stdpath.Split(dstPath)
	if srcDir == dstDir {
		if !user.CanRename() || !user.CanFTPManage() {
			return errs.PermissionDenied
		}
		return fs.Rename(ctx, srcPath, dstBase)
	} else {
		if !user.CanFTPManage() || !user.CanMove() || (srcBase != dstBase && !user.CanRename()) {
			return errs.PermissionDenied
		}
		if err = fs.Move(ctx, srcPath, dstDir); err != nil {
			if srcBase != dstBase {
				return err
			}
			if _, err1 := fs.Copy(ctx, srcPath, dstDir); err1 != nil {
				return fmt.Errorf("failed move for %+v, and failed try copying for %+v", err, err1)
			}
			return nil
		}
		if srcBase != dstBase {
			return fs.Rename(ctx, stdpath.Join(dstDir, srcBase), dstBase)
		}
		return nil
	}
}
