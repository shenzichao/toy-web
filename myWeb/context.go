package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// Context 选择 定义为接口还是结构体
type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func (c *Context) ReadJson(req interface{}) error {
	r := c.R
	body, err := io.ReadAll(r.Body)
	if err != nil {
		//fmt.Fprintf(c.w, "get the body failed, error: %v", err)
		return err
	}
	err = json.Unmarshal(body, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) WriteJson(code int, resp interface{}) error {
	c.W.WriteHeader(code) // 写入HTTP状态码
	respJson, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = c.W.Write(respJson)
	if err != nil {
		return err
	}
	return nil

}

// 对WriteJson进行再次封装

func (c *Context) OkJson(resp interface{}) error {
	// 委托
	return c.WriteJson(http.StatusOK, resp)
}

func (c *Context) SystemErrorJson(resp interface{}) error {
	return c.WriteJson(http.StatusInternalServerError, resp)
}

func (c *Context) BadRequestJson(resp interface{}) error {
	return c.WriteJson(http.StatusBadRequest, resp)
}

func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		W: writer,
		R: request,
	}
}
