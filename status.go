package main

import "sync"

type Metric struct {
	Name   string  `json:"name"`
	Value  float64 `json:"value"`
	Unit   string  `json:"unit"`
	Status string  `json:"status"`
}

var mu sync.RWMutex

var data = struct {
	PM1, PM25, PM4, PM10           Metric
	TVOC, Temp, Humidity, Pressure Metric
	Battery                        Metric
	BluetoothConnection            Metric
}{
	PM1:                 Metric{"pm1", 0, "µg/m³", "none"},
	PM25:                Metric{"pm2_5", 0, "µg/m³", "none"},
	PM4:                 Metric{"pm4", 0, "µg/m³", "none"},
	PM10:                Metric{"pm10", 0, "µg/m³", "none"},
	TVOC:                Metric{"tvoc", 0, "ppb", "none"},
	Temp:                Metric{"temp", 0, "celsius", "none"},
	Humidity:            Metric{"humidity", 0, "%", "none"},
	Pressure:            Metric{"pressure", 0, "hPa", "none"},
	Battery:             Metric{"battery", 0, "%", "none"},
	BluetoothConnection: Metric{"bluetooth_connection", 0, "bool", "none"},
}

func UpdatePM(pm1, pm25, pm4, pm10 float64) {
	mu.Lock()
	data.PM1.Value = pm1
	data.PM1.Status = statusPM1(pm1)
	data.PM25.Value = pm25
	data.PM25.Status = statusPM25(pm25)
	data.PM4.Value = pm4
	data.PM4.Status = statusPM4(pm4)
	data.PM10.Value = pm10
	data.PM10.Status = statusPM10(pm10)
	mu.Unlock()
}

func statusPM1(pm float64) string {
	switch {
	case pm <= 0:
		return "none"
	case pm <= 12:
		return "ok"
	case pm <= 35.4:
		return "warn"
	default:
		return "critical"
	}
}

func statusPM25(pm float64) string {
	switch {
	case pm <= 0:
		return "none"
	case pm <= 12:
		return "ok"
	case pm <= 35.4:
		return "warn"
	default:
		return "critical"
	}
}

func statusPM4(pm float64) string {
	switch {
	case pm <= 0:
		return "none"
	case pm <= 12:
		return "ok"
	case pm <= 35.4:
		return "warn"
	default:
		return "critical"
	}
}

func statusPM10(pm float64) string {
	switch {
	case pm <= 0:
		return "none"
	case pm <= 54:
		return "ok"
	case pm <= 154:
		return "warn"
	default:
		return "critical"
	}
}

func UpdateTVOC(tvoc float64) {
	mu.Lock()
	data.TVOC.Value = tvoc
	data.TVOC.Status = statusTVOC(tvoc)
	mu.Unlock()
}

func statusTVOC(ppb float64) string {
	switch {
	case ppb <= 0:
		return "none"
	case ppb <= 250:
		return "ok"
	case ppb <= 500:
		return "warn"
	default:
		return "critical"
	}
}

func UpdateBME280(temp, humidity, pressure float64) {
	mu.Lock()
	data.Temp.Value = temp
	data.Temp.Status = statusTemperature(temp)
	data.Humidity.Value = humidity
	data.Humidity.Status = statusHumidity(humidity)
	data.Pressure.Value = pressure
	data.Pressure.Status = statusPressure(pressure)
	mu.Unlock()
}
func statusTemperature(celsius float64) string {
	switch {
	case celsius <= -100:
		return "none"
	case celsius >= 18 && celsius <= 26:
		return "ok"
	case (celsius >= 14 && celsius < 18) || (celsius > 26 && celsius <= 30):
		return "warn"
	default:
		return "critical"
	}
}

func statusHumidity(percent float64) string {
	switch {
	case percent <= 0:
		return "none"
	case percent >= 30 && percent <= 60:
		return "ok"
	case (percent >= 20 && percent < 30) || (percent > 60 && percent <= 70):
		return "warn"
	default:
		return "critical"
	}
}

func statusPressure(hpa float64) string {
	switch {
	case hpa <= 0:
		return "none"
	case hpa >= 980 && hpa <= 1020:
		return "ok"
	case (hpa >= 950 && hpa < 980) || (hpa > 1020 && hpa <= 1050):
		return "warn"
	default:
		return "critical"
	}
}

func UpdateBattery(percent float64) {
	mu.Lock()
	data.Battery.Value = percent
	data.Battery.Status = statusBattery(percent)
	mu.Unlock()
}

func statusBattery(p float64) string {
	switch {
	case p < 0:
		return "none"
	case p >= 50:
		return "ok"
	case p >= 20:
		return "warn"
	default:
		return "critical"
	}
}

func UpdateBluetoothStatus(status string) {
	v := float64(0)
	if status == "ok" {
		v = 1
	}
	mu.Lock()
	data.BluetoothConnection = Metric{
		Name:   "bluetooth_connection",
		Value:  v,
		Unit:   "connected",
		Status: status,
	}
	mu.Unlock()
}
