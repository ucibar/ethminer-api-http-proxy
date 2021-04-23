# HTTP Proxy For Ethminer API 

> HTTP Proxy For [ethminer](https://github.com/ethereum-mining/ethminer) TCP json-rpc 2.0 API

Since [**Ethminer**](https://github.com/ethereum-mining/ethminer) does not have a HTTP API, this small go app allows you to consume the **Ethminer** API as HTTP.

## Usage

Build the project and run from console:
```sh
go build .
```
```sh
ethminer-api-http-proxy --help
```
Only required flag is `-miner`, this your **ethminer** API address. Eg:
```sh
ethminer-api-http-proxy -miner 127.0.0.1:3333
```
If you want to change HTTP server address, use `-serve` flag. Eg:
```sh
ethminer-api-http-proxy -miner 127.0.0.1:3333 -serve :8085
ethminer-api-http-proxy -miner 127.0.0.1:3333 -serve 192.168.1.36:8085
```
And you can consume **ethminer** API as HTTP:

```sh
curl--request POST '127.0.0.1:8085' --header 'Content-Type: application/json' --data-raw '{"id":0,"jsonrpc":"2.0","method":"miner_getstatdetail"}'
```

## HTTP API
There is no endpoints except `/`. Just Send `POST` requests to `/` with **ethminer** JSON request body.

You can find all **ethminer** requests from [**ethminer API Docs**](https://github.com/ethereum-mining/ethminer/blob/master/docs/API_DOCUMENTATION.md).

### Note:
Since **ethminer** API using [**json-rpc 2.0**](https://www.jsonrpc.org/specification) and TCP connection, API requests must have `id` field for request-response mapping. But you don't have to worry about it when you using this proxy. Because, this is HTTP.

So, just set `id:0` when sending requests and ignore the `id` in responses, Eg:
```json
//Request Body
{
  "id": 0, //doesn't matter what you send, just send integer
  "jsonrpc": "2.0",
  "method": "miner_ping"
}

//Response Body
{
  "id": 1241, //just ignore it
  "jsonrpc": "2.0",
  "result": "pong"
}
```

## Caveats:
- For now, only tested on Windows 10, but probably works on linux either.
- **ethminer** API authorize not implemented. For safety, use your ethminer API in read-only mode. Details in [**ethminer API Docs**](https://github.com/ethereum-mining/ethminer/blob/master/docs/API_DOCUMENTATION.md).
- **Important:** Since this app have no tests and bad error handling, if this proxy app die because of an error or something and can't properly close the TCP connection with **ethminer** API, **ethminer** dies too. This is probably a bug in **ethminer.**

## TODO:
- Write tests
- Write better error handling and make sure to close the TCP connection properly.
- Write better log messages.(actualy, write log messages...)
- Implement and test **ethminer** authorize for write oriented requests, provide safety.
- Proper HTTP status codes for bad, invalid **ethminer** API requests. 