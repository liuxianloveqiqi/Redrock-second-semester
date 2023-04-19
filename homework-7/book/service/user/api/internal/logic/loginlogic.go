package logic

import (
	"book/service/user/api/internal/svc"
	"book/service/user/api/internal/types"
	"book/service/user/model"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginReply, error) {
	if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		return nil, errors.New("参数错误")
	}

	userInfo, err := l.svcCtx.UserModel.FindOneByNumber(l.ctx, req.Username)
	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errors.New("用户名不存在")
	default:
		return nil, err
	}

	if userInfo.Password != req.Password {
		return nil, errors.New("用户密码不正确")
	}
	fmt.Println(9999999999999999)
	// ---start---
	accessTokenString, refreshTokenString := l.GetToken(userInfo.Id, uuid.New().String())
	if accessTokenString == "" || refreshTokenString == "" {
		return nil, errors.New("生成jwt错误")
	}

	// ---end---
	return &types.LoginReply{
		Id:           userInfo.Id,
		Name:         userInfo.Name,
		Gender:       userInfo.Gender,
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
