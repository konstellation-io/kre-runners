class CompressedMessageTooLargeException(Exception):
    def __init__(self, max_size):
        self.max_size = max_size
        message = f"Compressed message exceeds maximum size allowed: {max_size} bytes"
        super().__init__(message)