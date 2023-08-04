import abc
import asyncio
import logging
import sys
import time
import traceback
import pymongo

from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext

from config import Config
from exceptions import ProcessMessagesNotImplemented
from context_configuration import ContextConfiguration, new_context_configuration


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
        self.loop: asyncio.AbstractEventLoop
        self.nc: NatsClient = NatsClient()
        self.js: JetStreamContext
        self.config: Config = config
        self.subscription_sids = []
        self.runner_name: str = runner_name
        self.mongo_conn = None
        self.configuration: ContextConfiguration = None

    @staticmethod
    def _get_stream_name(version_id: str, workflow_name: str) -> str:
        return f"{version_id.replace('.', '-')}-{workflow_name}"

    async def connect(self) -> None:
        """
        Connect to NATS.
        """

        self.logger.info(f"Connecting to MongoDB...")
        self.mongo_conn = pymongo.MongoClient(
            self.config.mongo_uri, socketTimeoutMS=10000, connectTimeoutMS=10000
        )

        self.logger.info(f"Connecting to NATS {self.config.nats_server}...")
        self.js = self.nc.jetstream(timeout=self.config.js_request_timeout)
        await self.nc.connect(self.config.nats_server, name=self.runner_name)

        self.logger.info(f"Creating context configuration...")
        self.configuration = await new_context_configuration(self.config, self.logger, self.js)

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
