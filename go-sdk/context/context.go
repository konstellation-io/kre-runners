package context

import (
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/context/centralized_configuration"
	msg "github.com/konstellation-io/kre-runners/go-sdk/v1/context/messaging"
	meta "github.com/konstellation-io/kre-runners/go-sdk/v1/context/metadata"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/context/object_store"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/context/path_utils"
	kai "github.com/konstellation-io/kre-runners/go-sdk/v1/protos"
	"os"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type pathUtils interface {
	GetBasePath() string
	ComposeBasePath(relativePath string) string
}

type messaging interface {
	SendOutput(response proto.Message, channelOpt ...string) error
	SendOutputWithRequestID(response proto.Message, requestID string, channelOpt ...string) error
	SendAny(response *anypb.Any, channelOpt ...string)
	SendEarlyReply(response proto.Message, channelOpt ...string) error
	SendEarlyExit(response proto.Message, channelOpt ...string) error

	IsMessageOK() bool
	IsMessageError() bool
	IsMessageEarlyReply() bool
	IsMessageEarlyExit() bool
}

type metadata interface {
	GetProcess() string
	GetWorkflow() string
	GetProduct() string
	GetVersion() string
	GetObjectStoreName() string
	GetKeyValueStoreProductName() string
	GetKeyValueStoreWorkflowName() string
	GetKeyValueStoreProcessName() string
}

type objectStore interface {
	Get(key string) ([]byte, error)
	Save(key string, value []byte) error
	Delete(key string) error
}

type centralizedConfig interface {
	GetConfig(key string, scope ...msg.Scope) (string, error)
	SetConfig(key, value string, scope ...msg.Scope) error
	DeleteConfig(key string, scope msg.Scope) (string, error)
}

//TODO add metrics interface

//TODO add storage interface

//TODO add MongoDB interface

type KaiContext struct {
	// Needed deps
	nats      *nats.Conn
	jetstream nats.JetStreamContext
	//repository     mongodb.Manager
	//measurements   influxdb2.Client
	requestMessage *kai.KaiNatsMessage

	// Main methods
	Logger            logr.Logger
	PathUtils         pathUtils
	Metadata          metadata
	Messaging         messaging
	ObjectStore       objectStore
	CentralizedConfig centralizedConfig
}

type ContextWrapper KaiContext

func NewContext(logger logr.Logger, natsCli *nats.Conn, jetstreamCli nats.JetStreamContext) KaiContext {
	centralizedConfigInst, err := centralized_configuration.NewCentralizedConfiguration(logger, jetstreamCli)
	if err != nil {
		logger.Error(err, "Error initializing Centralized Configuration")
		os.Exit(1)
	}

	objectStoreInst, err := object_store.NewObjectStore(logger, jetstreamCli)
	if err != nil {
		logger.Error(err, "Error initializing Object Store")
		os.Exit(1)
	}

	messagingInst := msg.NewMessaging(logger, natsCli, jetstreamCli, nil)

	ctx := KaiContext{
		nats:              natsCli,
		jetstream:         jetstreamCli,
		Logger:            logger,
		PathUtils:         path_utils.NewPathUtils(logger),
		Metadata:          meta.NewMetadata(logger),
		Messaging:         messagingInst,
		ObjectStore:       objectStoreInst,
		CentralizedConfig: centralizedConfigInst,
	}

	return ctx
}

func (ctx KaiContext) GetRequestID() string {
	return ctx.requestMessage.RequestId
}

func ShallowCopyWithRequest(ctx KaiContext, requestMsg *kai.KaiNatsMessage) KaiContext {
	hCtx := ctx
	hCtx.requestMessage = requestMsg

	return hCtx
}

func (ctx KaiContext) SetRequest(requestMsg *kai.KaiNatsMessage) {
	ctx.requestMessage = requestMsg
}

func (ctx ContextWrapper) SetRequest(requestMessage *kai.KaiNatsMessage) {
	ctx.requestMessage = requestMessage
}
