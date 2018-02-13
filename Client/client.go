package main

import (
	"google.golang.org/grpc"
	"log"
	pb "grpc/rpc"
	"context"
	"os"
	"io"
)

const(
	address = "localhost:50051"
)

func send_audio( vprclient pb.Vpr_VoiceprintRecognizeClient){

	log.Println("send_audio")
	pf, err := os.Open("f:/request.mp3")
	if err != nil {
		log.Println("file open err: ", err)
	}
	defer pf.Close()

	tmpdata := make([]byte, 1000)
	idx := 1
	var allsize int64
	allsize = 0
	log.Println("asize1:", allsize)

	for {
		dsize, err := pf.Read(tmpdata)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("read err: %v", err)
			break
		}
		vprclient.Send(&pb.VoiceprintRecognizeRequest{
			&pb.VoiceprintRecognizeRequest_Body{&pb.RequestBody{tmpdata[:dsize]}}})
		log.Println("Send : idx:", idx, ", dsize: ",  int64(dsize), ", err:", err)
		idx += 1
	}

	//
	//for dsize, err := pf.Read(tmpdata); dsize != 0 && err != io.EOF; {
	//	log.Println("asize:", allsize)
	//
	//	vprclient.Send(&pb.VoiceprintRecognizeRequest{
	//		&pb.VoiceprintRecognizeRequest_Body{&pb.RequestBody{tmpdata}}})
	//
	//	log.Println("Send : idx:", idx,", allsize:",allsize, ", dsize: ",  int64(dsize), ", err:", err)
	//	idx += 1
	//	allsize += int64(dsize)
	//	//pf.Seek(allsize, 0)
	//}
}

func main() {
	conn, err :=  grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewVprClient(conn)

	vprclient, err := c.VoiceprintRecognize(context.Background())
	if err != nil {
		log.Fatalf("could not recognize: %v", err)
	}

	conf := &pb.RequestConfig{1, "GIEIWB225B2H", "1234", "AAAAAA", "124356757878"}

	vprclient.Send(&pb.VoiceprintRecognizeRequest{&pb.VoiceprintRecognizeRequest_Config{conf}})

	send_audio(vprclient)

	reply, err := vprclient.CloseAndRecv()
	if err != nil {
		log.Println("Failed to recv : %v", err)
	}

	log.Println(reply.Information)

	//conf := pb.RequestConfig{1, "FH2HY352GQ", "1234", "AAAAA11111BBBBB", "139040274825"}
	//req := pb.VoiceprintRecognizeRequest{&conf}
	//r, err := c.VoiceprintRecognize(context.Background(), &req)
}