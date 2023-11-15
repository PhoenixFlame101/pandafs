package mule

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

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

	fmt.Println("mule started")

	_, _ = n.c.ReportStatus(context.Background(), &core.Request{})

	stream, _ := n.c.AssignTask(context.Background(), &core.Request{})
	for {

		res, err := stream.Recv()
		if err != nil {
			return
		}

		fmt.Println("received command: ", res.Data)

		parts := strings.Split(res.Data, " ")
		if err := exec.Command(parts[0], parts[1:]...).Run(); err != nil {
			fmt.Println(err)
		}
	}
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
