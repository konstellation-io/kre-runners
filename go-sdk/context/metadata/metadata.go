package metadata

import (
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

type Metadata struct {
	logger logr.Logger
}

func NewMetadata(logger logr.Logger) *Metadata {
	return &Metadata{
		logger: logger,
	}
}

func (md Metadata) GetProcess() string {
	return viper.GetString("krt_process_name")
}

func (md Metadata) GetWorkflow() string {
	return viper.GetString("krt_workflow_name")
}

func (md Metadata) GetProduct() string {
	return viper.GetString("krt_product_name")
}

func (md Metadata) GetVersion() string {
	return viper.GetString("krt_version")
}

func (md Metadata) GetObjectStoreName() string {
	return viper.GetString("krt_nats_object_store")
}

func (md Metadata) GetKeyValueStoreProductName() string {
	return viper.GetString("krt_nats_key_value_store_product")
}

func (md Metadata) GetKeyValueStoreWorkflowName() string {
	return viper.GetString("krt_nats_key_value_store_workflow")
}

func (md Metadata) GetKeyValueStoreProcessName() string {
	return viper.GetString("krt_nats_key_value_store_process")
}
