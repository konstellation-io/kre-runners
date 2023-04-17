import json

from nats.aio.client import Client as NatsClient
from nats.errors import TimeoutError
from nats.js.client import JetStreamContext
from config import Config

from compression import compress_if_needed


class ContextData:
    def __init__(self, config, nc, js, mongo_conn, logger):
        self.__config__: Config = config
        self.__nc__: NatsClient = nc
        self.__js__: JetStreamContext = js
        self.__mongo_conn__ = mongo_conn
        self.__logger__ = logger

    async def save(self, coll, data):
        if not isinstance(coll, str) or coll == "":
            raise Exception(
                f"[ctx.db.save] invalid 'collection'='{coll}', must be a nonempty string"
            )

        stream_info = await self.__js__.stream_info(self.__config__.nats_stream)
        stream_max_size = stream_info.config.max_msg_size or -1
        server_max_size = self.__nc__.max_payload

        max_size = (
            min(stream_max_size, server_max_size) if stream_max_size != -1 else server_max_size
        )

        try:
            subject = self.__config__.nats_mongo_writer
            payload = bytes(json.dumps({"coll": coll, "doc": data}), encoding="utf-8")
            out = compress_if_needed(payload, max_size=max_size, logger=self.__logger__)
            response = await self.__nc__.request(subject, out, timeout=60)
            res_json = json.loads(response.data.decode())
            if not res_json["success"]:
                self.__logger__.error("Unexpected error saving data")
        except TimeoutError:
            self.__logger__.error("Error saving data: request timed out")

    async def find(self, coll, query):
        if not isinstance(coll, str) or coll == "":
            raise Exception(
                f"[ctx.db.save] invalid 'collection'='{coll}', must be a nonempty string"
            )

        if not isinstance(query, dict) or not query:
            raise Exception(f"[ctx.db.find] invalid 'query'='{query}', must be a nonempty dict")

        try:
            collection = self.__mongo_conn__[self.__config__.mongo_data_db_name][coll]
            self.__logger__.debug(
                f"call to mongo to get data on{self.__config__.mongo_data_db_name}.{coll}: {query}"
            )
            cursor = collection.find(query)
            return list(cursor)

        except Exception as err:
            raise Exception(f"[ctx.db.find] error getting data from MongoDB: {err}")
