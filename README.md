# Ender's Game

<div style="display: flex; flex-direction: row;">
<img src="./space_gopher.png" width="90" height="120"/>
<pre>
 ____  __ _  ____  ____  ____  _ ____     ___   __   _  _  ____ 
(  __)(  ( \(    \(  __)(  _ \(// ___)   / __) / _\ ( \/ )(  __)
 ) _) /    / ) D ( ) _)  )   /  \___ \  ( (_ \/    \/ \/ \ ) _) 
(____)\_)__)(____/(____)(__\_)  (____/   \___/\_/\_/\_)(_/(____)
</pre>
</div>

## Play Now!

Use this [LURK client](https://isoptera.lcsc.edu:5068) and connect to port **5069** with the default hostname to play in your browser.

## Description

[Enders's Game](https://en.wikipedia.org/wiki/Ender%27s_Game) is the first of many books in a series detailing the journeys of a young man as he both saves the universe and becomes the greatest monster in all of history. Explore battle school, command school and other worlds in this MMO text-based dungeon crawler game.

## Dev Info

### Server

The _Ender's Game_ LURK server is built entirely in Go utilizing the [LURK protocol](https://isoptera.lcsc.edu/~seth/cs435/lurk_2.3.html) for client / server communication.

### Client

The client is built with a Go backend and vanilla Javascript font-end with a REST API for communication. All of the text processing is done in the backend to limit Javascript processing.

For decently real time updates, one of the endpoints is meant for long-polling. One goroutine is constantly checking if anything can be dequeued from a queue which gets populated by a goroutine reading from the server socket. Upon dequeuing a message from the server, a response is written to client.

### Bugs

There will likely be lots of bugs in the server and or client due to the protocol not having very strict rules. The client is built with the _Ender's Game_ server in mind.

### Other LURK Info

This [Wireshark Dissector](https://github.com/Clayal10/lurk_dissector) was very helpful while developing both the server and client. 