# whereyouat

Determine a client’s location based on their IP address via TCP RPC or HTTP.

I recently found out that a company I’ve been working for has been using an external _paid_ API to check if an IP is in the EU (GDPR). The main reason? To figure out if a cookie banner needs to be shown. I went down a rabbit hole learning how this works offline and decided to sharpen my Go skills along the way.  

This tool can also do full IP tracking and provide city-level geolocation, but my goal was just a fast way to check a client’s country.  

> Note: This won’t return results for localhost—I patched my IPv4 as the remote address for testing purposes.

Big thanks to:  
- [oschwald/maxminddb-golang](https://github.com/oschwald/maxminddb-golang) – MMDB parser, because writing my own felt like a lot lol  
- [P3TERX/GeoLite.mmdb](https://github.com/P3TERX/GeoLite.mmdb) – CI/CD releases of updated databases  

## Features

- **TCP RPC and HTTP JSON-RPC** – Includes a Go client example  
- **IP Geolocation** – Automatic location detection from client IP  
- **Auto-updating Database** – Fetches the latest GeoLite2 database automatically  
- **Environment Configuration** – Simple `.env` support  

## Quick Start

### Install Dependencies

```bash
make deps
````

### Run the Server

```bash
make run
```

Server endpoints:

* TCP RPC: `localhost:1234`
* HTTP JSON-RPC: `localhost:8080`

### Test with Clients

**TCP Client:**

```bash
make tcp-client
```

**cURL:**

```bash
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{"method":"LocationService.Calculate","params":[{}],"id":1}'
```

## Configuration

Create a `.env` file (see `.env.example`):

```env
TCP_ADDRESS=:1234
HTTP_ADDRESS=:8080
MMDB_GITHUB_API_URL=https://api.github.com/repos/P3TERX/GeoLite.mmdb/releases/latest
DB_DIR=db
DB_FILENAME=GeoLite2-Country.mmdb
TAG_FILE=GeoLite2-Country.mmdb.tag
```

## API

### LocationService.Calculate

Determine a client’s location based on their IP address.

**Request:**

```json
{
  "method": "LocationService.Calculate",
  "params": [{}],
  "id": 1
}
```

**Response:**

```json
{
  "result": {
    "Location": "United States",
    "IsoCode": "US"
  },
  "error": null,
  "id": 1
}
```

## Makefile Commands

```bash
make help         # Show available commands
make build        # Build server and client binaries
make run          # Run the server
make tcp-client   # Run TCP client example
make test         # Run tests
make clean        # Remove build artifacts
make deps         # Download dependencies
```
