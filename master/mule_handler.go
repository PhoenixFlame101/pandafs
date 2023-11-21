package master

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	pb "pandafs/core"

	"google.golang.org/grpc"
)

type MuleNode struct {
	ID   string
	Addr string
	// You may need to include additional fields based on your requirements
}

// selectMuleNode selects a mule node for the given chunk.
// This is a simple example, and you may want to implement your own logic.
func selectMuleNode(muleNodes []MuleNode) MuleNode {
	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Generate a random index within the range of the muleNodes slice
	randomIndex := rand.Intn(len(muleNodes))

	// Return the randomly selected mule node
	return muleNodes[randomIndex]
}

// sendChunkToMuleNode sends a chunk to the specified mule node using gRPC.
func sendChunkToMuleNode(chunk []byte, muleNode MuleNode) error {
	conn, err := grpc.Dial(muleNode.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewNodeServiceClient(conn)

	// Replace "YourGRPCMethod" with the actual gRPC method you have in your service
	response, err := client.AssignTask(context.Background(), &pb.Request{Action: string(chunk)})
	if err != nil {
		return err
	}

	fmt.Println("Response from mule node:", response)

	return nil
}

// updateDatabaseWithChunkInfo updates the SQLite database with chunk information.
func updateDatabaseWithChunkInfo(chunkID int, muleNode MuleNode) error {
	result, err := db.Exec(`
      INSERT INTO inode (filename, filesize, isDirectory) VALUES (?, 0, false);
   `, chunkID)
	if err != nil {
		log.Fatal(err)
	}

	inodeID, _ := result.LastInsertId()

	id, err := GetDirID("/")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
      INSERT INTO directory (file_id, dir_id) VALUES (?, ?);
   `, inodeID, id)
	if err != nil {
		log.Fatal(err)
	}
}

func SplitAndDistribute(fileContent []byte, muleNodes []MuleNode) error {
	chunkSize := 10 * 1024 // 10KB

	// Split the file content into chunks
	for i := 0; i < len(fileContent); i += chunkSize {
		end := i + chunkSize
		if end > len(fileContent) {
			end = len(fileContent)
		}
		chunk := fileContent[i:end]

		// Send the chunk to a mule node using gRPC
		muleNode := selectMuleNode(muleNodes)
		if err := sendChunkToMuleNode(chunk, muleNode); err != nil {
			return err
		}

		// Update SQLite database with chunk ID and mule node ID
		if err := updateDatabaseWithChunkInfo(i/chunkSize, muleNode); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Example usage in main
	muleNodes := []MuleNode{
		{ID: "1", Addr: "localhost:50051"},
		// Add more mule nodes as needed
	}

	// Read the file content (replace this with your own file reading logic)
	fileContent := []byte("This is the content of the file.")

	// Split and distribute the file to mule nodes
	err := SplitAndDistribute(fileContent, muleNodes)
	if err != nil {
		log.Fatal(err)
	}
}
