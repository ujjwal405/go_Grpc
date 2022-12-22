package service

import (
	"context"
	"errors"
	"fmt"
	"grpc_go/pbs/pb"
	"sync"

	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("this laptop already exist")

type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Find(id string) *pb.Laptop
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}
type MemoryStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

func NewMemory() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]*pb.Laptop),
	}
}
func (store *MemoryStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("cannot copy :%v", err)
	}
	store.data[other.Id] = other
	return nil
}
func (store *MemoryStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy %v", err)
	}
	return other, nil
}
func (store *MemoryStore) Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	for _, laptop := range store.data {
		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			return errors.New("context is cancelled")
		}
		if Qualified(filter, laptop) {
			other, err := copied(laptop)
			if err != nil {
				return err
			}
			err = found(other)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func Qualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPrice() {
		return false
	}
	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}
	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}
	if tobit(laptop.GetRam()) < tobit(filter.GetMinRam()) {
		return false
	}
	return true
}
func tobit(memory *pb.Memory) uint64 {
	value := memory.GetValue()
	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value * 8
	case pb.Memory_KILOBYTE:
		return value * 1024 * 8
	case pb.Memory_MEGABYTE:
		return value << 23
	case pb.Memory_GIGABYTE:
		return value << 33
	case pb.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}
func copied(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy %v", err)
	}
	return other, nil
}
