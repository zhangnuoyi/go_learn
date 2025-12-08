package router

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"strings"

	"gateway/config"
)

// Router 路由管理器
type Router struct {
	config *config.Config
	// 存储上游服务的代理实例
	proxies map[string]*httputil.ReverseProxy
}

// NewRouter 创建路由管理器实例
func NewRouter(config *config.Config) *Router {
	router := &Router{
		config:  config,
		proxies: make(map[string]*httputil.ReverseProxy),
	}

	// 初始化所有上游服务的代理
	for _, upstream := range config.Upstreams {
		targetURL, err := url.Parse(fmt.Sprintf("http://%s", upstream.Http.Target))
		if err != nil {
			continue
		}

		// 创建反向代理
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		router.proxies[upstream.Name] = proxy
	}

	return router
}

// MatchRoute 匹配路由，返回匹配的上游服务名称
func (r *Router) MatchRoute(method, path string) string {
	for _, upstream := range r.config.Upstreams {
		for _, mapping := range upstream.Mappings {
			if mapping.Method == method && r.matchPath(mapping.Path, path) {
				return upstream.Name
			}
		}
	}
	return ""
}

// matchPath 匹配路径，支持精确匹配、路径参数和通配符
func (r *Router) matchPath(pattern, path string) bool {
	// 精确匹配
	if pattern == path {
		return true
	}

	// 通配符匹配
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		if strings.HasPrefix(path, prefix) {
			// 确保是完整路径段匹配
			if len(prefix) == 0 {
				return true
			}
			if len(path) == len(prefix) {
				return true
			}
			if path[len(prefix)] == '/' {
				return true
			}
		}
		return false
	}

	// 路径参数匹配（如 /v1/post/:id）
	if strings.Contains(pattern, ":") {
		patternParts := strings.Split(pattern, "/")
		pathParts := strings.Split(path, "/")

		if len(patternParts) != len(pathParts) {
			return false
		}

		for i, patternPart := range patternParts {
			if strings.HasPrefix(patternPart, ":") {
				continue // 跳过参数部分
			}
			if patternPart != pathParts[i] {
				return false
			}
		}
		return true
	}

	return false
}

// GetProxy 获取上游服务的代理实例
func (r *Router) GetProxy(upstreamName string) *httputil.ReverseProxy {
	return r.proxies[upstreamName]
}
