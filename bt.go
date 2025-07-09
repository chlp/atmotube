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

func connectToAtmotube() {
	must("enable adapter", adapter.Enable())

	UpdateBluetoothStatus("critical")

	fmt.Println("🔍 Searching for Atmotube...")
	ch := make(chan bluetooth.ScanResult, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
			if strings.TrimSpace(result.LocalName()) == "ATMOTUBE" {
				fmt.Printf("✅ Found: %s [%s]\n", result.LocalName(), result.Address.String())
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

	device, err := adapter.Connect(result.Address, bluetooth.ConnectionParams{})
	must("connecting", err)
	fmt.Println("🔗 Connected. Discovering services...")

	UpdateBluetoothStatus("warn")

	services, err := device.DiscoverServices([]bluetooth.UUID{uuidFromString("DB450001-8E9A-4818-ADD7-6ED94A328AB4")})
	must("discovering services", err)

	for _, service := range services {
		chars, err := service.DiscoverCharacteristics(nil)
		must("discovering characteristics", err)
		for _, char := range chars {
			switch char.UUID() {
			case uuidFromString("DB450002-8E9A-4818-ADD7-6ED94A328AB4"):
				subscribeSGPC3(char)
			case uuidFromString("DB450003-8E9A-4818-ADD7-6ED94A328AB4"):
				subscribeBME280(char)
			case uuidFromString("DB450004-8E9A-4818-ADD7-6ED94A328AB4"):
				subscribeStatus(char)
			case uuidFromString("DB450005-8E9A-4818-ADD7-6ED94A328AB4"):
				subscribePM(char)
			default:
				fmt.Printf("ℹ️ Unknown characteristic: %s\n", char.UUID().String())
			}
		}
	}

	UpdateBluetoothStatus("ok")
}

func uuidFromString(s string) bluetooth.UUID {
	s = strings.ReplaceAll(s, "-", "")
	b, err := hex.DecodeString(s)
	must("decode UUID", err)
	var uuid [16]byte
	copy(uuid[:], b)
	return bluetooth.NewUUID(uuid)
}

func subscribeSGPC3(char bluetooth.DeviceCharacteristic) {
	char.EnableNotifications(func(buf []byte) {
		if len(buf) >= 2 {
			tvocPpb := binary.LittleEndian.Uint16(buf[0:2])
			UpdateTVOC(float64(tvocPpb) / 1000.0)
		}
	})
}

func subscribeBME280(char bluetooth.DeviceCharacteristic) {
	char.EnableNotifications(func(buf []byte) {
		if len(buf) >= 8 {
			humidity := float64(buf[0])
			pressurePa := binary.LittleEndian.Uint32(buf[2:6])
			pressure := float64(pressurePa) / 100.0
			temp := float64(binary.LittleEndian.Uint16(buf[6:8])) / 100.0
			UpdateBME280(temp, humidity, pressure)
		}
	})
}

func subscribeStatus(char bluetooth.DeviceCharacteristic) {
	char.EnableNotifications(func(buf []byte) {
		if len(buf) >= 2 {
			battery := buf[1]
			UpdateBattery(float64(battery))
			fmt.Printf("🔋 Battery: %d%%\n", battery)
		}
	})
}

func parse3BytesLE(b []byte) uint32 {
	if len(b) != 3 {
		return 0
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

func subscribePM(char bluetooth.DeviceCharacteristic) {
	char.EnableNotifications(func(buf []byte) {
		if len(buf) >= 12 {
			pm1 := float64(parse3BytesLE(buf[0:3])) / 100.0
			pm25 := float64(parse3BytesLE(buf[3:6])) / 100.0
			pm10 := float64(parse3BytesLE(buf[6:9])) / 100.0
			pm4 := float64(parse3BytesLE(buf[9:12])) / 100.0
			UpdatePM(pm1, pm25, pm4, pm10)
		}
	})
}

func must(msg string, err error) {
	if err != nil {
		log.Fatalf("❌ %s: %v", msg, err)
	}
}
