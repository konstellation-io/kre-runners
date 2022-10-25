import logging
import time
import pytest

from datetime import datetime
from unittest.mock import MagicMock, ANY

from freezegun import freeze_time
from influxdb_client import Point

from context_measurement import ContextMeasurement

DATETIME_FORMAT = '%Y-%m-%d %H:%M:%S.%f'


@pytest.fixture
def simple_logger():
    return logging.getLogger("Tests")


@pytest.fixture
def config():
    config = MagicMock()
    config.krt_version = "test_version"
    config.krt_runtime_id = "test"
    config.krt_workflow_name = "test_workflow"
    config.krt_node_name = "test_node"
    config.influx_uri = "influx"
    return config


@pytest.fixture
def context_measurement(simple_logger, config):
    context_measurement = ContextMeasurement(config, simple_logger)
    context_measurement.__write_api__ = MagicMock()
    context_measurement.__write_api__.write = MagicMock()
    return context_measurement


@pytest.mark.unittest
def test_measurement_save_with_custom_timestamp_expect_ok(context_measurement, config):
    # Create new test Point
    measurement = "test_measurement"
    fields = {"field1": "test1", "field2": "test2"}
    tags = {"tag1": "test1", "tag2": "test2"}
    timestamp = datetime.strptime('2018-06-29 08:15:27.243860', DATETIME_FORMAT)

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
    assert args[2]._write_precision == ContextMeasurement.PRECISION_NS
    assert args[2]._tags['tag1'] == 'test1'
    assert args[2]._tags['tag2'] == 'test2'
    assert args[2]._tags['version'] == config.krt_version
    assert args[2]._tags['workflow'] == config.krt_workflow_name
    assert args[2]._tags['node'] == config.krt_node_name
    assert args[2]._fields['field1'] == 'test1'
    assert args[2]._fields['field2'] == 'test2'


@pytest.mark.unittest
@freeze_time("2020-01-01 01:02:03.000004")
def test_measurement_save_with_custom_timestamp_precision_expect_ok(context_measurement, config):
    # Create new test Point
    measurement = "test_measurement"
    fields = {"field1": "test1", "field2": "test2"}
    tags = {"tag1": "test1", "tag2": "test2"}

    point = Point(measurement)
    point.field("field", "test")
    point.tag("tag", "test")
    point.tag("version", config.krt_version)
    point.tag("workflow", config.krt_workflow_name)
    point.tag("node", config.krt_node_name)

    context_measurement.save(measurement, fields, tags, precision=ContextMeasurement.PRECISION_S)

    context_measurement.__write_api__.write.assert_called_once()
    context_measurement.__write_api__.write.assert_called_with(config.krt_runtime_id, "", ANY)
    args = context_measurement.__write_api__.write.call_args.args
    assert len(args) == 3
    assert args[2]._name == measurement
    assert args[2]._time == datetime.strptime("2020-01-01 01:02:03.000004", DATETIME_FORMAT)
    assert args[2]._write_precision == ContextMeasurement.PRECISION_S
    assert args[2]._tags['tag1'] == 'test1'
    assert args[2]._tags['tag2'] == 'test2'
    assert args[2]._tags['version'] == config.krt_version
    assert args[2]._tags['workflow'] == config.krt_workflow_name
    assert args[2]._tags['node'] == config.krt_node_name
    assert args[2]._fields['field1'] == 'test1'
    assert args[2]._fields['field2'] == 'test2'


@pytest.mark.unittest
@freeze_time("2020-01-01 00:00:00.000000")
def test_measurement_save_with_default_timestamp_expect_ok(context_measurement, config):
    # Create new test Point
    measurement = "test_measurement"
    fields = {"field1": "test1", "field2": "test2"}
    tags = {"tag1": "test1", "tag2": "test2"}

    point = Point(measurement)
    point.field("field", "test")
    point.tag("tag", "test")
    point.tag("version", config.krt_version)
    point.tag("workflow", config.krt_workflow_name)
    point.tag("node", config.krt_node_name)

    context_measurement.save(measurement, fields, tags)

    context_measurement.__write_api__.write.assert_called_once()
    context_measurement.__write_api__.write.assert_called_with(config.krt_runtime_id, "", ANY)
    args = context_measurement.__write_api__.write.call_args.args
    assert len(args) == 3
    assert args[2]._name == measurement
    assert args[2]._time == datetime.strptime("2020-01-01 00:00:00.000000", DATETIME_FORMAT)
    assert args[2]._write_precision == ContextMeasurement.PRECISION_NS
    assert args[2]._tags['tag1'] == 'test1'
    assert args[2]._tags['tag2'] == 'test2'
    assert args[2]._tags['version'] == config.krt_version
    assert args[2]._tags['workflow'] == config.krt_workflow_name
    assert args[2]._tags['node'] == config.krt_node_name
    assert args[2]._fields['field1'] == 'test1'
    assert args[2]._fields['field2'] == 'test2'
