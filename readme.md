Run with `go run .` you can then make requests to `localhost:3333`.

There are 3 endpoints that you can call
1. `POST /image` The payload should include
	- `patientId` the id of the patient that image is for.
	- `image` the DICOM image.
2. `GET /image/:image_id/header_attributes`
	- You can query the attributes by providing them in a query string eg. `?(0002,0000)&(0002,0001)&(0002,0002)`.
	If no query if provided then all the attributes will be returned.
3. `GET /image/:image_id/png`

I left many comments about how this would interact with other microservices and what a DB schema might be like.
