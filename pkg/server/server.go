package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

type HandlerFunc func(conn net.Conn)

type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}

func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go s.handle(conn)
	}
	return nil
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Print(err)
		}
	}()

	s.mu.RLock()
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err == io.EOF {
		log.Printf("%s", buf[:n])
	}
	if err != nil {
		log.Print(err)
	}

	data := buf[:n]
	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)
	if requestLineEnd == -1 {

	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {

	}

	_, path, _ := parts[0], parts[1], parts[2]
	handler := s.handlers[path]
	s.mu.RUnlock()
	handler(conn)
}

func (s *Server) PrintResponse(conn net.Conn, response string) (err error) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()

	_, err = conn.Write([]byte(
		"HTTP/1.1 200 OK\r\n" +
			"Content-Length: " + strconv.Itoa(len(response)) + "\r\n" +
			"Content-Type: text/html\r\n" +
			"Connection: close\r\n" +
			"\r\n" +
			response,
	))
	if err != nil {
		return err
	}

	return nil
}

func parseRequest(conn net.Conn) (path string) {
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err == io.EOF {
		log.Printf("%s", buf[:n])
	}
	if err != nil {
		log.Print(err)
	}

	data := buf[:n]
	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)
	if requestLineEnd == -1 {

	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {

	}

	method, path, version := parts[0], parts[1], parts[2]
	if method != "GET" {

	}

	if version != "HTTP/1.1" {

	}

	return path
}
