package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"grpc_go/pbs/pb"
	"grpc_go/sample"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func createlaptop(laptopclient pb.LaptopServiceClient, laptop *pb.Laptop) {
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := laptopclient.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Println("laptop already exists")
		} else {
			log.Println("laptop created successfully")
		}
		return
	}
	log.Printf("created laptop with id :%s", res.Id)
}

func testcreatelaptop(laptopclient pb.LaptopServiceClient) {
	createlaptop(laptopclient, sample.NewLaptop())
}
func testsearchlaptop(laptopclient pb.LaptopServiceClient) {
	for i := 0; i < 10; i++ {
		createlaptop(laptopclient, sample.NewLaptop())
	}
	filter := &pb.Filter{
		MaxPrice:    2500,
		MinCpuCores: 4,
		MinCpuGhz:   2,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}
	searchlaptop(laptopclient, filter)
}
func uploadimage(laptopclient pb.LaptopServiceClient, laptopid string, imagepath string) {
	file, err := os.Open(imagepath)
	if err != nil {
		fmt.Println("cannot open file ", err)
	}
	defer file.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := laptopclient.UploadImage(ctx)
	if err != nil {
		log.Fatal("cannnot upload image ", stream.RecvMsg(nil))
	}
	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.Imageinfo{
				LaptopId:  laptopid,
				ImageType: filepath.Ext(imagepath),
			},
		},
	}
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send information ", stream.RecvMsg(nil))
	}
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read data into buffer", err)
		}
		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_Chunkdata{
				Chunkdata: buffer[:n],
			},
		}
		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send data to server", stream.RecvMsg(nil))
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response ", err)
	}
	log.Printf("uploaded an image with id %s  of size %d", res.GetId(), res.GetSize())
}
func testuploadimage(laptopclient pb.LaptopServiceClient) {
	laptop := sample.NewLaptop()
	createlaptop(laptopclient, laptop)
	uploadimage(laptopclient, laptop.GetId(), "C:\\photo\\pexels-life-of-pix-7974-ImResizer.jpg")
}

func searchlaptop(laptopclient pb.LaptopServiceClient, filter *pb.Filter) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopclient.SearchLaptop(ctx, req)
	if err != nil {
		log.Println("cannot search laptop")
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("cannnot receive laptop")
		}
		laptop := res.GetLaptop()
		log.Printf("found laptop with id %s", laptop.GetId())
	}
}
func ratelaptop(laptopclient pb.LaptopServiceClient, laptopid []string, score []float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := laptopclient.RateLaptop(ctx)
	if err != nil {
		return fmt.Errorf("cannot send ratelaptop request %v", err)
	}
	response := make(chan error)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("no more  data to receive")
				response <- nil
				return
			}
			if err != nil {
				response <- fmt.Errorf("cannot receive response ")
				return
			}
			log.Println("received a response ", res)
		}

	}()
	for index, laptopid := range laptopid {
		req := &pb.RateLaptopRequest{
			LaptopId: laptopid,
			Score:    score[index],
		}
		err = stream.Send(req)
		if err != nil {
			return fmt.Errorf("cannnot send response %v", stream.RecvMsg(nil))
		}
		fmt.Printf("sent a request with %s", laptopid)
	}
	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("error occured during closing stream %v", err)
	}
	err = <-response
	return err
}
func testratelaptop(laptopclient pb.LaptopServiceClient) {
	n := 3
	laptopid := make([]string, n)
	for i := 0; i < n; i++ {
		laptop := sample.NewLaptop()
		laptopid[i] = laptop.GetId()
		createlaptop(laptopclient, laptop)
	}
	score := make([]float64, n)
	for j := 0; j < n; j++ {
		score[j] = sample.RandomScore()
	}
	err := ratelaptop(laptopclient, laptopid, score)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	serveraddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serveraddress)

	conn, err := grpc.Dial(*serveraddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Cannot dial ")
	}
	laptopclient := pb.NewLaptopServiceClient(conn)
	testuploadimage(laptopclient)
	testratelaptop(laptopclient)
	testsearchlaptop(laptopclient)
	testcreatelaptop(laptopclient)

}
