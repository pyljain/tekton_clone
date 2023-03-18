package main

import (
	"context"
	"log"
	"os"
	"tektonclone/pkg/k8s"
	"tektonclone/pkg/pipelines"
	"tektonclone/pkg/proto"
	"tektonclone/pkg/repository"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
)

func main() {

	token := os.Getenv("CONNECTION_TOKEN")
	serverAddress := os.Getenv("SERVER_ADDRESS")
	err := getEvents(token, serverAddress)
	if err != nil {
		log.Fatalf("unable to connect to the server %s", err)
	}

}

func getEvents(token string, serverAddress string) error {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	defer conn.Close()
	client := proto.NewEventsClient(conn)
	ctx := context.Background()
	er := proto.GetEventRequest{
		Token: token,
	}

	stream, err := client.GetEvent(ctx, &er)
	if err != nil {
		return err
	}

	for {
		event, err := stream.Recv()
		if err != nil {
			return err
		}

		log.Printf("Event received from the server %+v", event)
		// Clone repository
		ctx := context.Background()
		pipelineDefBytes, err := repository.GetPipelineDef(ctx, event.RepositoryName, event.CommitRef)
		if err != nil {
			return err
		}

		pd := pipelines.PipelineDef{}

		err = yaml.Unmarshal(pipelineDefBytes, &pd)
		if err != nil {
			return err
		}

		podName, err := k8s.CreatePod(ctx, event.RepositoryName, pd)
		if err != nil {
			return err
		}

		log.Printf("Pod started %s", podName)
	}
}
