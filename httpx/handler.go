package httpx

import (
	"github.com/soyacen/gox/slicex"
	"net/http"
)

var Default405Body = []byte("405 method not allowed")

func GetHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

func GetHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

func HeadHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodHead {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

func HeadHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodHead {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

func PostHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

func PostHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

func PutHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

func PutHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

func PatchHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPatch {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

func PatchHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPatch {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

func DeleteHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

func DeleteHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

func ConnectHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodConnect {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

func ConnectHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodConnect {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

func OptionsHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodOptions {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

func OptionsHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodOptions {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

func TraceHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodTrace {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

func TraceHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodTrace {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

func Handler(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

func HandlerFunc(method string, handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

func MatchHandler(methods []string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if slicex.NotContains(methods, req.Method) {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

func MatchHandlerFunc(methods []string, handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if slicex.NotContains(methods, req.Method) {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}
