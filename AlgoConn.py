import asyncio
from typing import List

AlgoWriter: List[asyncio.StreamWriter] = []


def decodeInstruction(inst: int, dist: int) -> bytes:
    # TODO: Implement
    return b'0'


def generatePacketHandler(android_buffer: List[asyncio.StreamWriter], arduino_buffer: List[asyncio.StreamWriter]):
    async def handle_packet(reader: asyncio.StreamReader, writer: asyncio.StreamWriter):
        # TODO: Confirm endian of integer transmission.
        AlgoWriter.append(writer)
        try:
            while True:
                # Read instruction
                raw_instruction = await reader.read(1)
                if not raw_instruction:
                    break
                inst = int.from_bytes(raw_instruction, byteorder='big')
                # Read distance
                raw_dist = await reader.read(4)
                if not raw_dist:
                    break
                dist = int.from_bytes(raw_dist, byteorder='big')

                # Decode next instruction for arduino and send
                next_inst = decodeInstruction(inst, dist)
                if len(arduino_buffer) == 1 and not arduino_buffer[0].is_closing():
                    arduino_buffer[0].write(next_inst)
                    await arduino_buffer[0].drain()
                else:
                    raise Exception("Arduino not connected!")

                # Forward map to android
                mapdata = await reader.readline()
                if not mapdata:
                    break
                if len(android_buffer) == 1 and not android_buffer[0].is_closing():
                    android_buffer[0].write(mapdata)
                    await android_buffer[0].drain()
                else:
                    raise Exception("Android not connected!")

            print("Close the connection")
            writer.close()
        except Exception as e:
            print(f"Exception: {e}")

    return handle_packet


async def serve_algo_server(android_buffer: List[asyncio.StreamWriter], arduino_buffer: List[asyncio.StreamWriter],
                            host: str = '127.0.0.1', port: int = 8888) -> asyncio.tasks:
    server = await asyncio.start_server(
        generatePacketHandler(android_buffer, arduino_buffer), host, port)
    addr = server.sockets[0].getsockname()
    print(f'Serving on {addr}')
    async with server:
        await server.serve_forever()

