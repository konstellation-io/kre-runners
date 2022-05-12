import abc
import asyncio
import logging
import time

from nats.aio.client import Client as NATS


class Runner:
    def __init__(self, runner_name, config):
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
        self.nc = NATS()
        self.js = self.nc.jetstream()
        self.config = config
        self.subscription_sid = None
        self.runner_name = runner_name

    def _get_stream_name(self, version_id: str, workflow_name: str):
        return f"{version_id.replace('.', '-')}-{workflow_name}"

    def start(self):
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

    async def connect(self):
        self.logger.info(f"Connecting to NATS {self.config.nats_server}...")
        await self.nc.connect(
            self.config.nats_server, name=self.runner_name
        )
        
        #subjects = [self.config.nats_input]

        #if self.config.nats_output != "":
        #    subjects.append(self.config.nats_output)

        #self.logger.info(f"Add stream {self.config.nats_stream} subjects: {subjects}")
        #await self.js.add_stream(name=self.config.nats_stream, subjects=subjects)

    async def stop(self):
        if not self.nc.is_closed:
            self.logger.info("closing NATS connection")
            await self.nc.close()

        # await self.js.delete_stream(name=self.runner_name)
        self.logger.info("stop loop")
        self.loop.stop()

    @abc.abstractmethod
    async def process_messages(self):
        raise Exception(f"process_messages should be implemented.")
