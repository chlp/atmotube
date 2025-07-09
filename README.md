# Atmotube PRO – Air Quality Monitor Bluetooth Reader

A minimal Go-based tool for reading real-time air quality data from the **Atmotube PRO** device via Bluetooth Low Energy (BLE).  
The tool connects, subscribes to BLE characteristics, and serves the latest environmental data via a local HTTP API.

## Features

- Connects to Atmotube PRO via Bluetooth (BLE)
- Subscribes to key GATT characteristics
- Receives live environmental data, including:
    - Particulate Matter (PM1, PM2.5, PM4, PM10)
    - Total Volatile Organic Compounds (TVOC)
    - Temperature, Humidity, Pressure
- Launches a local HTTP server on port 8092
- Serves current sensor data as structured JSON

## API Response Format

The HTTP API returns the most recent measurement in the following JSON structure:

```json
{
  "PM1": {
    "name": "pm1",
    "value": 4.38,
    "unit": "µg/m³",
    "status": "ok"
  },
  "PM25": {
    "name": "pm2_5",
    "value": 5.93,
    "unit": "µg/m³",
    "status": "ok"
  },
  "PM4": {
    "name": "pm4",
    "value": 4.63,
    "unit": "µg/m³",
    "status": "ok"
  },
  "PM10": {
    "name": "pm10",
    "value": 7.02,
    "unit": "µg/m³",
    "status": "ok"
  },
  "TVOC": {
    "name": "tvoc",
    "value": 0.109,
    "unit": "mg/m³",
    "status": "ok"
  },
  "Temp": {
    "name": "temp",
    "value": 24.4,
    "unit": "celsius",
    "status": "ok"
  },
  "Humidity": {
    "name": "humidity",
    "value": 37,
    "unit": "%",
    "status": "ok"
  },
  "Pressure": {
    "name": "pressure",
    "value": 1013.8,
    "unit": "hPa",
    "status": "ok"
  },
  "Battery": {
    "name": "battery",
    "value": 44,
    "unit": "%",
    "status": "warn"
  },
  "BluetoothConnection": {
    "name": "bluetooth_connection",
    "value": 1,
    "unit": "connected",
    "status": "ok"
  }
}
```

## Example Terminal Output

```
🔍 Searching for Atmotube...  
✅ Found: ATMOTUBE [35705aeb-5c28-a8a4-b1a4-3b1370060b09]  
🔗 Connected. Searching for services...  
🌫️ TVOC (SGPC3): 298 ppb  
🌫️ TVOC (SGPC3): 294 ppb  
🌡 Temperature: 25.0°C, 💧 Humidity: 39%, 📟 Pressure: 1012.1 hPa  
🌫️ TVOC (SGPC3): 293 ppb  
🌡 Temperature: 25.0°C, 💧 Humidity: 39%, 📟 Pressure: 1012.0 hPa  
🌫️ TVOC (SGPC3): 291 ppb  
🌡 Temperature: 25.0°C, 💧 Humidity: 39%, 📟 Pressure: 1012.1 hPa  
🌁 PM1: 3.33, PM2.5: 5.66, PM4: 5.13, PM10: 7.55 µg/m³
```

## Platform

Tested and verified on **macOS only** using the built-in Bluetooth stack.

## Protocol Reference

This project is based on the official Bluetooth API specification provided by Atmotube:  
[ATMO Bluetooth API.pdf](ATMO Bluetooth API.pdf) — stored in the root of this repository and originally available at https://support.atmotube.com/en/articles/10364981-bluetooth-api

## Requirements

- Atmotube PRO device
- BLE support enabled

## License

MIT License

## Author

[Aleksei Rytikov](https://github.com/chlp)