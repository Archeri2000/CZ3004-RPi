"""
this file serves as an entry point with which we will invoke the image recognition submodule
    the main functionalities of this file are 2-fold: 
    1. capturing the actual image and resizing
    2. invoke submodule with image captured
"""

import fileinput
import numpy as np 
import json
from typing import List, Tuple 
import socket
import sys
# from picamera import PiCamera # INSTALL PI CAMERA
# import picamera.array
import argparse
import constants 
# import mdp_imagerec # TODO: set wanting's directory to be in same file hierachy 

# maybe can take return dict?
def takePhoto(metadata):
    """takes a photo and returns a tuple of photo, meta data 
@params: metadata of the image 
@return: tuple of (numpy array, meta-data)
idk man, meta data might be useful - orientation, pos of robot
"""
    with picamera.PiCamera() as camera:
        with picamera.array.PiRGBArray(camera) as output:
            camera.resolution = (constants.CAMERA_WIDTH, constants.CAMERA_HEIGHT)
            camera.capture(output, 'rgb')
            return (output.array, metadata)

if __name__ == "__main__":
    # parser = argparse.ArgumentParser(description='RPi camera interface')
    # parser.add_argument('o', type = int, help = 'orientation of the robot')
    # parser.add_argument('c', type = int, nargs = 2, help = "co-ordinate of the robot")
    # args = parser.parse_args()
    # metadata = [args.o] + args.c
    # pass image to downstream model 
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        # Connect to server and send data
        sock.connect((constants.HOST, constants.PORT))
        for line in fileinput.input():
        # keep alive and persistently read 
        # format: o, x, y
            metadata = line.split()
            img = takePhoto(metadata)
            sock.sendall(json.dumps({"data":img[0], "metadata":img[1]}).encode(encoding="UTF-8"))
            received = json.loads(str(sock.recv(1024), "utf-8")) 
            print(received) # idk if this will actually work but we are using stdout as handler b/w go and py 
