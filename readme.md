Run with `go run .` you can then make requests to `localhost:3333`.

There are 3 endpoints that you can call
1. `POST /image`
2. `GET /image/:image_id/header_attributes`
	- Add a query string to filter the attributes eg. `?(0002,0000)&(0002,0001)&(0002,0002)` and you receive the associated data.
3. `GET /image/:image_id/png`

I left many comments about how this would interact with other microservices and what a DB schema might be like.
