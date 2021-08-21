
package xterminus

import (
    "net"
    "context"
    "log"
    "errors"

    "github.com/spf13/viper"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"

    "chat.com/xcode"
    "chat.com/xmodel"
    "chat.com/xproto"
)

type GrpcS struct {
}

func StartGrpcServer() {
    rpcPort := viper.GetString("rpc.port")

    l, err := net.Listen("tcp", ":" + rpcPort)
    if err != nil {
        log.Fatalf("-FATA- listen failed: %v", err)
    }

    s := grpc.NewServer()
    xproto.RegisterGreeterServer(s, &GrpcS{})
    // Register reflection service on gRPC server.
    reflection.Register(s)
    if err :=  s.Serve(l); err != nil {
        log.Fatalf("-FATA- serve failed: %v", err)
    }
}

func (this *GrpcS) QueryUsersOnline(ctx context.Context, req *xproto.QueryUsersOnlineReq) (*xproto.QueryUsersOnlineRsp, error) {
    isOnline := CheckUserOnline(req.GetAppId(), req.GetUserId())

    return &xproto.QueryUsersOnlineRsp{ RetCode: xcode.OK, Online: isOnline }, nil
}

func (this *GrpcS) SendMsg(ctx context.Context, req *xproto.SendMsgReq) (*xproto.SendMsgRsp, error) {
    if req.GetIsLocal() {
        return &xproto.SendMsgRsp{ RetCode: xcode.ParameterIllegal }, errors.New("unsupport local")
    }

    data := xmodel.GetMsgData(req.GetUserId(), req.GetSeq(), req.GetCms(), req.GetMsg())
    bl, err := SendMessageToLocal(req.GetAppId(), req.GetUserId(), data)
    if err != nil {
        return &xproto.SendMsgRsp{ RetCode: xcode.ServerError }, errors.New("server error")
    }

    if ! bl {
        return &xproto.SendMsgRsp{ RetCode: xcode.OperationFailure }, errors.New("operation failure")
    }

    return &xproto.SendMsgRsp{ RetCode: xcode.OK }, nil
}

func (this *GrpcS) SendMsgAll(ctx context.Context, req *xproto.SendMsgAllReq) (*xproto.SendMsgAllRsp, error) {
    data := xmodel.GetMsgData(req.GetUserId(), req.GetSeq(), req.GetCms(), req.GetMsg())
    sendMessageToAppId(req.GetAppId(), req.GetUserId(), data)

    return &xproto.SendMsgAllRsp{ RetCode: xcode.OK }, nil
}

func (this *GrpcS) GetUserList(ctx context.Context, req *xproto.GetUserListReq) (*xproto.GetUserListRsp, error) {
    appId := req.GetAppId()
    userList := getUserList(appId)

    return &xproto.GetUserListRsp{ RetCode: xcode.OK, UserId: userList }, nil
}