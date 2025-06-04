package znet

import (
	"log/slog"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	routes         map[uint32]ziface.IRouter
	workerPoolSize uint32
	taskQueue      []chan ziface.IRequest
}

func NewMsgHandler() ziface.IMsgHandler {
	return &MsgHandler{
		routes:         make(map[uint32]ziface.IRouter),
		workerPoolSize: utils.GlobalObject.WorkerPoolSize,
		taskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

func (mh *MsgHandler) CreatWorkerPool() {
	for i := range mh.taskQueue {
		mh.taskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskQueueSize)
		go mh.StartWorker(i, mh.taskQueue[i])
	}
}

func (mh *MsgHandler) StartWorker(workerID int, workerQueue chan ziface.IRequest) {
	slog.Debug("Worker", "workerID", workerID, "msg", "worker start")
	for {
		select {
		case req := <-workerQueue:
			mh.DoMsgHandler(req)
		}
	}
}

func (mh *MsgHandler) SendReqToWorker(req ziface.IRequest) {
	workerID := uint32(req.GetConnection().GetConnID()) % mh.workerPoolSize
	slog.Debug("SendMsgToWorker", "ConnID", req.GetConnection().GetConnID(), "workerID", workerID, "msg", "send req to worker")
	mh.taskQueue[workerID] <- req
}
