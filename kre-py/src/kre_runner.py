import abc
import asyncio
import logging
import sys
import time
import traceback
from typing import List

import pymongo
from nats.aio.client import Client as NatsClient
from nats.aio.subscription import Subscription
from nats.js.client import JetStreamContext

from config import Config
from exceptions import ProcessMessagesNotImplemented


class Runner:
    def __init__(self, runner_name: str, config: Config):
        logging.basicConfig(
            level=logging.INFO,
            format="%(asctime)s.%(msecs)03dZ %(levelname)s %(message)s",
            datefmt="%Y-%m-%dT%H:%M:%S",
        )
        logging.Formatter.converter = time.gmtime
        logging.addLevelName(logging.DEBUG, "DEBUG")
        logging.addLevelName(logging.WARNING, "WARN")
        logging.addLevelName(logging.FATAL, "ERROR")
        logging.addLevelName(logging.CRITICAL, "ERROR")

        self.logger = logging.getLogger(runner_name)
        self.loop = asyncio.get_event_loop()
        self.nc: NatsClient = NatsClient()
        self.js: JetStreamContext
        self.config: Config = config
        self.subscription_sids: List[Subscription]
        self.runner_name: str = runner_name
        self.mongo_conn = None

    @staticmethod
    def _get_stream_name(version_id: str, workflow_name: str) -> str:
        return f"{version_id.replace('.', '-')}-{workflow_name}"

    def start(self) -> None:

        """
        Run the python node service in an asyncio loop and also listen to new NATS messages.
        """

        try:
            asyncio.ensure_future(self.connect())
            asyncio.ensure_future(self.process_messages())
            self.loop.run_forever()
        except KeyboardInterrupt:
            self.logger.info("process interrupted")
        finally:
            self.loop.run_until_complete(self.stop())
            self.logger.info("closing loop")
            self.loop.close()

    async def connect(self) -> None:

        """
        Connect to NATS.
        """

        self.logger.info(f"Connecting to MongoDB...")
        self.mongo_conn = pymongo.MongoClient(
            self.config.mongo_uri, socketTimeoutMS=10000, connectTimeoutMS=10000
        )

        self.logger.info(f"Connecting to NATS {self.config.nats_server}...")
        self.js = self.nc.jetstream()
        await self.nc.connect(self.config.nats_server, name=self.runner_name)

    async def stop(self) -> None:

        """
        Stop the NATS connection and the asyncio loop.
        """

        # nats.aio.subscription Subscription
        for sub in self.subscription_sids:
            try:
                await sub.unsubscribe()
            except Exception as err:
                tb = traceback.format_exc()
                self.logger.error(
                    f"Error unsubscribing from the NATS subject {sub.subject}: {err}\n\n{tb}"
                )
                sys.exit(1)

            if not self.nc.is_closed:
                await self.nc.close()
                self.logger.info("NATS connection closed")

        self.loop.stop()

    @abc.abstractmethod
    async def process_messages(self):
        raise ProcessMessagesNotImplemented(Exception)
