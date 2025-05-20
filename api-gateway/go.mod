module api-gateway

go 1.24.2

require (
	card-service v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.72.0
	user-service v0.0.0-00010101000000-000000000000
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace card-service => ../card-service

replace user-service => ../user-service
