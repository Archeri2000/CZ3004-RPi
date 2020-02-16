import asyncio
from typing import List

import serial_asyncio

ArduinoWriter: List[asyncio.StreamWriter] = []


def processData(raw_data: bytes) -> int:
    # TODO: Write conversion from raw to number of squares
    return 0


async def handle_serial(algo_buffer: List[asyncio.StreamWriter], reader: asyncio.StreamReader, writer: asyncio.StreamWriter):
    # TODO: Confirm endian of integer transmission.
    ArduinoWriter.append(writer)
    try:
        while True:
            # Read sensor data
            raw_data = await reader.readline()
            if not raw_data:
                break
            sensor_data = processData(raw_data)

            # Forward sensor data to Algo
            if algo_buffer is not None and not algo_buffer[0].is_closing():
                algo_buffer[0].write(bytes(sensor_data))
                await algo_buffer[0].drain()
            else:
                raise Exception("PC not connected!")

        print("Close the connection")
        writer.close()
    except Exception as e:
        print(f"Exception: {e}")


async def serve_serial_conn(algo_buffer: List[asyncio.StreamWriter],
                            serialURL: str = '/dev/tty', baud: int = 115200) -> asyncio.tasks:
    reader, writer = await serial_asyncio.open_serial_connection(url=serialURL, baudrate=baud)
    await handle_serial(algo_buffer, reader, writer)

