
package xmodel

// request message
type Request struct {
    Seq   string        `json:"seq"`            // unique sequence for message
    Cmd   string        `json:"cmd"`            // command word
    Data  interface{}   `json:"data,omitempty"` // payload
}

// login request
type Login struct {
    ServiceToken    string  `json:"serviceToken"`      // verify user whether logined or not
    AppId           uint32  `json:"appId,omitempty"`
    UserId          string  `json:"userId,omitempty"`
}

// heartbeat request
type HeartBeat struct {
    UserId string `json:"userId,omitempty"`
}