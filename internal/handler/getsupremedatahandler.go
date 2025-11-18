package handler

import (
	"net/http"
	"time"

	"imd-seat-be/internal/logic"
	"imd-seat-be/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getSupremeDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetSupremeDataLogic(r.Context(), svcCtx)
		resp, err := l.GetSupremeData()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func getSupremeListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetSupremeDataLogic(r.Context(), svcCtx)
		resp, err := l.GetSupremeList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func download(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		l := logic.NewGetSupremeDataLogic(r.Context(), svcCtx)
		file, err := l.Download(filename)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// 设置下载响应头
			w.Header().Set("Content-Disposition", "attachment; filename="+filename)
			w.Header().Set("Content-Type", "application/octet-stream")

			// 将文件内容写入响应
			http.ServeContent(w, r, filename, time.Now(), file)
		}
		defer file.Close()
	}
}
