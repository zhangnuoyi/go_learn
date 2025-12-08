// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"06-blog-cloud/blog_post_api/api/internal/logic"
	"06-blog-cloud/blog_post_api/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func listHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewListLogic(r.Context(), svcCtx)
		resp, err := l.List()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
