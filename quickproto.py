
import base64
import threading
import socket

import os

class DataMissingError(Exception):
    pass

class DelimiterMixin:
    def __init__(self, delimiter):
        self.DELIMITER = delimiter

    @property
    def DELIMITER_HEADER(self) -> bytes:
        return self.DELIMITER + self.DELIMITER

    @property
    def DELIMITER_FILE(self) -> bytes:
        return self.DELIMITER_BODY + self.DELIMITER_HEADER

    @property
    def DELIMITER_BODY(self) -> bytes:
        return self.DELIMITER_HEADER + self.DELIMITER_HEADER
    
    @property
    def DELIMITER_END(self) -> bytes:
        return self.DELIMITER_BODY + self.DELIMITER_BODY

class MessageFile:
    def __init__(self, name, data):
        self.name: str = name
        self.data: bytes = data

    def __str__(self):
        return str(self.name)

    def __repr__(self):
        return f"MessageFile({self.name}, has_data: {not not self.data})"
        
class Message(DelimiterMixin):
    HEADERS: dict[str, list[str]]
    BODY: bytes
    _data: bytes
    DELIMITER: bytes
    use_B64: bool
    FILES = dict[str, MessageFile]

    def __init__(self, headers: dict[str, list[str]] = None, body: bytes = None, files={}, delimiter: bytes=b"$", use_b64: bool=False, raw: bytes=None) -> None:
        super().__init__(delimiter)
        self.HEADERS = headers
        self.BODY = body
        self._data = raw
        self.use_B64 = use_b64
        self.FILES = files

    def AddHeader(self, key: str, value: str) -> None:
        self.HEADERS[key].append(value)

    def GetHeader(self, key: str) -> list[str]:
        return self.HEADERS[key]

    def AddFile(self, name: str, data: bytes, mfile: MessageFile = None):
        if mfile is None:
            mfile = MessageFile(name, data)
            self.FILES[name] = mfile
        else:
            self.FILES[name] = mfile

    def Parse(self):
        if self._data is None:
            raise DataMissingError("No data to parse")
        headers, body = self._data.split(self.DELIMITER_BODY, 1)
        header_dict: dict[str, list[str]] = {}
        for header in headers.split(self.DELIMITER_HEADER):
            data = header.split(self.DELIMITER)
            curr_key = data[0].decode('utf-8')
            header_dict[curr_key] = []
            for value in data[1:]:
                header_dict[curr_key].append(value.decode('utf-8'))
        body = body[:-len(self.DELIMITER_END)]
        if self.use_B64:
            body = base64.b64decode(body)

        # Split files from request
        data = body.split(self.DELIMITER_FILE)
        body = data[len(data)-1]
        files = data[:len(data)-1]
        for file in files:
            name, data = file.split(self.DELIMITER_HEADER, 1)
            fdata = base64.b64decode(data)
            fname = name.decode('utf-8')
            self.FILES[fname] = MessageFile(fname, fdata)

        self.HEADERS = header_dict
        self.BODY = body
        return self.HEADERS, self.BODY

    def Generate(self):
        if self.HEADERS is None:
            raise DataMissingError("Headers missing")
        delim = self.DELIMITER
        data = b""
        for key, values in self.HEADERS.items():
            data += key.encode('utf-8') + delim
            for value in values:
                data += value.encode('utf-8') + delim
            data += delim

        data += self.DELIMITER_HEADER
        for _, file in self.FILES.items():
            data += file.name.encode('utf-8') + self.DELIMITER_HEADER
            data += base64.b64encode(file.data) + self.DELIMITER_FILE
        data += self.DELIMITER_END
        self._data = data
        return data


class Client(DelimiterMixin):
    def __init__(self, host: str, port: int, delimiter: bytes) -> None:
        super().__init__(delimiter)
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.connect((host, port))

    def Send(self, msg: Message):
        self.sock.send(msg.Generate())

    def Recv(self) -> Message:
        data = b""
        while True:
            data += self.sock.recv(1024)
            if data.endswith(self.DELIMITER_END):
                break
        msg = Message(raw=data)
        msg.Parse()
        return msg

class Server(DelimiterMixin):
    def __init__(self, host: str, port: int, delimiter: bytes) -> None:
        super().__init__(delimiter)
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.bind((host, port))

    def Listen(self):
        self.sock.listen()
        while True:
            client, addr = self.sock.accept()
            self.handle_client(client, addr)
        
    def handle_client(self, client: socket.socket, addr):
        while True:
            data: bytes = b""
            while True:
                data += client.recv(1024)
                if data.endswith(self.DELIMITER_END):
                    break
            msg = Message(raw=data, delimiter=self.DELIMITER)
            msg.Parse()
            msg = self.handle_message(msg, client, addr)
            self.send(msg, client)

    def send(self, msg: Message, client: socket.socket):
        client.send(msg.Generate())
    
    def handle_message(self, msg: Message, client: socket.socket, addr):
        return msg


if __name__ == "__main__":
    body = "BODY"*3
    msg = Message(headers={"key": ["value", "value", "value", "value1"]}, body=bytes(body.encode('utf-8')))
    msg.AddFile("test.txt", b"test")
    msg.AddFile("test2.txt", b"test2")
    msg.AddFile("test3.txt", b"test3")
    msg.Generate()
    msg.Parse()

    def startserver():
        server = Server("localhost", 1234, b"$")
        server.Listen()

    threading.Thread(target=startserver).start()
    
    client = Client("localhost", 1234, b"$")
    client.Send(msg)
    msg = client.Recv()
    print(msg._data)
    print(msg.HEADERS, msg.BODY, msg.FILES)
    [print(file.name, not not file.data) for k, file in msg.FILES.items()]
    os._exit(0)

