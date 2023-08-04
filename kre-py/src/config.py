import os


class Config:
    def __init__(self) -> None:
        self.js_request_timeout = int(os.getenv("KRT_JS_REQUEST_TIMEOUT", 30000))

        # Mandatory environment variables
        try:
            self.krt_workflow_name = os.environ["KRT_WORKFLOW_NAME"]
            self.krt_runtime_id = os.environ["KRT_RUNTIME_ID"]
            self.krt_version_id = os.environ["KRT_VERSION_ID"]
            self.krt_version = os.environ["KRT_VERSION"]
            self.krt_node_name = os.environ["KRT_NODE_NAME"]
            self.nats_server = os.environ["KRT_NATS_SERVER"]
            self.nats_inputs = os.environ["KRT_NATS_INPUTS"].split(",")
            self.nats_output = os.environ["KRT_NATS_OUTPUT"]
            self.nats_stream = os.environ["KRT_NATS_STREAM"]
            self.nats_mongo_writer = os.environ["KRT_NATS_MONGO_WRITER"]
            self.nats_object_store = os.getenv("KRT_NATS_OBJECT_STORE", default=None)
            self.nats_key_value_store_project = os.environ["KRT_NATS_KEY_VALUE_STORE_PROJECT"]
            self.nats_key_value_store_workflow = os.environ["KRT_NATS_KEY_VALUE_STORE_WORKFLOW"]
            self.nats_key_value_store_node = os.environ["KRT_NATS_KEY_VALUE_STORE_NODE"]
            self.base_path = os.environ["KRT_BASE_PATH"]
            self.handler_path = os.environ["KRT_HANDLER_PATH"]
            self.mongo_data_db_name = "data"
            self.mongo_uri = os.environ["KRT_MONGO_URI"]
            self.influx_uri = os.environ["KRT_INFLUX_URI"]
        except Exception as err:
            raise Exception(f"error reading config: the {str(err)} env var is missing")
