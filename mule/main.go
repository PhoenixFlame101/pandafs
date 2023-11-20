// mule/main.go
package mule

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"google.golang.org/grpc"

	"pandafs/core"
)

type Mule struct {
	conn *grpc.ClientConn
	c    core.NodeServiceClient
}

func (n *Mule) Init() (err error) {
	n.conn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return err
	}
	n.c = core.NewNodeServiceClient(n.conn)
	return nil
}

func (n *Mule) Start() {
	fmt.Println("Mule started")

	// Periodically ping the master every 10 seconds
	go func() {
		for {
			n.pingMaster()
			time.Sleep(10 * time.Second)
		}
	}()

	// Handle incoming commands from the master
	stream, _ := n.c.AssignTask(context.Background(), &core.Request{})
	for {
		res, err := stream.Recv()
		if err != nil {
			return
		}

		fmt.Println("Received command:", res.GetData())
		parts := strings.Split(res.GetData(), " ")
		if err := exec.Command(parts[0], parts[1:]...).Run(); err != nil {
			fmt.Println(err)
		}
	}
}

func (n *Mule) pingMaster() {
	response, err := n.c.ReportStatus(context.Background(), &core.Request{})
	if err != nil {
		fmt.Println("Error pinging master:", err)
		return
	}
	fmt.Println("Ping response from master:", response.GetData())
}

var mule *Mule

func GetMule() *Mule {
	if mule == nil {
		mule = &Mule{}
		if err := mule.Init(); err != nil {
			panic(err)
		}
	}
	return mule
}

func main() {
	GetMule().Start()
}
