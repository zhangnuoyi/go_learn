// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"api/internal/logic/user"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterDto
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		logic := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := logic.Register(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// 确保响应包含id字段
			if resp != nil {
				httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
					"id":    resp.Id,
					"name":  resp.Name,
					"email": resp.Email,
				})
			} else {
				httpx.OkJsonCtx(r.Context(), w, resp)
			}
		}
	}
}
