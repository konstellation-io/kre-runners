import asyncio
import logging
import abc

NATS_FLUSH_TIMEOUT = 10


class Runner:
    def __init__(self, runner_name, config):
        logging.basicConfig(
            level=logging.INFO,
            format="%(asctime)s %(levelname)s %(message)s",
            datefmt="%Y-%m-%dT%H:%M:%S%z"
        )
        logging.addLevelName(logging.DEBUG, 'DEBUG')
        logging.addLevelName(logging.WARNING, 'WARN')
        logging.addLevelName(logging.FATAL, 'ERROR')
        logging.addLevelName(logging.CRITICAL, 'ERROR')

        self.logger = logging.getLogger(runner_name)
        self.loop = asyncio.get_event_loop()
        self.config = config

    def run(self):
        try:
            asyncio.ensure_future(self.start())
            self.loop.run_forever()
        except KeyboardInterrupt:
            self.logger.info("Process interrupted")
        finally:
            self.loop.run_until_complete(self.stop())
            self.loop.close()

    @abc.abstractmethod
    async def stop(self):
        raise Exception(f"Stop should be implemented.")

    @abc.abstractmethod
    async def start(self):
        raise Exception(f"Start should be implemented.")
