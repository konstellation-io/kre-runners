from typing import NewType, Any, Dict
from collections.abc import Callable
from context import HandlerContext

Handler = NewType("Handler", Callable[[HandlerContext, Any]])


class HandlerManager:
    def __init__(self, default_handler: Handler, custom_handlers: Dict[str, Handler]):
        self.default_handler = default_handler
        self.custom_handlers = custom_handlers

    def get_handler(self, node: str) -> Handler:
        try:
            handler = self.custom_handlers.get(node)
            if handler is not None:
                return handler
            return self.default_handler
        except Exception:
            return self.default_handler
