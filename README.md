Yo'ed server
===========

Yo'ed server which handles incoming calls from Yo API.

#Installation
You need the go package on your machine to get the source

`go get github.com/yoed/yoed-server`


#Configuration

Create a `config.json` file aside the executable program.

##listen

The `ip:port` to listen to.

#Protocol

The Yo'ed client/server protocol is a simple protocol based on HTTP.

Clients can subscribe to notifications for one or many handles by calling the `yo` endpoint with the `handles` (separated by commas) they want to listen to and the `callback_url` to call when a Yo is received.

The `yoed` endpoint is the server entry point. For each Yo account you own, you can setup a callback URL like this: http://my-server-ip:port/yoed/:handle

When a Yo is received, the Yo API will call you on this URL (HTTP/GET) with a `username` parameter corresponding to the user who have Yo'ed you. This value will be passed to all the clients which have subscribed.

The server can handle and dispatch actions to clients for many handles, just set one URL per handle in the [Yo API dashboard](http://developer.justyo.co).

#Todo
* Add a security layer for preventing clients to subscribe to unauthorized handles
* Use (libchan by docker)[https://github.com/docker/libchan] instead of pure HTTP protocol to be able to have in memory clients, over Unix socket, etc.
* Add some unit tests
