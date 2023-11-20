package master

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"pandafs/core"
)

type MasterNode struct {
	api     *gin.Engine
	ln      net.Listener
	svr     *grpc.Server
	nodeSvr *core.NodeServiceGrpcServer
	mu      sync.Mutex
	mules   map[string]struct{}
}

func (m *MasterNode) AddMuleID(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.mules[id]; !exists {
		m.mules[id] = struct{}{}
		fmt.Printf("Added Mule ID: %s\n", id)
	}
}

func (n *MasterNode) ReportStatus(ctx context.Context, req *core.Request) (*core.Response, error) {
	muleID := req.GetId()

	n.AddMuleID(muleID)

	return &core.Response{
		Data: "Status received",
	}, nil
}

func (n *MasterNode) Init() (err error) {
	n.ln, err = net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}
	n.svr = grpc.NewServer()
	n.nodeSvr = core.GetNodeServiceGrpcServer()
	n.mules = make(map[string]struct{})
	core.RegisterNodeServiceServer(n.svr, n.nodeSvr)

	n.api = gin.Default()
	n.api.POST("/tasks", func(c *gin.Context) {
		rawData, err := c.GetRawData()
		payload := string(rawData)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		HandleCommand(payload)
		n.nodeSvr.CmdChannel <- payload

		response := <-n.nodeSvr.ResponseChannel
		c.JSON(http.StatusOK, gin.H{"response": response.GetData()})
	})

	n.api.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.String(400, "Bad request: no file provided")
			return
		}
		defer file.Close()

		// Create a new file on the server to store the uploaded file
		dst, err := os.Create("./uploads/" + header.Filename)
		if err != nil {
			c.String(500, "Failed to create the file on server")
			return
		}
		defer dst.Close()

		// Copy the file content from the request to the newly created file
		_, err = io.Copy(dst, file)
		if err != nil {
			c.String(500, "Failed to copy the file content")
			return
		}

		c.String(200, "File uploaded successfully: %s", header.Filename)
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
