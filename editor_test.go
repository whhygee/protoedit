package protoedit

import (
	"strings"
	"testing"

	"github.com/whhygee/protoedit/fixtures"
)

func TestNew(t *testing.T) {
	type args struct {
		proto string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid proto",
			args: args{
				proto: fixtures.TestProto,
			},
		},
		{
			name: "empty proto",
			args: args{
				proto: "",
			},
		},
		{
			name: "minimal proto",
			args: args{
				proto: fixtures.TestProtoMinimal,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			editor, err := New(tt.args.proto)
			if err != nil {
				t.Errorf("New() error = %v", err)
				return
			}
			if editor == nil {
				t.Error("New() returned nil editor")
			}
		})
	}
}

func TestNew_Error(t *testing.T) {
	type args struct {
		proto string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "invalid proto syntax",
			args: args{
				proto: fixtures.TestProtoInvalid,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(tt.args.proto)
			if err == nil {
				t.Error("New() expected error, got nil")
			}
		})
	}
}

func TestAppendToService(t *testing.T) {
	type args struct {
		proto       string
		serviceName string
		content     string
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, result string)
	}{
		{
			name: "appends RPC to end of service",
			args: args{
				proto:       fixtures.TestProto,
				serviceName: "TestService",
				content:     fixtures.TestRPCContent,
			},
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
			name: "preserves all comments",
			args: args{
				proto:       fixtures.TestProtoWithComments,
				serviceName: "MyService",
				content:     "  rpc NewRPC(NewRequest) returns (NewResponse) {}",
			},
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
			t.Parallel()

			editor, err := New(tt.args.proto)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if err := editor.AppendToService(tt.args.serviceName, tt.args.content); err != nil {
				t.Fatalf("AppendToService() error = %v", err)
			}

			tt.check(t, editor.String())
		})
	}
}

func TestAppendToService_Error(t *testing.T) {
	type args struct {
		proto       string
		serviceName string
		content     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "service not found",
			args: args{
				proto:       fixtures.TestProto,
				serviceName: "NonExistentService",
				content:     fixtures.TestRPCContent,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			editor, err := New(tt.args.proto)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if err := editor.AppendToService(tt.args.serviceName, tt.args.content); err == nil {
				t.Error("AppendToService() expected error, got nil")
			}
		})
	}
}

func TestAppendToMessage(t *testing.T) {
	type args struct {
		proto       string
		messageName string
		content     string
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, result string)
	}{
		{
			name: "appends field to end of message",
			args: args{
				proto:       fixtures.TestProto,
				messageName: "GetFooRequest",
				content:     fixtures.TestFieldContent,
			},
			check: func(t *testing.T, result string) {
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			editor, err := New(tt.args.proto)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if err := editor.AppendToMessage(tt.args.messageName, tt.args.content); err != nil {
				t.Fatalf("AppendToMessage() error = %v", err)
			}

			tt.check(t, editor.String())
		})
	}
}

func TestAppendToMessage_Error(t *testing.T) {
	type args struct {
		proto       string
		messageName string
		content     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "message not found",
			args: args{
				proto:       fixtures.TestProto,
				messageName: "NonExistentMessage",
				content:     fixtures.TestFieldContent,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			editor, err := New(tt.args.proto)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if err := editor.AppendToMessage(tt.args.messageName, tt.args.content); err == nil {
				t.Error("AppendToMessage() expected error, got nil")
			}
		})
	}
}

func TestAppendToEnum(t *testing.T) {
	type args struct {
		proto    string
		enumName string
		content  string
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, result string)
	}{
		{
			name: "appends value to end of enum",
			args: args{
				proto:    fixtures.TestProto,
				enumName: "Status",
				content:  fixtures.TestEnumValueContent,
			},
			check: func(t *testing.T, result string) {
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			editor, err := New(tt.args.proto)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if err := editor.AppendToEnum(tt.args.enumName, tt.args.content); err != nil {
				t.Fatalf("AppendToEnum() error = %v", err)
			}

			tt.check(t, editor.String())
		})
	}
}

func TestAppendToEnum_Error(t *testing.T) {
	type args struct {
		proto    string
		enumName string
		content  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "enum not found",
			args: args{
				proto:    fixtures.TestProto,
				enumName: "NonExistentEnum",
				content:  fixtures.TestEnumValueContent,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			editor, err := New(tt.args.proto)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if err := editor.AppendToEnum(tt.args.enumName, tt.args.content); err == nil {
				t.Error("AppendToEnum() expected error, got nil")
			}
		})
	}
}

func TestAppend(t *testing.T) {
	type args struct {
		proto   string
		content string
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, result string)
	}{
		{
			name: "appends message to end of file",
			args: args{
				proto:   fixtures.TestProto,
				content: fixtures.TestMessageContent,
			},
			check: func(t *testing.T, result string) {
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			editor, err := New(tt.args.proto)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			editor.Append(tt.args.content)
			tt.check(t, editor.String())
		})
	}
}

func TestCombinedOperations(t *testing.T) {
	type args struct {
		proto       string
		serviceName string
		rpcContent  string
		msgContent  string
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, result string)
	}{
		{
			name: "add RPC and messages",
			args: args{
				proto:       fixtures.TestProto,
				serviceName: "TestService",
				rpcContent:  fixtures.TestRPCContent,
				msgContent:  fixtures.TestMessagesContent,
			},
			check: func(t *testing.T, result string) {
				wantContains := []string{
					"rpc CreateBar",
					"message CreateBarRequest",
					"message CreateBarResponse",
				}
				for _, want := range wantContains {
					if !strings.Contains(result, want) {
						t.Errorf("result should contain %q", want)
					}
				}

				if _, err := New(result); err != nil {
					t.Errorf("result should be valid proto: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			editor, err := New(tt.args.proto)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if err := editor.AppendToService(tt.args.serviceName, tt.args.rpcContent); err != nil {
				t.Fatalf("AppendToService() error = %v", err)
			}

			editor.Append(tt.args.msgContent)
			tt.check(t, editor.String())
		})
	}
}
