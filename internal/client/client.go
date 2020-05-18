package client

import (
	"bytes"
	"fmt"
	"github.com/zdarovich/grpc-server-client/internal/api"
	"log"
)

type Client struct {
}

func (c *Client) Message(srv api.ProxyCaller_MessageServer) error {

	buf := new(bytes.Buffer)

	for {
		log.Print("recv resp")
		req, err := srv.Recv()
		if err != nil {
			log.Printf("receive error %v", err)
			return err
		}

		if req.FinishWrite {
			log.Print("finish msg")
			break
		} else {
			log.Print("write buf")
			_, err := buf.Write(req.Data)
			if err != nil {
				log.Print(err)
			}
			log.Print("end buf")

		}
		log.Print("send resp")
		dataLen := int64(len(req.Data))
		resp := api.Response{CommittedSize: dataLen}
		if err := srv.Send(&resp); err != nil {
			log.Printf("send error %v", err)
			return err
		}
		log.Printf("proxied data len=%d", dataLen)
	}

	resp := api.Response{CommittedSize: 0}
	if err := srv.Send(&resp); err != nil {
		return err
	}
	log.Print("read final buf")
	b := make([]byte, 100)
	n, err := buf.Read(b)
	if err != nil {
		return nil
	}
	fmt.Print(string(b[:n]))

	return nil
}

