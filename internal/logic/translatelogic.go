// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/mustang5910/deeplx-translategemma/internal/svc"
	"github.com/mustang5910/deeplx-translategemma/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TranslateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTranslateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TranslateLogic {
	return &TranslateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TranslateLogic) Translate(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
