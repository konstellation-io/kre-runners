import gzip
import logging

from exceptions import CompressedMessageTooLargeException

# The compresslevel argument is an integer from 0 to 9 controlling the level of compression;
# 1 is fastest and produces the least compression, and 9 is slowest and produces the most compression.
# 0 is no compression.
COMPRESS_LEVEL = 9

# This is a magic number in the file header that compressed files contains.
GZIP_HEADER = b"\x1f\x8b"


def bytes_to_kb(size_bytes: int) -> str:
    """Transform bytes to KiloBytes.

    Args:
        size_bytes (int): amount of bytes.

    Returns:
        str: Returns an string indicating the size of the input data.
    """
    return f"{(size_bytes / 1024):.2f} KB"


def size_in_kb(data: bytes) -> str:
    """Gets the data size in KB.

    Args:
        data (bytes): any input data.

    Returns:
        str: Returns an string indicating the size of the input data.
    """
    return bytes_to_kb(len(data))


def is_compressed(data: bytes) -> bool:
    """Compares the first two bytes to check if the data is compressed.

    Args:
        data (bytes): any input data.

    Returns:
        bool: If the data is compressed returns true.
    """
    return data.startswith(GZIP_HEADER)


def compress(data: bytes) -> bytes:
    """Compresses the input data using GZIP.

    Args:
        data (bytes): any input data.

    Returns:
        bytes: Returns the compressed data.
    """
    return gzip.compress(data, compresslevel=COMPRESS_LEVEL)


def uncompress(data: bytes) -> bytes:
    """Uncompresses the input data.

    Args:
        data (bytes): any input data.

    Returns:
        bytes: Returns the uncompressed data.
    """
    return gzip.decompress(data)


def compress_if_needed(
    data: bytes,
    max_size: int,
    logger: logging.Logger = logging.getLogger(),
) -> bytes:
    """If the msg is bigger than the allowed max_size, compresses the msg.
    In other case returns the data without any modifications.

    Args:
        data (bytes): Any input data.
        logger (logging.Logger, optional): A logger instance. Defaults to logging.getLogger().
        max_size (int, optional): The maximum size allowed in bytes. Defaults to MESSAGE_THRESHOLD.

    Raises:
        Exception: If the compresses msg is still bigger than 1MB throws an exception.

    Returns:
        bytes: The output message ensuring that the size is lower than max_size.
    """
    if len(data) <= max_size:
        return data

    out = compress(data)

    if len(out) > max_size:
        data_size_mb = bytes_to_mb(len(data))
        max_size_mb = bytes_to_mb(len(max_size))
        logger.debug("compressed message exceeds maximum size allowed: current" +
                    f"message size {data_size_mb}MB, max allowed size {max_size_mb}MB")

        raise CompressedMessageTooLargeException("compressed message exceeds maximum size allowed")

    logger.info(f"Original message size: {size_in_kb(data)}. Compressed{size_in_kb(out)}")

    return out


def bytes_to_mb(size_in_bytes: int) -> float:
    return float("{:.1f}".format(size_in_bytes / 1024 / 1024))
