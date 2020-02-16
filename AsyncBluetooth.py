from asyncio import *

from bluetooth import *


class bluetoothComm(object):
    # initialize client & server to none and connection to false.
    def __init__(self, callback):
        self.server_sock: BluetoothSocket = None
        self.client_sock: BluetoothSocket = None
        self.bluetooth_connected = False
        self.callback = callback

    # configure bluetooth port and listen for bluetooth connection
    async def init_bluetooth_comm(self):
        while True:
            try:
                self.server_sock = BluetoothSocket(RFCOMM)
                self.server_sock.bind(("", PORT_ANY))
                self.server_sock.listen(1)
                self.server_sock.setblocking(0)
                self.port = self.server_sock.getsockname()[1]

                uuid = "00001101-0000-1000-8000-00805F9B34FB"

                advertise_service(self.server_sock, "MDPGrp4_Server", service_id=uuid,
                                  service_classes=[uuid, SERIAL_PORT_CLASS],
                                  profiles=[SERIAL_PORT_PROFILE]
                                  )

                print("Waiting for connection on RFCOMM channel %d" % self.port)
                while self.client_sock is None:
                    (self.client_sock, client_info) = await self.start_connection()
                print("Accepted connection from " + str(client_info))
                self.bluetooth_connected = True

                # uncomment if testing one-to-one communication
                """
                while True:
                     data = self.reading_from_bluetooth()
                     if(data != None):
                         print(data)
                         sendData = data
                         self.writing_to_bluetooth(sendData)
                         retry = False
                """
                return

            except Exception:
                print("Bluetooth-Connection Error")

    async def start_connection(self):
        await sleep(0)
        return self.server_sock.accept()

    def bluetooth_is_connected(self):
        return self.bluetooth_connected

    async def serve_forever(self):
        while True:
            await self.reading_from_bluetooth()

    async def get_data(self):
        await sleep(0)
        return self.client_sock.recv(1024)

    async def reading_from_bluetooth(self):
        try:
            # Receive data from the socket with 1024 bytes
            while True:
                data = await self.get_data()
                if not data:
                    break

            # strData = "RPI: " + data.decode('utf-8')
            # rstrip() method returns a copy of the string with trailing characters removed
            print("Read From Bluetooth: " + data.decode("ascii"))
            self.callback(data)
            return data.decode("ascii")

        except BluetoothError as e:
            print("Bluetooth-COMM/Bluetooth-RECV Error: ")
            if ('Connection reset by peer' in str(e)):
                self.disconnect_bluetooth()
                print("Assuming Bluetooth is disconnected, attempting to reconnect...")
                self.init_bluetooth_comm()

    def writing_to_bluetooth(self, message):
        try:
            if not self.bluetooth_connected:
                print("Bluetooth not connected! Unable to transmit")
                return

            # strMessage = "RPI: " + message.decode('utf-8')
            strMessage = message.encode('ascii')
            self.client_sock.send(strMessage)
            message1 = strMessage.decode('ascii')
            print("Writing to Bluetooth: " + message1)

        except BluetoothError as e:
            print("Bluetooth-SEND Error: " + str(e))

    def disconnect_bluetooth(self):
        try:
            if not (self.server_socket is None):
                self.server_socket.close()
                print("Closing Bluetooth Server Socket...")

            if not (self.client_sock is None):
                self.client_sock.close()
                print("Closing Bluetooth Client Socket...")

        except Exception:
            pass
        self.bluetooth_connected = False
        print("ALL DONE!!!")



# uncomment if testing one-to-one communication
"""
bluetoothObj = bluetoothComm()
bluetoothObj.init_bluetooth_comm()
"""

