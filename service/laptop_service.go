package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"grpc_go/pbs/pb"
	"io"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxsize = 1 << 20

type LaptopServer struct {
	pb.UnimplementedLaptopServiceServer
	store      MemoryStore
	imagestore Imagememory
	ratestore  Ratingmemory
}

func NewLaptop() *LaptopServer {
	return &LaptopServer{store: *NewMemory(), imagestore: *Newimagestore("C:\\photos"), ratestore: *NewRatingMemory()}
}
func (server *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("received a laptop with id :%s", laptop.Id)
	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop uuid is not valid :%v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot create new uid :%v", err)
		}
		laptop.Id = id.String()
	}
	if ctx.Err() == context.DeadlineExceeded {
		log.Println("deadline exceeded")
		return nil, status.Errorf(codes.Internal, "deadline exceeded")
	}
	err := server.store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannnot create laptop %v", err)
	}
	fmt.Printf("laptop created with id:%s", laptop.Id)
	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil
}
func (server *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("received a laptop with filter: %v", filter)
	err := server.store.Search(stream.Context(), filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{
				Laptop: laptop,
			}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			log.Printf("sent a laptop with id %v", laptop.GetId())
			return nil
		})

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error %v", err)
	}
	return nil
}
func (server *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot receive image")
	}
	laptopid := req.GetInfo().GetLaptopId()
	imagetype := req.GetInfo().GetImageType()
	log.Printf("received a laptop with id %s with imagetype %s", laptopid, imagetype)
	laptop, err := server.store.Find(laptopid)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot find laptop %v", err)
	}
	if laptop == nil {
		return status.Errorf(codes.InvalidArgument, "cannot find laptop wiht id %s", laptopid)
	}
	imagedata := bytes.Buffer{}
	imagesize := 0

	for {
		if err = ctxerr(stream.Context()); err != nil {
			return err
		}
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("end of file")
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannnot receive chunk %v", err)
		}
		chunk := req.GetChunkdata()
		size := len(chunk)
		imagesize += size
		if imagesize > maxsize {
			return status.Errorf(codes.InvalidArgument, "imagesize is too large %d", imagesize)
		}
		_, err = imagedata.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot write in buffer")
		}

	}
	imageid, err := server.imagestore.Save(laptopid, imagetype, imagedata)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot save image to store %v", err)
	}
	res := &pb.ImageResponse{
		Id:   imageid,
		Size: uint32(imagesize),
	}
	err = stream.SendAndClose(res)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot send response %v", err)
	}
	log.Printf("successfull saved with laptop id %s", laptopid)
	return nil
}
func (server *LaptopServer) RateLaptop(stream pb.LaptopService_RateLaptopServer) error {
	for {
		if err := ctxerr(stream.Context()); err != nil {
			return err
		}
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive data")
		}
		laptopid := req.GetLaptopId()
		score := req.GetScore()
		found, err := server.store.Find(laptopid)
		if err != nil {
			return status.Errorf(codes.Internal, "error while finding laptop")
		}
		if found == nil {
			return status.Errorf(codes.NotFound, "searched result not found with laptop id %s", laptopid)
		}
		rating, err := server.ratestore.Add(laptopid, score)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot add to rate store")
		}
		res := &pb.RateLaptopResponse{
			LaptopId:  laptopid,
			RateCount: rating.Count,
			Score:     rating.Sum / float64(rating.Count),
		}
		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot send response to the client")
		}

	}
	return nil
}
func ctxerr(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return (status.Errorf(codes.Canceled, "request is cancelled"))
	case context.DeadlineExceeded:
		return (status.Errorf(codes.DeadlineExceeded, "deadline exceeded"))
	default:
		return nil
	}
}
