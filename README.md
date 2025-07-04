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
    "pm1":      {"name": "PM1",        "value": 3.33,   "unit": "µg/m³"},
    "pm2_5":    {"name": "PM2.5",      "value": 5.66,   "unit": "µg/m³"},
    "pm4":      {"name": "PM4",        "value": 5.13,   "unit": "µg/m³"},
    "pm10":     {"name": "PM10",       "value": 7.55,   "unit": "µg/m³"},
    "tvoc":     {"name": "TVOC",       "value": 291.0,  "unit": "ppb"},
    "temp":     {"name": "Temperature","value": 25.0,   "unit": "°C"},
    "humidity": {"name": "Humidity",   "value": 39.0,   "unit": "%"},
    "pressure": {"name": "Pressure",   "value": 1012.1, "unit": "hPa"}
}
```

## Example Terminal Output

```bash
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
**ATMO Bluetooth API.pdf** — stored in the root of this repository and originally available at https://support.atmotube.com/en/articles/10364981-bluetooth-api

## Requirements

- Atmotube PRO device
- BLE support enabled

## License

MIT License

## Author

[Aleksei Rytikov](https://github.com/chlp)