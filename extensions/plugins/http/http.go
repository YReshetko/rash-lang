package main

import (
	"errors"
	"fmt"
	"net/http"
)

var HttpPlugin = httpPlugin{
	servers: map[string]*server{},
}

const (
	pkg  = "http"
	ver  = "0.0.1"
	desc = "provides http server functions"
)

type httpPlugin struct {
	servers map[string]*server
}
type Callback func(args ...interface{}) ([]interface{}, error)

type server struct {
	mux  *http.ServeMux
	port string
	// map[path][method]handler
	routes map[string]map[string]Callback
}

func (s httpPlugin) Eval(fnName string, args ...interface{}) ([]interface{}, error) {
	switch fnName {
	case "new":
		return s.newServer(args...)
	case "start":
		return s.start(args...)
	default:
		return nil, fmt.Errorf("function %s not found in %s extension", fnName, pkg)
	}
}

func (s httpPlugin) Call(fnName string, callback func(args ...interface{}) ([]interface{}, error), args ...interface{}) ([]interface{}, error) {
	switch fnName {
	case "register":
		return s.register(callback, args...)
	default:
		return nil, fmt.Errorf("callback function %s not found in %s extension", fnName, pkg)
	}
}

func (s httpPlugin) Package() string {
	return pkg
}

func (s httpPlugin) Version() string {
	return ver
}

func (s httpPlugin) Description() string {
	return desc
}

func (s httpPlugin) newServer(args ...interface{}) ([]interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("expected at least port")
	}

	port := args[0].(string)
	serverName := fmt.Sprintf("http_server_%d", len(s.servers)+1)
	mux := http.NewServeMux()

	s.servers[serverName] = &server{
		mux:    mux,
		port:   "localhost:" + port,
		routes: map[string]map[string]Callback{},
	}

	return []interface{}{serverName}, nil
}

func (s httpPlugin) register(callback func(args ...interface{}) ([]interface{}, error), args ...interface{}) ([]interface{}, error) {
	if len(args) < 3 {
		return nil, errors.New("expected server name, http method and path pattern")
	}
	serverName := args[0].(string)
	method := args[1].(string)
	pattern := args[2].(string)

	m, ok := s.servers[serverName].routes[pattern]
	if !ok {
		m = map[string]Callback{}
		s.servers[serverName].routes[pattern] = m
	}

	m[method] = callback

	return nil, nil
}

func (s httpPlugin) start(args ...interface{}) ([]interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("expected at least server name")
	}

	serverName := args[0].(string)
	server, ok := s.servers[serverName]
	if !ok {
		return nil, errors.New("server not found")
	}

	for route, value := range server.routes {
		v := value
		server.mux.HandleFunc(route, func(writer http.ResponseWriter, request *http.Request) {
			callback, ok := v[request.Method]
			if !ok {
				writer.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			values, err := callback()
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			writer.WriteHeader(http.StatusOK)
			if len(values) > 0 {
				if v, ok := values[0].(string); ok {
					writer.Write([]byte(v))
				}
			}

		})
	}

	go func() {
		err := http.ListenAndServe(server.port, server.mux)
		fmt.Println(err)
	}()

	return nil, nil
}
