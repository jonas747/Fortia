#fnet

This is the reference client and server package for the fnet protocol (short for fortia net), which basically is protocol buffers rpc (or very similar)

##fnet protocol
fnet is a very lightweight and simple networking protocol (on top of tcp, websocket or whatever you want to implement yourself) which uses protocol buffers to encode the payload.

Basically this is it:

     ----------------------------------
    | evtid | payload length | payload |
     ----------------------------------

 - evtid: Unsingned 32 bit integer representing what kind of message is received
 - payload length: Unsigned 32 bit integer representing the length of the payload in bytes
 - payload: The actual payload encoded in the procotol buffer format, the type of the protocol buffer is decided by the event id

As you can see the header is only 64 bits(8 bytes) long

##Example
Examples can be found in the examples folder