package logic

import (
	"book/service/user/rpc/types/user"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"book/service/search/api/internal/svc"
	"book/service/search/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchReq) (*types.SearchReply, error) {
	fmt.Println("jjjjjjjjj")
	userIdNumber := json.Number(fmt.Sprintf("%v", l.ctx.Value("userID")))
	logx.Infof("userId: %s", userIdNumber)
	userId, err := userIdNumber.Int64()
	if err != nil {
		return nil, err
	}
	fmt.Println("冯绍峰奶粉服务器呢法去年范围内法球呢")
	// 使用user rpc
	one, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user.IdReq{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("woshunubabba")
	searchBook, err := l.svcCtx.SearchModel.FindOne(l.ctx, 1)
	if err != nil {
		return nil, err
	}
	return &types.SearchReply{
		Id:      strconv.FormatInt(searchBook.Id, 10),
		Name:    searchBook.Name,
		Message: fmt.Sprintf("name为%v查询了%v", one.Name, searchBook.Name),
	}, nil
}
