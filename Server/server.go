package main

import (
	"net"
	"log"
	"google.golang.org/grpc"
	pb "grpc/rpc"
	"net/http"
	"io"
	"fmt"
	"os"
)

const (
	port = ":50051"
)

type server struct{
}

func (s *server) VoiceprintRecognize(stream pb.Vpr_VoiceprintRecognizeServer) error {
	idx := 0

	gid := ""
	pid := ""
	imei := ""

	fd, err := os.OpenFile("f:/server.recv", os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Println("open file err", err)
	}
	defer fd.Close()

	for {
		idx += 1
		audio, err := stream.Recv()
		if idx == 1 {
			log.Println("server gid", audio.GetConfig().GetGid())
			log.Println("server pid", audio.GetConfig().GetPid())
			log.Println("server imei", audio.GetConfig().GetImei())
			gid = audio.GetConfig().GetGid()
			pid = audio.GetConfig().GetPid()
			imei = audio.GetConfig().GetImei()
		}

		if err == io.EOF {

			stas := pb.Status{http.StatusOK, http.StatusText(http.StatusOK), nil}
			return stream.SendAndClose(&pb.VoiceprintRecognizeResponse{
				&stas,
				"echo:"+gid + ", "+ pid +", "+ imei,
			})
		}
		if err != nil {
			log.Println("Failed to recv : %v", err)
			return err
		}

		if idx != 1 {

			num, err := fd.Write(audio.GetBody().Body)
			if err != nil {
				log.Println("write err", err)
			}
			log.Println("idx:", idx, ", save num: ", num, ", body_len:", len(audio.GetBody().Body))
		}
	}
	fmt.Println("voiceprint success")
	return nil
}

func (s *server) NewServer(){

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen :%v", err)
	}

	//var opts []grpc.ServerOption

	ser := grpc.NewServer()
	pb.RegisterVprServer(ser, s)

	if err := ser.Serve(lis); err != nil{
		log.Fatal("failed to serve: %v", err)
	}
}

func main(){

	srv := server{}
	srv.NewServer()
}