package sample

import (
	"grpc_go/pbs/pb"
	"math/rand"

	"github.com/google/uuid"
)

func randomKeyboardlayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}
func randomBool() bool {
	return rand.Intn(2) == 1
}
func randomset(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}
func randomBrand() string {
	return randomset("Intel", "AMD")
}
func randomName(brand string) string {
	if brand == "Intel" {
		return randomset("3xx - Celeron D",
			"4xx - Celeron",
			"5xx - Pentium 4",
			"6xx - Pentium 4",
			"8xx - Pentium D and Pentium Extreme Edition",
			"9xx - Pentium D and Pentium Extreme Edition",
			"E1xxx - Celeron Dual-Core",
			"E2xxx - Pentium Dual-Core")
	}
	return randomset("Ryzen 3", "Ryzen 5", "Ryzen 7")
}
func randomcore(min int, max int) int {
	return min + (rand.Intn(max - min + 1))
}
func randomfloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
func randomGPUBrand() string {
	return randomset("NIVDIA", "AMD")
}
func randomGPUName(brand string) string {
	if brand == "NIVDIA" {
		return randomset(
			" GTX 560 Ti",
			" RTX 3080",
			"RTX 3070 ")
	}
	return randomset(
		"AMD Radeon™ RX 6700",
		"AMD Radeon™ RX 6600 XT",
		"AMD Radeon™ RX 6600",
	)
}
func randomFloat32(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
func randomscreenpanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}
func randomscreenresolution() *pb.Screen_Resolution {
	height := randomcore(1000, 4320)
	width := height * 16 / 9
	resolution := &pb.Screen_Resolution{
		Width:  uint32(width),
		Height: uint32(height),
	}
	return resolution
}
func randomUid() string {
	return uuid.New().String()
}
func randomlaptopBrand() string {
	return randomset("Apple", "Dell", "Lenovo")
}
func randomlaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomset("Macbook Air", "Macbook Pro")
	case "Dell":
		return randomset("Latitude", "Vostro")
	default:
		return randomset("Thinkpad X1", "Thinkpad P3")
	}
}
func RandomScore() float64 {
	return float64(rand.Intn(10))
}
