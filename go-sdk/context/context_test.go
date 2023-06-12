package context

import (
	"reflect"
	"testing"

	"github.com/go-logr/logr"
	kai "github.com/konstellation-io/kre-runners/go-sdk/v1/protos"
	"github.com/nats-io/nats.go"
)

func TestNewContext(t *testing.T) {
	type args struct {
		logger       logr.Logger
		natsCli      *nats.Conn
		jetstreamCli nats.JetStreamContext
	}
	tests := []struct {
		name string
		args args
		want KaiContext
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContext(tt.args.logger, tt.args.natsCli, tt.args.jetstreamCli); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_GetRequestID(t *testing.T) {
	type fields struct {
		nats              *nats.Conn
		jetstream         nats.JetStreamContext
		requestMessage    *kai.KaiNatsMessage
		Logger            logr.Logger
		PathUtils         pathUtils
		Metadata          metadata
		Messaging         messaging
		ObjectStore       objectStore
		CentralizedConfig centralizedConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := KaiContext{
				nats:              tt.fields.nats,
				jetstream:         tt.fields.jetstream,
				requestMessage:    tt.fields.requestMessage,
				Logger:            tt.fields.Logger,
				PathUtils:         tt.fields.PathUtils,
				Metadata:          tt.fields.Metadata,
				Messaging:         tt.fields.Messaging,
				ObjectStore:       tt.fields.ObjectStore,
				CentralizedConfig: tt.fields.CentralizedConfig,
			}
			if got := ctx.GetRequestID(); got != tt.want {
				t.Errorf("KaiContext.GetRequestID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShallowCopyWithRequest(t *testing.T) {
	type args struct {
		ctx        KaiContext
		requestMsg *kai.KaiNatsMessage
	}
	tests := []struct {
		name string
		args args
		want KaiContext
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShallowCopyWithRequest(tt.args.ctx, tt.args.requestMsg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShallowCopyWithRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_SetRequest(t *testing.T) {
	type fields struct {
		nats              *nats.Conn
		jetstream         nats.JetStreamContext
		requestMessage    *kai.KaiNatsMessage
		Logger            logr.Logger
		PathUtils         pathUtils
		Metadata          metadata
		Messaging         messaging
		ObjectStore       objectStore
		CentralizedConfig centralizedConfig
	}
	type args struct {
		requestMsg *kai.KaiNatsMessage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := KaiContext{
				nats:              tt.fields.nats,
				jetstream:         tt.fields.jetstream,
				requestMessage:    tt.fields.requestMessage,
				Logger:            tt.fields.Logger,
				PathUtils:         tt.fields.PathUtils,
				Metadata:          tt.fields.Metadata,
				Messaging:         tt.fields.Messaging,
				ObjectStore:       tt.fields.ObjectStore,
				CentralizedConfig: tt.fields.CentralizedConfig,
			}
			ctx.SetRequest(tt.args.requestMsg)
		})
	}
}

func TestContextWrapper_SetRequest(t *testing.T) {
	type fields struct {
		nats              *nats.Conn
		jetstream         nats.JetStreamContext
		requestMessage    *kai.KaiNatsMessage
		Logger            logr.Logger
		PathUtils         pathUtils
		Metadata          metadata
		Messaging         messaging
		ObjectStore       objectStore
		CentralizedConfig centralizedConfig
	}
	type args struct {
		requestMessage *kai.KaiNatsMessage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWrapper{
				nats:              tt.fields.nats,
				jetstream:         tt.fields.jetstream,
				requestMessage:    tt.fields.requestMessage,
				Logger:            tt.fields.Logger,
				PathUtils:         tt.fields.PathUtils,
				Metadata:          tt.fields.Metadata,
				Messaging:         tt.fields.Messaging,
				ObjectStore:       tt.fields.ObjectStore,
				CentralizedConfig: tt.fields.CentralizedConfig,
			}
			ctx.SetRequest(tt.args.requestMessage)
		})
	}
}
