
package xmodel

import (
    "errors"
    "fmt"
    "strings"
)

type Server struct {
    Ip		string	`json:"ip"`
    Port    string  `json:"port"`
}

func NewServer(ip, port string) *Server {
    return &Server{Ip: ip, Port: port}
}

func (s *Server) String() (string) {
    if s == nil {
        return ""
    }
    return fmt.Sprintf("%s:%s", s.Ip, s.Port)
}

func ToServer(s string) (*Server, error) {
    list := strings.Split(s, ":")
    if len(list) != 2 {
        return nil, errors.New(fmt.Sprintf("ToServer: illegal format [%s]", s))
    }

    return &Server{ Ip: list[0], Port: list[1] }, nil
}