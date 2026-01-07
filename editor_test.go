package protoedit

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/whhygee/protoedit/testdata"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		proto   string
		wantErr bool
	}{
		{
			name:    "valid proto",
			proto:   testdata.TestProto,
			wantErr: false,
		},
		{
			name:    "empty proto",
			proto:   "",
			wantErr: false,
		},
		{
			name:    "minimal proto",
			proto:   testdata.TestProtoMinimal,
			wantErr: false,
		},
		{
			name:    "invalid proto",
			proto:   testdata.TestProtoInvalid,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor, err := New(tt.proto)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && editor == nil {
				t.Error("New() returned nil editor")
			}
		})
	}
}

func TestAppendToService(t *testing.T) {
	tests := []struct {
		name        string
		proto       string
		serviceName string
		content     string
		check       func(t *testing.T, result string)
	}{
		{
			name:        "appends RPC to end of service",
			proto:       testdata.TestProto,
			serviceName: "TestService",
			content:     testdata.TestRPCContent,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "rpc CreateBar") {
					t.Error("result should contain the new RPC")
				}

				createBarPos := strings.Index(result, "rpc CreateBar")
				getFooPos := strings.Index(result, "rpc GetFoo")
				if createBarPos < getFooPos {
					t.Error("CreateBar should be after GetFoo")
				}

				if !strings.Contains(result, "// GetFoo gets a foo.") {
					t.Error("original comments should be preserved")
				}

				if _, err := New(result); err != nil {
					t.Errorf("result should be valid proto: %v", err)
				}
			},
		},
		{
			name:        "preserves all comments",
			proto:       testdata.TestProtoWithComments,
			serviceName: "MyService",
			content:     "  rpc NewRPC(NewRequest) returns (NewResponse) {}",
			check: func(t *testing.T, result string) {
				comments := []string{
					"// File comment",
					"// Service comment line 1",
					"// Service comment line 2",
					"// RPC comment",
				}
				for _, c := range comments {
					if !strings.Contains(result, c) {
						t.Errorf("comment %q should be preserved", c)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor, err := New(tt.proto)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if err := editor.Append.ToService(tt.serviceName, tt.content); err != nil {
				t.Fatalf("ToService() error = %v", err)
			}

			tt.check(t, editor.String())
		})
	}
}

func TestAppendToService_Error(t *testing.T) {
	editor, err := New(testdata.TestProto)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := editor.Append.ToService("NonExistentService", testdata.TestRPCContent); err == nil {
		t.Error("ToService() expected error, got nil")
	}
}

func TestAppendToMessage(t *testing.T) {
	editor, err := New(testdata.TestProto)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := editor.Append.ToMessage("GetFooRequest", testdata.TestFieldContent); err != nil {
		t.Fatalf("ToMessage() error = %v", err)
	}

	result := editor.String()

	if !strings.Contains(result, "extra_field") {
		t.Error("result should contain the new field")
	}

	fieldPos := strings.Index(result, "extra_field")
	idFieldPos := strings.Index(result, "string id = 1")
	if fieldPos < idFieldPos {
		t.Error("extra_field should be after id")
	}

	if _, err := New(result); err != nil {
		t.Errorf("result should be valid proto: %v", err)
	}
}

func TestAppendToMessage_Error(t *testing.T) {
	editor, err := New(testdata.TestProto)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := editor.Append.ToMessage("NonExistentMessage", testdata.TestFieldContent); err == nil {
		t.Error("ToMessage() expected error, got nil")
	}
}

func TestAppendToEnum(t *testing.T) {
	editor, err := New(testdata.TestProto)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := editor.Append.ToEnum("Status", testdata.TestEnumValueContent); err != nil {
		t.Fatalf("ToEnum() error = %v", err)
	}

	result := editor.String()

	if !strings.Contains(result, "STATUS_PENDING") {
		t.Error("result should contain the new enum value")
	}

	pendingPos := strings.Index(result, "STATUS_PENDING")
	activePos := strings.Index(result, "STATUS_ACTIVE")
	if pendingPos < activePos {
		t.Error("STATUS_PENDING should be after STATUS_ACTIVE")
	}

	if _, err := New(result); err != nil {
		t.Errorf("result should be valid proto: %v", err)
	}
}

func TestAppendToEnum_Error(t *testing.T) {
	editor, err := New(testdata.TestProto)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := editor.Append.ToEnum("NonExistentEnum", testdata.TestEnumValueContent); err == nil {
		t.Error("ToEnum() expected error, got nil")
	}
}

func TestAppendToFile(t *testing.T) {
	editor, err := New(testdata.TestProto)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := editor.Append.ToFile(testdata.TestMessageContent); err != nil {
		t.Fatalf("ToFile() error = %v", err)
	}
	result := editor.String()

	if !strings.Contains(result, "message CreateBarRequest") {
		t.Error("result should contain CreateBarRequest")
	}

	statusEnumPos := strings.Index(result, "enum Status")
	createBarRequestPos := strings.Index(result, "message CreateBarRequest")
	if createBarRequestPos < statusEnumPos {
		t.Error("new message should be after original content")
	}

	if _, err := New(result); err != nil {
		t.Errorf("result should be valid proto: %v", err)
	}
}

func TestIndentation(t *testing.T) {
	tests := []struct {
		name  string
		proto string
		op    func(*Editor) error
		want  string
	}{
		{
			name: "empty block gets indented",
			proto: `syntax = "proto3";
message Foo {}`,
			op: func(e *Editor) error {
				return e.Append.ToMessage("Foo", "string name = 1;")
			},
			want: `syntax = "proto3";
message Foo {
  string name = 1;
}`,
		},
		{
			name: "closing brace on own line",
			proto: `syntax = "proto3";
service Svc {
	rpc Get(Req) returns (Res) {}
}`,
			op: func(e *Editor) error {
				return e.Append.ToService("Svc", "rpc Create(Req) returns (Res) {}")
			},
			want: `syntax = "proto3";
service Svc {
	rpc Get(Req) returns (Res) {}

	rpc Create(Req) returns (Res) {}
}`,
		},
		{
			name: "over-indented content normalized",
			proto: `syntax = "proto3";
service Svc {
	rpc Get(Req) returns (Res) {}
}`,
			op: func(e *Editor) error {
				return e.Append.ToService("Svc", "\t\t\t\trpc Create(Req) returns (Res) {}")
			},
			want: `syntax = "proto3";
service Svc {
	rpc Get(Req) returns (Res) {}

	rpc Create(Req) returns (Res) {}
}`,
		},
		{
			name: "under-indented content normalized",
			proto: `syntax = "proto3";
service Svc {
	rpc Get(Req) returns (Res) {}
}`,
			op: func(e *Editor) error {
				return e.Append.ToService("Svc", "rpc Create(Req) returns (Res) {}")
			},
			want: `syntax = "proto3";
service Svc {
	rpc Get(Req) returns (Res) {}

	rpc Create(Req) returns (Res) {}
}`,
		},
		{
			name: "nested message appends",
			proto: `syntax = "proto3";
message Foo {
	message Bar {}
}`,
			op: func(e *Editor) error {
				if err := e.Append.ToMessage("Bar", "string a = 1;"); err != nil {
					return err
				}
				return e.Append.ToMessage("Foo", "string b = 2;")
			},
			want: `syntax = "proto3";
message Foo {
	message Bar {
		string a = 1;
	}

	string b = 2;
}`,
		},
		{
			name: "sequential service appends",
			proto: `syntax = "proto3";
service Svc {
	rpc Get(Req) returns (Res) {}
}`,
			op: func(e *Editor) error {
				if err := e.Append.ToService("Svc", "  rpc One(Req) returns (Res) {}"); err != nil {
					return err
				}
				return e.Append.ToService("Svc", "rpc Two(Req) returns (Res) {}")
			},
			want: `syntax = "proto3";
service Svc {
	rpc Get(Req) returns (Res) {}

	rpc One(Req) returns (Res) {}

	rpc Two(Req) returns (Res) {}
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor, err := New(tt.proto)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if err := tt.op(editor); err != nil {
				t.Fatalf("op() error = %v", err)
			}

			if diff := cmp.Diff(tt.want, editor.String()); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}

			if _, err := New(editor.String()); err != nil {
				t.Errorf("result should be valid proto: %v", err)
			}
		})
	}
}
