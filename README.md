# whereyouat

Determines a user's location based on their IP address from an TCP RPC request or HTTP

I recently found out a company i have been working for has been using an external _paid_ API to determine if an IP is inside of the EU (GDPR) is coming onto our sites to determine if we need to render a cookie banner. I went down a rabbit hole learning how this is done on how this is done (offline) and decided to sharpen my golang on the journey.

This wont give a result on local host- i just patched my IPv4 as the remote address for testing purposes!

Thanks oschwald for [maxminddb-golang](https://github.com/oschwald/maxminddb-golang) mmdb parser, i didnt feel like writing my own parser lol
Thanks P3TERX for [GeoLite.mmdb](https://github.com/P3TERX/GeoLite.mmdb) CI/CD file releases on new database update

## Features

- **TCP RPC and HTTP JSON-RPC**: RPC with client example written in go
- **IP Geolocation**: Automatic location detection from client IP
- **Auto-updating Database**: Fetches latest GeoLite2 database automatically
- **Environment Configuration**: Simple `.env` file support

## Quick Start

### Install Dependencies

```bash
make deps
```

### Run the Server

```bash
make run
```

The server will start on:

- TCP RPC: `localhost:1234`
- HTTP JSON-RPC: `localhost:8080`

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
MMDB_GITHUB_API_URL=https://api.github.com/repos/yourrepo/releases/latest
DB_DIR=data
DB_FILENAME=GeoLite2-Country.mmdb
TAG_FILE=db_version.txt
```

## API

### LocationService.Calculate

Determines the location of the client based on their IP address.

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

## License

MIT
