package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

var (
	battery uint8

	tvocPpb uint16

	humidity    uint8
	pressureHpa float64
	temp        float64

	pm1  float64
	pm25 float64
	pm10 float64
	pm4  float64
)

type Measurement struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type SensorData struct {
	PM1      Measurement `json:"pm1"`
	PM25     Measurement `json:"pm2_5"`
	PM4      Measurement `json:"pm4"`
	PM10     Measurement `json:"pm10"`
	TVOC     Measurement `json:"tvoc"`
	Temp     Measurement `json:"temp"`
	Humidity Measurement `json:"humidity"`
	Pressure Measurement `json:"pressure"`
}

func uuidFromString(s string) bluetooth.UUID {
	s = strings.ReplaceAll(s, "-", "")
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Fatalf("❌ hex.DecodeString(%s) error: %v", s, err)
	}
	if len(b) != 16 {
		log.Fatalf("❌ UUID should be 16 byte, received %d", len(b))
	}

	var uuid [16]byte
	copy(uuid[:], b)

	return bluetooth.NewUUID(uuid)
}

var (
	uuidAtmotubeService      = uuidFromString("DB450001-8E9A-4818-ADD7-6ED94A328AB4")
	uuidSGPC3Characteristic  = uuidFromString("DB450002-8E9A-4818-ADD7-6ED94A328AB4")
	uuidBME280Characteristic = uuidFromString("DB450003-8E9A-4818-ADD7-6ED94A328AB4")
	uuidStatusCharacteristic = uuidFromString("DB450004-8E9A-4818-ADD7-6ED94A328AB4")
	uuidPMCharacteristic     = uuidFromString("DB450005-8E9A-4818-ADD7-6ED94A328AB4")
)

func main() {
	must("turn on adapter", adapter.Enable())

	fmt.Println("🔍 Searching for Atmotube...")
	var device bluetooth.Device
	ch := make(chan bluetooth.ScanResult, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
			if strings.TrimSpace(result.LocalName()) == "ATMOTUBE" {
				fmt.Printf("✅ Найдено: %s [%s]\n", result.LocalName(), result.Address.String())
				ch <- result
				adapter.StopScan()
			}
		})
	}()

	var result bluetooth.ScanResult
	select {
	case result = <-ch:
	case <-ctx.Done():
		log.Fatalln("❌ Device not found")
	}

	var err error
	device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
	must("подключение", err)

	fmt.Println("🔗 Connected. Searching for services...")

	services, err := device.DiscoverServices([]bluetooth.UUID{uuidAtmotubeService})
	must("searching for services", err)

	for _, service := range services {
		chars, err := service.DiscoverCharacteristics(nil)
		must("searching for characteristics", err)

		for _, char := range chars {
			switch char.UUID().String() {
			case uuidSGPC3Characteristic.String():
				subscribeSGPC3(char)
			case uuidBME280Characteristic.String():
				subscribeBME280(char)
			case uuidStatusCharacteristic.String():
				subscribeStatus(char)
			case uuidPMCharacteristic.String():
				subscribePM(char)
			default:
				fmt.Printf("ℹ️ Wrong characteristic: %s\n", char.UUID().String())
			}
		}
	}

	select {}
}

func must(msg string, err error) {
	if err != nil {
		log.Fatalf("❌ Fail %s: %v", msg, err)
	}
}

func subscribeSGPC3(char bluetooth.DeviceCharacteristic) {
	err := char.EnableNotifications(func(buf []byte) {
		if len(buf) >= 2 {
			tvocPpb = binary.LittleEndian.Uint16(buf[0:2])
			fmt.Printf("🌫️ TVOC (SGPC3): %d ppb\n", tvocPpb)
		}
	})
	must("SGPC3", err)
}

func subscribeBME280(char bluetooth.DeviceCharacteristic) {
	err := char.EnableNotifications(func(buf []byte) {
		if len(buf) >= 8 {
			humidity = buf[0]
			pressurePa := binary.LittleEndian.Uint32(buf[2:6])
			pressureHpa = float64(pressurePa) / 100.0
			tp := binary.LittleEndian.Uint16(buf[6:8])
			temp = float64(tp) / 100.0
			fmt.Printf("🌡 Temperature: %.1f°C, 💧 Humidity: %d%%, 📟 Pressure: %.1f hPa\n", temp, humidity, pressureHpa)
		}
	})
	must("BME280", err)
}

func subscribeStatus(char bluetooth.DeviceCharacteristic) {
	err := char.EnableNotifications(func(buf []byte) {
		if len(buf) >= 2 {
			battery = buf[1]
			fmt.Printf("🔋 Battery: %d%%\n", battery)
		}
	})
	must("Status", err)
}

func parse3BytesLE(b []byte) uint32 {
	if len(b) != 3 {
		return 0
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

func subscribePM(char bluetooth.DeviceCharacteristic) {
	err := char.EnableNotifications(func(buf []byte) {
		if len(buf) >= 12 {
			pm1 = float64(parse3BytesLE(buf[0:3])) / 100.0
			pm25 = float64(parse3BytesLE(buf[3:6])) / 100.0
			pm10 = float64(parse3BytesLE(buf[6:9])) / 100.0
			pm4 = float64(parse3BytesLE(buf[9:12])) / 100.0
			fmt.Printf("🌁 PM1: %.2f, PM2.5: %.2f, PM4: %.2f, PM10: %.2f µg/m³\n", pm1, pm25, pm4, pm10)
		} else {
			fmt.Printf("⚠️ Wrong data for PM: %d byte\n", len(buf))
		}
	})
	must("PM", err)
}
