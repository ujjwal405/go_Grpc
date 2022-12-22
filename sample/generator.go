//
package sample

import (
	"grpc_go/pbs/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewKeyboard() *pb.Keyboard {
	keyboard := &pb.Keyboard{
		Layout:   randomKeyboardlayout(),
		Blacklit: randomBool(),
	}
	return keyboard
}
func NewCPU() *pb.CPU {
	brand := randomBrand()
	name := randomName(brand)
	numberofcores := randomcore(2, 8)
	numberofthreads := randomcore(numberofcores, 12)
	minghz := randomfloat(2.0, 3.5)
	maxghz := randomfloat(minghz, 5.0)
	cpu := &pb.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(numberofcores),
		NumberTherads: uint32(numberofthreads),
		MinGhz:        minghz,
		MaxGhz:        maxghz,
	}
	return cpu
}
func NewGPU() *pb.GPU {
	brand := randomGPUBrand()
	name := randomGPUName(brand)
	minghz := randomfloat(1.0, 1.5)
	maxghz := randomfloat(minghz, 2.0)
	memory := &pb.Memory{
		Value: uint64(randomcore(2, 6)),
		Unit:  pb.Memory_GIGABYTE,
	}
	gpu := &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minghz,
		MaxGhz: maxghz,
		Memory: memory,
	}
	return gpu
}
func NewRAM() *pb.Memory {
	ram := &pb.Memory{
		Value: uint64(randomcore(2, 6)),
		Unit:  pb.Memory_GIGABYTE,
	}
	return ram
}
func NewSSD() *pb.Storage {
	ssd := &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(randomcore(128, 1024)),
			Unit:  pb.Memory_GIGABYTE,
		},
	}
	return ssd
}
func NewHDD() *pb.Storage {
	hdd := &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(randomcore(1, 6)),
			Unit:  pb.Memory_TERABYTE,
		},
	}
	return hdd
}
func NewScreen() *pb.Screen {
	screen := &pb.Screen{
		SizeInch:   randomFloat32(13, 17),
		Resolution: randomscreenresolution(),
		Panel:      randomscreenpanel(),
		Multitouch: randomBool(),
	}
	return screen
}
func NewLaptop() *pb.Laptop {
	brand := randomlaptopBrand()
	name := randomlaptopName(brand)
	laptop := &pb.Laptop{
		Id:       randomUid(),
		Brand:    brand,
		Name:     name,
		Cpu:      NewCPU(),
		Ram:      NewRAM(),
		Gpus:     []*pb.GPU{NewGPU()},
		Storages: []*pb.Storage{NewSSD(), NewHDD()},
		Screen:   NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomfloat(1.0, 3.0),
		},
		PriceUsd:    randomfloat(1500, 3000),
		ReleaseYear: uint32(randomcore(2015, 2019)),
		UpdatedAt:   timestamppb.Now(),
	}
	return laptop
}
