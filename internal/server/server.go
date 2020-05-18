package server

import (
	"bufio"
	"context"
	"github.com/zdarovich/grpc-server-client/internal/api"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"net/http"
)
const (
	chunksize int = 1024
)

type Server struct {
}

func (s *Server) Init(ctx context.Context, msg *api.UrlMessage) (*api.UrlMessage, error) {
	log.Print("start stream")
	conn, err := grpc.Dial(":7778", grpc.WithInsecure())
	if err != nil {
		log.Printf("can not connect with server %v", err)
		return nil, err
	}

	// create stream
	client := api.NewProxyCallerClient(conn)
	stream, err := client.Message(ctx)
	if err != nil {
		log.Printf("open stream error %v", err)
		return nil, err
	}

	resp, err := http.Get(msg.Url)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer resp.Body.Close()


	var (
		part  = make([]byte, chunksize)
		messageId = rand.Int31()
		streamErr error
		count = 0
		reader = bufio.NewReader(resp.Body)
		totalSentLen int64 = 0
	)
	for {
		log.Print("read chunk")
		if count, streamErr = reader.Read(part); streamErr != nil {
			break
		}
		log.Print("send chunk")

		err = stream.Send(&api.Request{
			Data:        part[:count],
			MessageId:   messageId,
			FinishWrite: false,
		})
		if err != nil {
			log.Print(err)
			return nil, err
		}
		log.Print("recv resp")

		resp, err := stream.Recv()
		if err == io.EOF {
			continue
		}
		if err != nil {
			return nil, err
		}
		log.Printf("client received %d", resp.CommittedSize)
		if resp.CommittedSize == 0 {
			break
		}
		totalSentLen += resp.CommittedSize
	}
	if streamErr == io.EOF {
		log.Print("finish stream")
		err = stream.Send(&api.Request{
			Data:        nil,
			MessageId:   messageId,
			FinishWrite: true,
		})
		if err != nil {
			return nil, err
		}
		_, err := stream.Recv()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, streamErr
	}
	log.Printf("finished with totalSent=%d", totalSentLen)

	return msg, nil

}





