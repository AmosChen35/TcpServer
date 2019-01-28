# TcpServer
## Simple Usage
* This program is include server mode and client mode and please use the `-h` to see more information about additionally flags.

## Execute Program With Server Mode
* The local console started and the default ip-address and port is "127.0.0.1:8080".
  ```sh
  $ server console
  ```
  <img src="https://github.com/AmosChen35/TcpServer/blob/master/screenshot/execute_server_mode.gif" height="374" width="641">


## Execute Program With Client Mode
* The remote console will start and You can use the `attach` flag connect to the server.
  ```sh
  $ server attach attach 127.0.0.1:8080
  ```
  <img src="https://github.com/AmosChen35/TcpServer/blob/master/screenshot/execute_client_mode.gif" height="374" width="641">

## Call Contorl Through Console
* Execute the local sample contorl from the client side and that will return a sample string "HelloBridge".
  ```
  > personal.Hello()
  ```
  <img src="https://github.com/AmosChen35/TcpServer/blob/master/screenshot/hello.gif" height="374" width="641">
  
* Execute the remote sample contorl from the client side and that will return a sample json.
  ```
  > personal.HelloPRC()
  ```
  <img src="https://github.com/AmosChen35/TcpServer/blob/master/screenshot/hellorpc.gif" height="374" width="641">

* Execute the connections control from the client side, that will return how many connections was connected from server in current state.
  ```
  > admin.connections()
  ```
  <img src="https://github.com/AmosChen35/TcpServer/blob/master/screenshot/connections.gif" height="374" width="641">
