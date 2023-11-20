// master/main.go
package master

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"pandafs/core"
)

type MasterNode struct {
	api     *gin.Engine
	ln      net.Listener
	svr     *grpc.Server
	nodeSvr *core.NodeServiceGrpcServer
}

func (n *MasterNode) Init() (err error) {
	n.ln, err = net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}
	n.svr = grpc.NewServer()
	n.nodeSvr = core.GetNodeServiceGrpcServer()
	core.RegisterNodeServiceServer(n.svr, n.nodeSvr)

	n.api = gin.Default()
	n.api.POST("/tasks", func(c *gin.Context) {
		var payload struct {
			Cmd string `json:"cmd"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		n.nodeSvr.CmdChannel <- payload.Cmd

		response := <-n.nodeSvr.ResponseChannel
		c.JSON(http.StatusOK, gin.H{"response": response.GetData()})
	})

	return nil
}

func (n *MasterNode) Start() {
	go n.svr.Serve(n.ln)
	_ = n.api.Run(":9092")
	n.svr.Stop()
}

var node *MasterNode

func GetMasterNode() *MasterNode {
	if node == nil {
		node = &MasterNode{}
		if err := node.Init(); err != nil {
			panic(err)
		}
	}
	return node
}

func main() {
	GetMasterNode().Start()
}
