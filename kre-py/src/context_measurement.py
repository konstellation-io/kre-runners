from datetime import datetime

from influxdb_client import InfluxDBClient, Point, WritePrecision
from influxdb_client.client.write_api import WriteOptions, WriteType

# We are using influxdb-client-python that is compatible with 1.8+ versions:
# https://github.com/influxdata/influxdb-client-python#influxdb-1-8-api-compatibility
INFLUX_ORG = ""  # not used
INFLUX_BUCKET = "kre"  # kre is the created database during the deployment
INFLUX_TOKEN = ""  # we don't need authentication


class ContextMeasurement:
    PRECISION_S = WritePrecision.S
    PRECISION_MS = WritePrecision.MS
    PRECISION_US = WritePrecision.US
    PRECISION_NS = WritePrecision.NS

    def __init__(self, config, logger):
        self.__config__ = config
        self.__logger__ = logger

        client = InfluxDBClient(url=self.__config__.influx_uri, token=INFLUX_TOKEN)

        opts = WriteOptions(write_type=WriteType.batching, batch_size=1_000, flush_interval=1_000)
        self.__write_api__ = client.write_api(write_options=opts)

    def save(
        self,
        measurement: str,
        fields: dict,
        tags: dict,
        time: datetime = None,
        precision: str = PRECISION_NS,
    ):
        """
        Save will save a metric into this runtime's influx bucket.
        Default tags will be added:

        'version' - The version's name

        'workflow' - The workflow's name

        'node' - This node's name
        """
        point = Point(measurement)

        for key in fields:
            point.field(key, fields[key])

        for key in tags:
            point.tag(key, tags[key])

        point.tag("version", self.__config__.krt_version)
        point.tag("workflow", self.__config__.krt_workflow_name)
        point.tag("node", self.__config__.krt_node_name)

        if time is None:
            time = datetime.utcnow()

        point.time(time, precision)

        self.__write_api__.write(self.__config__.krt_runtime_id, INFLUX_ORG, point)
