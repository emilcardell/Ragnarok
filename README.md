Ragnarok
========

Some kind of message based data base built in Go. Just because I want to learn Go.
This i the fist basic implementation just for handling file operations.

## Build
Run build.bat

## Run
main.exe

## Commands

Server will run on port 9991

### Write to channel
GET localhost:9991/wc/ThisIsMyChannel

### Read first message
GET localhost:9991/rc/ThisIsMyChannel

### Read specific message
GET localhost:9991/rc/ThisIsMyChannel/12345

12345 = offset

