package mule

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc"

	"pandafs/core"
)

type Mule struct {
	id   string
	conn *grpc.ClientConn
	c    core.NodeServiceClient
}

func generateRandomID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a byte slice to store the generated ID
	id := make([]byte, length)

	// Generate random ID characters
	for i := range id {
		id[i] = charset[rand.Intn(len(charset))]
	}

	return string(id)
}

func (n *Mule) Init() (err error) {
	n.id = generateRandomID(16)
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
		// parts := strings.Split(res.GetData(), " ")
		// if err := exec.Command(parts[0], parts[1:]...).Run(); err != nil {
		// 	fmt.Println(err)
		// }
	}
}

func (n *Mule) pingMaster() {
	response, err := n.c.ReportStatus(context.Background(), &core.Request{Id: mule.id})
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
