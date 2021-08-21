
package xmodel

import (
    "encoding/json"
    
    "chat.com/xcode"
)

const (
    MessageTypeText = "text"
    
    MessageCmdMsg   = "msg"
    MessageCmdEnter = "enter"
    MessageCmdExit  = "exit"
)

type Message struct {
    Target  string `json:"target"`
    Type    string `json:"type"`        // text/img/...
    Msg     string `json:"msg"`
    From    string `json:"from"`        // who send this 
}

func NewTestMsg(from string, msg string) (*Message) {
    return &Message{
        Type:   MessageTypeText,
        From:   from,
        Msg:    msg,
    }
}

func _getTextMsgData(cmd, userId, msgId, message string) (string) {
    textMsg := &Message{
        Type: MessageTypeText,
        From: userId,
        Msg:  message,
    }
    header := NewResponseHeader(msgId, cmd, xcode.OK, "OK", textMsg)
    b, _ := json.Marshal(header)

    return string(b)
}

func GetMsgData(userId, msgId, cmd, message string) (string) {
    return _getTextMsgData(cmd, userId, msgId, message)
}

func GetTextMsgData(userId, msgId, message string) (string) {
    return _getTextMsgData("msg", userId, msgId, message)
}

func GetTextMsgDataEnter(userId, msgId, message string) (string) {
    return _getTextMsgData("enter", userId, msgId, message)
}

func GetTextMsgDataExit(userId, msgId, message string) (string) {
    return _getTextMsgData("exit", userId, msgId, message)
}
