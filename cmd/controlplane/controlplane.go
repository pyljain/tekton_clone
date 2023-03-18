package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"tektonclone/pkg/db"
	"tektonclone/pkg/proto"
	"tektonclone/pkg/runner"
	"tektonclone/pkg/signing"

	"google.golang.org/grpc"
)

func main() {

	cs := os.Getenv("CONN_STRING")
	database, err := db.NewPostgres(cs)
	if err != nil {
		log.Fatalf("Error occured while connecting to the database %s", err)
	}

	runnerChannels := make(map[int]chan PipelineInvocationRequest)

	http.HandleFunc("/invoke", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		payload := PipelineInvocationRequest{}
		payloadBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(payloadBytes, &payload)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		runnerId, err := database.FindRunnerByRepository(r.Context(), payload.RepositoryURL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if runnerChannel, ok := runnerChannels[runnerId]; ok {
			runnerChannel <- payload
		}
	})

	http.HandleFunc("/link", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		payload := runner.RunnerMetadata{}
		payloadBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(payloadBytes, &payload)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		err = database.CreateLink(ctx, payload.RunnerId, payload.Repository)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	// Create a new runner
	http.HandleFunc("/runner", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		runnerBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newRunner := runner.Runner{}
		err = json.Unmarshal(runnerBytes, &newRunner)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error %s", err)
			return
		}

		token, err := newRunner.CreateToken(r.Context(), database)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error %s", err)
			return
		}

		w.Write([]byte(token))
	})

	go StartGRPCServer(runnerChannels)
	http.ListenAndServe(":8090", nil)
}

type PipelineInvocationRequest struct {
	RepositoryURL string `json:"repository"`
	Ref           string `json:"ref"`
}

func StartGRPCServer(runnerChannels map[int]chan PipelineInvocationRequest) error {
	lis, err := net.Listen("tcp", ":8091")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	proto.RegisterEventsServer(s, &GrpcServer{
		runnerChannels: runnerChannels,
	})

	err = s.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

type GrpcServer struct {
	runnerChannels map[int]chan PipelineInvocationRequest
	proto.UnimplementedEventsServer
}

func (s *GrpcServer) GetEvent(req *proto.GetEventRequest, stream proto.Events_GetEventServer) error {
	runnerId, err := signing.ValidateToken(req.Token)
	if err != nil {
		return err
	}

	s.runnerChannels[runnerId] = make(chan PipelineInvocationRequest)

	for message := range s.runnerChannels[runnerId] {
		stream.Send(&proto.Event{
			EventId:        "1",
			Type:           proto.EventType_executePipeline,
			RepositoryName: message.RepositoryURL,
			CommitRef:      message.Ref,
		})
	}

	return nil
}

