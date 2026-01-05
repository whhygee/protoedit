package fixtures

const (
	TestProto = `syntax = "proto3";

package test.v1;

// TestService is a test service.
service TestService {
  // GetFoo gets a foo.
  rpc GetFoo(GetFooRequest) returns (GetFooResponse) {}
}

message GetFooRequest {
  string id = 1;
}

message GetFooResponse {
  string name = 1;
}

enum Status {
  STATUS_UNSPECIFIED = 0;
  STATUS_ACTIVE = 1;
}
`

	TestProtoMinimal = `syntax = "proto3";`

	TestProtoInvalid = "this is not valid proto {{{"

	TestProtoWithComments = `syntax = "proto3";

// File comment
package test.v1;

// Service comment line 1
// Service comment line 2
service MyService {
  // RPC comment
  rpc MyRPC(MyRequest) returns (MyResponse) {}
}
`

	TestRPCContent = `  // CreateBar creates a bar.
  rpc CreateBar(CreateBarRequest) returns (CreateBarResponse) {}`

	TestFieldContent = `  // extra_field is an additional field.
  string extra_field = 99;`

	TestEnumValueContent = "  STATUS_PENDING = 99;"

	TestMessageContent = `// CreateBarRequest is the request for CreateBar.
message CreateBarRequest {
  string name = 1;
}`

	TestMessagesContent = `message CreateBarRequest {
  string name = 1;
}

message CreateBarResponse {
  string id = 1;
}`
)
