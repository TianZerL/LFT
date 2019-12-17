# LFT - *LAN files transferor*
LFT is a simple cli tool to transfer your files in LAN, also works in WAN, based on Go.  
It can be called as *LAN files transferor* :joy:.

# Feature
- [x] Send a file
- [x] Send directory
- [ ] Scan server

# Example
Start a server:  
![How to start a server](example/server.gif)
Send a file:
![How to send a file](example/sendFile.gif)
Send a directory
![How to send a directory](example/sendDir.gif)

# Usage
Start a server  
```
LFT -w  
```
Send a file or directory  
```
LFT -d [source path] -ip [server ip]  
```
More argument 
```
    -?    Display help information
    -d string
            Source or destination (default "./receive/")
    -h    Display help information
    -ip string
            Server IP address (default "0.0.0.0")
    -port string
            Server Port (default "6981")
    -scan
            Scan Lan to find servers(TODO)
    -w    Start a server
```

# Install
```
go get -v github.com/TianZerL/LFT
```

# Author
TianZerL

# License
LGPL-3.0