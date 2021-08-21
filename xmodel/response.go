
package xmodel

type Response struct {
    Code    uint32      `json:"code"`
    CodeMsg string      `json:"codeMsg"`
    Data    interface{} `json:"data"`
}

type Header struct {
    Seq         string      `json:"seq"`
    Cmd         string      `json:"cmd"`
    Response   *Response    `json:"response"`
}

func NewResponseHeader(seq string, cmd string, code uint32, codeMsg string,data interface{}) (*Header) {
    response := &Response{ Code: code, CodeMsg: codeMsg, Data: data }
    return &Header{Seq: seq, Cmd: cmd, Response: response}
}
