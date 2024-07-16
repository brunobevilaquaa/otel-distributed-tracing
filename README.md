# Weather Check

This project aims to develop a system in Go that receives a ZIP code (CEP), identifies the city, and returns the current weather (temperature in Celsius, Fahrenheit, and Kelvin) along with the city name. This system implement OTEL (Open Telemetry) and Zipkin.

## Features

- **Receive ZIP Code:** The system accepts a valid 8-digit ZIP code (CEP) as input.
- **Location Lookup:** It looks up the ZIP code to determine the location name.
- **Temperature Conversion:** It fetches the current temperature and returns it in Celsius, Fahrenheit, and Kelvin.
- **Success Response:**
    - **HTTP Code:** 200
    - **Response Body:** `{ "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65 }`
- **Invalid ZIP Code Response:**
    - **HTTP Code:** 422
    - **Message:** `invalid zipcode`
- **ZIP Code Not Found Response:**
    - **HTTP Code:** 404
    - **Message:** `cannot find zipcode`

## Getting Started

### Prerequisites

- Ensure you have Docker and Docker Compose installed.
- Obtain an API key for the weather service and export it as an environment variable.

### Running the Project Locally

1. Export the `WEATHER_API_KEY` environment variable:
   ```sh
   export WEATHER_API_KEY=your_api_key_here
   ```
2. Start the project using Docker Compose:
   ```sh
   docker-compose up
   ```

### Usage

To get weather information for a specific ZIP code (CEP), make an HTTP POST request to the following endpoint:

```sh
curl -XPOST -H "Content-type: application/json" -d '{"cep": "99999999"}' 'http://localhost:8080/api/v1/weather-check'
```

Replace `99999999` with the ZIP code you want to check.
