// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"06-blog-cloud/blog_post_api/api/internal/logic"
	"06-blog-cloud/blog_post_api/api/internal/svc"
	"06-blog-cloud/blog_post_api/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func saveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostDto
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSaveLogic(r.Context(), svcCtx)
		resp, err := l.Save(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
