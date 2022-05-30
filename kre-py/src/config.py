import os

NATS_FLUSH_TIMEOUT = 10


class Config:
    def __init__(self):
        # Mandatory variables
        try:
            self.krt_workflow_name = os.environ.get('KRT_WORKFLOW_NAME', "workflow1")
            self.krt_runtime_id = os.environ.get('KRT_RUNTIME_ID', "runtime1")
            self.krt_version_id = os.environ.get('KRT_VERSION_ID', "version.1234")
            self.krt_version = os.environ.get('KRT_VERSION', "testVersion1")
            self.krt_node_name = os.environ.get('KRT_NODE_NAME', "nodeA")
            self.nats_server = os.environ.get('KRT_NATS_SERVER', "nats:4222")
            self.nats_input = os.environ.get('KRT_NATS_INPUT', "test-subject-input")
            self.nats_output = os.environ.get('KRT_NATS_OUTPUT', "")
            self.nats_mongo_writer = os.environ.get('KRT_NATS_MONGO_WRITER', "mongo_writer")
            self.base_path = os.environ.get('KRT_BASE_PATH', "/tmp/myvol")
            self.handler_path = os.environ.get('KRT_HANDLER_PATH', "src/node/node_handler.py")
            self.mongo_data_db_name = "data"
            self.mongo_uri = os.environ.get('KRT_MONGO_URI', "mongodb://admin:123456@mongo:27017/admin")
            self.influx_uri = os.environ.get('KRT_INFLUX_URI', "http://influx:8086")
            self.nats_flush_timeout = NATS_FLUSH_TIMEOUT

        except Exception as err:
            raise Exception(f"error reading config: the {str(err)} env var is missing")
