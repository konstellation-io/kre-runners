import logging
from unittest.mock import MagicMock

import pytest


@pytest.fixture
def simple_logger():
    logging.basicConfig(level=logging.DEBUG)
    logging.getLogger().setLevel(logging.DEBUG)
    return logging.getLogger("Tests")


@pytest.fixture
def config():
    config = MagicMock()
    config.krt_version = "test_version"
    config.krt_runtime_id = "test"
    config.krt_workflow_name = "test_workflow"
    config.krt_node_name = "test_node"
    config.influx_uri = "influx"
    config.nats_mongo_writer = "mongo_writer"
    config.mongo_data_db_name = "test_db_name"
    return config
