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
	uuidAtmotubeService      = "db450001-8e9a-4818-add7-6ed94a328ab4"
	uuidSGPC3Characteristic  = "db450002-8e9a-4818-add7-6ed94a328ab4"
	uuidBME280Characteristic = "db450003-8e9a-4818-add7-6ed94a328ab4"
	uuidStatusCharacteristic = "db450004-8e9a-4818-add7-6ed94a328ab4"
	uuidPMCharacteristic     = "db450005-8e9a-4818-add7-6ed94a328ab4"
)

func connectToAtmotube() {
	must("enable adapter", adapter.Enable())

	UpdateBluetoothStatus("critical")

	fmt.Println("🔍 Searching for Atmotube...")
	ch := make(chan bluetooth.ScanResult, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		_ = adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
			if strings.TrimSpace(result.LocalName()) == "ATMOTUBE" {
				fmt.Printf("✅ Found: %s [%s]\n", result.LocalName(), result.Address.String())
				ch <- result
				_ = adapter.StopScan()
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

	services, err := device.DiscoverServices([]bluetooth.UUID{uuidFromString(uuidAtmotubeService)})
	must("discovering services", err)

	for _, service := range services {
		chars, err := service.DiscoverCharacteristics(nil)
		must("discovering characteristics", err)
		for _, char := range chars {
			switch char.UUID().String() {
			case uuidSGPC3Characteristic:
				subscribeSGPC3(char)
			case uuidBME280Characteristic:
				subscribeBME280(char)
			case uuidStatusCharacteristic:
				subscribeStatus(char)
			case uuidPMCharacteristic:
				subscribePM(char)
			default:
				fmt.Printf("ℹ️ Wrong characteristic: %s\n", char.UUID().String())
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
	err := char.EnableNotifications(func(buf []byte) {
		if len(buf) >= 2 {
			tvocPpb := binary.LittleEndian.Uint16(buf[0:2])
			UpdateTVOC(float64(tvocPpb))
			fmt.Printf("🌫️ TVOC (SGPC3): %d ppb\n", tvocPpb)
		}
	})
	must("SGPC3", err)
}

func subscribeBME280(char bluetooth.DeviceCharacteristic) {
	err := char.EnableNotifications(func(buf []byte) {
		if len(buf) >= 8 {
			humidity := buf[0]
			pressurePa := binary.LittleEndian.Uint32(buf[2:6])
			pressureHpa := float64(pressurePa) / 100.0
			tp := binary.LittleEndian.Uint16(buf[6:8])
			temp := float64(tp) / 100.0
			UpdateBME280(temp, float64(humidity), pressureHpa)
			fmt.Printf("🌡 Temperature: %.1f°C, 💧 Humidity: %d%%, 📟 Pressure: %.1f hPa\n", temp, humidity, pressureHpa)
		}
	})
	must("BME280", err)
}

func subscribeStatus(char bluetooth.DeviceCharacteristic) {
	err := char.EnableNotifications(func(buf []byte) {
		if len(buf) >= 2 {
			battery := buf[1]
			UpdateBattery(float64(battery))
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
			raw1 := parse3BytesLE(buf[0:3])
			raw25 := parse3BytesLE(buf[3:6])
			raw10 := parse3BytesLE(buf[6:9])
			raw4 := parse3BytesLE(buf[9:12])

			if raw1 == 0xFFFFFF && raw25 == 0xFFFFFF && raw10 == 0xFFFFFF && raw4 == 0xFFFFFF {
				fmt.Println("⚠️ Skipping invalid PM data (0xFFFFFF)")
				return
			}

			pm1 := float64(raw1) / 100.0
			pm25 := float64(raw25) / 100.0
			pm10 := float64(raw10) / 100.0
			pm4 := float64(raw4) / 100.0

			UpdatePM(pm1, pm25, pm4, pm10)
			fmt.Printf("🌁 PM1: %.2f, PM2.5: %.2f, PM4: %.2f, PM10: %.2f µg/m³\n", pm1, pm25, pm4, pm10)
		} else {
			fmt.Printf("⚠️ Wrong data for PM: %d byte\n", len(buf))
		}
	})
	must("PM", err)
}

func must(msg string, err error) {
	if err != nil {
		log.Fatalf("❌ %s: %v", msg, err)
	}
}
