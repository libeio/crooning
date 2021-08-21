
package xterminus

import (
    "time"
    "context"
    "errors"
    "log"
    "fmt"

    "google.golang.org/grpc"

    "chat.com/xcode"
    "chat.com/xmodel"
    "chat.com/xproto"
)

func grpcSendMessageToAll(server *xmodel.Server, seq string, appId uint32, userId string, cmd string, message string) (string, error) {
    conn, err := grpc.Dial(server.String(), grpc.WithTimeout(5 * time.Second), grpc.WithBlock(), grpc.WithInsecure())
    if err != nil {
        return "", err
    }
    defer conn.Close()

    c := xproto.NewGreeterClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    req := xproto.SendMsgAllReq{   // protobuf
        Seq:    seq,
        AppId:  appId,
        UserId: userId,
        Cms:    cmd,
        Msg:    message,
    }
    rsp, err := c.SendMsgAll(ctx, &req)
    if err != nil {
        return "", err
    }

    if rsp.GetRetCode() != xcode.OK {
        return "", errors.New(fmt.Sprintf("Unexpected return code:%d", rsp.GetRetCode()))
    }

    return rsp.GetSendMsgId(), nil
}

func grpcGetUserList(server *xmodel.Server, appId uint32) ([]string, error) {
    userIds := make([]string, 0)

    conn, err := grpc.Dial(server.String(), grpc.WithTimeout(5 * time.Second), grpc.WithBlock(), grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    c := xproto.NewGreeterClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    req := xproto.GetUserListReq{
        AppId: appId,
    }
    rsp, err := c.GetUserList(ctx, &req)
    if err != nil {
        return nil, err
    }

    if rsp.GetRetCode() != xcode.OK {
        return nil, errors.New(fmt.Sprintf("Unexpected return code:%d", rsp.GetRetCode()))
    }

    userIds = rsp.GetUserId()

    return userIds, nil
}

func grpcSendMsg(server *xmodel.Server, seq string, appId uint32, userId string, cmd string, msgType string, message string) (string, error) {
    conn, err := grpc.Dial(server.String(), grpc.WithTimeout(5 * time.Second), grpc.WithBlock(), grpc.WithInsecure())
    if err != nil {
        return "", err
    }
    defer conn.Close()

    c := xproto.NewGreeterClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    req := xproto.SendMsgReq{
        Seq:        seq,
        AppId:      appId,
        UserId:     userId,
        Cms:        cmd,
        Type:       msgType,
        Msg:        message,
        IsLocal:    false,
    }
    rsp, err := c.SendMsg(ctx, &req)
    if err != nil {
        return "", err
    }

    if rsp.GetRetCode() != xcode.OK {
        return "", errors.New(fmt.Sprintf("Unexpected return code:%d", rsp.GetRetCode()))
    }

    log.Printf("-INFO- grpc server send to other success")

    return rsp.GetSendMsgId(), nil
}