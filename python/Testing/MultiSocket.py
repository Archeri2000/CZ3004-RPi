import asyncio


async def handle_echo(reader: asyncio.StreamReader, writer: asyncio.StreamWriter):
    try:
        while True:
            data = await reader.read(100)
            if not data:
                break
            message = data.decode()
            addr = writer.get_extra_info('peername')

            print(f"Received {message!r} from {addr!r}")

            print(f"Send: {message!r}")
            writer.write(data)
            await writer.drain()
        print("Close the connection")
        writer.close()
    except Exception as e:
        print(f"Exception: {e}")


async def serve_echo_server(host: str = '127.0.0.1', port: int = 8888) -> asyncio.tasks:
    server = await asyncio.start_server(
        handle_echo, host, port)

    addr = server.sockets[0].getsockname()
    print(f'Serving on {addr}')
    async with server:
        await server.serve_forever()


async def main():
    s1 = asyncio.create_task(serve_echo_server(port=8080))
    s2 = asyncio.create_task(serve_echo_server(port=8081))
    await s1
    await s2

asyncio.run(main())