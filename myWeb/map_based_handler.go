package main

import "net/http"

type HandlerBasedOnMap struct {
	// key method + "#" + url
	handlers map[string]func(ctx *Context)
}

func (h *HandlerBasedOnMap) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := h.key(request.Method, request.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		ctx := NewContext(writer, request)
		handler(ctx)
	} else {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("no router match"))
	}
}

func (h *HandlerBasedOnMap) key(method, pattern string) string {
	return method + "#" + pattern
}
