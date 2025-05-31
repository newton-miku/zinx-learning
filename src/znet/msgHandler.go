package znet

import (
	"log/slog"
	"zinx/ziface"
)

type MsgHandler struct {
	routes map[uint32]ziface.IRouter
}

func NewMsgHandler() ziface.IMsgHandler {
	return &MsgHandler{
		routes: make(map[uint32]ziface.IRouter),
	}
}
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	if handler, ok := mh.routes[request.GetMsgID()]; ok {
		handler.PreHandle(request)
		handler.Handle(request)
		handler.PostHandle(request)
	} else {
		slog.Error("Router", "msgID", request.GetMsgID(), "msg", "router not found")
	}
}
func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.routes[msgID]; ok {
		slog.Error("Router", "msgID", msgID, "msg", "router already exists")
		return
	} else {
		mh.routes[msgID] = router
		slog.Debug("Router", "msgID", msgID, "msg", "add router successfully")
	}
}
