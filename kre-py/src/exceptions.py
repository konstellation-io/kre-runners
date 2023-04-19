class CompressedMessageTooLargeException(Exception):
    def __init__(self, message):
        # Call the base class constructor with the parameters it needs
        super().__init__(message)


class ProcessMessagesNotImplemented(Exception):
    """Exception thrown when the process_messages method is not implemented."""

    pass
