package utils

import "net/http"

//解决跨域

func handleInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 允许访问资源服务
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许跨域传递http头部
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		// 返回数据类型
		w.Header().Set("content-type", "application/json")
		// 添加允许请求方法
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE")
		// TODO 异常处理拦截器
		// TODO 日志记录拦截器
		// TODO 有没有更优雅的方式？
		h(w, r)
	}
}