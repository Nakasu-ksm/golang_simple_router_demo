package simplerouter

import (
	"GO_project/router/Model"
	"encoding/json"
	"net/http"
)

// Custom
type HttpResponse struct {
	writer http.ResponseWriter
}
type HttpReceiveFunction func(response HttpResponse, r *http.Request)

type HttpReceive struct {
	handle HttpReceiveFunction
	method string
	s      *srouter
}

type srouter struct {
	server              *http.ServeMux
	customerMethodError HttpReceiveFunction
}

//type srouter http.ServeMux

func (s *srouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ser := s.server
	ser.ServeHTTP(writer, request)
}

// New initial
func New() *srouter {
	router := srouter{}
	nsm := http.NewServeMux()
	router.server = nsm
	return &router
}

func (s *srouter) GET(path string, function HttpReceiveFunction) {
	s.setHandle(path, function, "GET")
}
func (s *srouter) POST(path string, function HttpReceiveFunction) {
	s.setHandle(path, function, "POST")
}
func (s *srouter) PUT(path string, function HttpReceiveFunction) {
	s.setHandle(path, function, "PUT")
}
func (s *srouter) PATCH(path string, function HttpReceiveFunction) {
	s.setHandle(path, function, "PATCH")
}
func (s *srouter) DELETE(path string, function HttpReceiveFunction) {
	s.setHandle(path, function, "DELETE")
}

// function set
func (s *srouter) setHandle(path string, function HttpReceiveFunction, method string) {
	s.server.Handle(path, extendHttp(function, method, s))
}

func extendHttp(function HttpReceiveFunction, method string, s *srouter) http.Handler {
	return &HttpReceive{handle: function, method: method, s: s}
}

func (h *HttpReceive) notFound(writer HttpResponse, request *http.Request) {
	if h.s.customerMethodError == nil {
		http.Error(writer.writer, "ご指定のページまたはファイルが見つかりませんでした", http.StatusNotFound)
		return
	}
	h.s.customerMethodError(writer, request)
}

func (h *HttpReceive) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	rs := HttpResponse{}
	rs.writer = writer
	if h.method != request.Method {
		h.notFound(rs, request)
		return
	}
	h.handle(rs, request)
}

func (h HttpResponse) ReturnJson(code int, message string) {
	marshal, err := json.Marshal(Model.Message{Code: code, Message: message})
	if err == nil {
		_, err := h.writer.Write(marshal)
		if err != nil {
			panic(err)
		}
		return
	}
	h.writer.Write([]byte("Error"))
	return
}

func (h HttpResponse) Write(mByte []byte) {
	_, err := h.writer.Write(mByte)
	if err != nil {
		panic(err)
	}
	return
}

func (s *srouter) SetMethodError(function HttpReceiveFunction) {
	s.customerMethodError = function
}
