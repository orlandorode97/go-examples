proto-gen:
     protoc \
     --proto_path proto \
     --go_out "generated/" --go_opt=paths=source_relative \
     --go-grpc_out "generated/" --go-grpc_opt=paths=source_relative \
     --grpc-gateway_out "generated/" --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=generate_unbound_methods=true \
     proto/**/*.proto
