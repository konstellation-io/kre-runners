import logging
import time
import pytest

from datetime import datetime
from unittest.mock import MagicMock, ANY

from influxdb_client import Point

from config import Config
from context_measurement import ContextMeasurement


@pytest.fixture
def simple_logger():
    return logging.getLogger("Tests")


@pytest.fixture
def config():
    config = MagicMock(Config)
    config.krt_version = "test_version"
    config.krt_runtime_id = "test"
    config.krt_workflow_name = "test_workflow"
    config.krt_node_name = "test_node"
    config.influx_uri = "influx"
    return config


@pytest.mark.unittest
def test_measurement_save(simple_logger, config):
    context_measurement = ContextMeasurement(config, simple_logger)
    context_measurement.__write_api__ = MagicMock()
    context_measurement.__write_api__.write = MagicMock()

    # Create new test Point
    measurement = "test_measurement"
    fields = {"field": "test"}
    tags = {"tag": "test"}
    timestamp = datetime.strptime('2018-06-29 08:15:27.243860', '%Y-%m-%d %H:%M:%S.%f')

    point = Point(measurement)
    point.field("field", "test")
    point.tag("tag", "test")
    point.tag("version", config.krt_version)
    point.tag("workflow", config.krt_workflow_name)
    point.tag("node", config.krt_node_name)
    point.time(time, timestamp)

    context_measurement.save(measurement, fields, tags, timestamp)

    context_measurement.__write_api__.write.assert_called_once()
    context_measurement.__write_api__.write.assert_called_with(config.krt_runtime_id, "", ANY)
    args = context_measurement.__write_api__.write.call_args.args
    assert len(args) == 3
    assert args[2]._name == measurement
    assert args[2]._time == timestamp
    assert args[2]._tags['tag'] == 'test'
    assert args[2]._tags['version'] == config.krt_version
    assert args[2]._tags['workflow'] == config.krt_workflow_name
    assert args[2]._tags['node'] == config.krt_node_name
    assert args[2]._fields['field'] == 'test'
