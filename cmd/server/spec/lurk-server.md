# Lurk Server

![Server Layout](./lurk-server.drawio.svg)

## Introduction

This server complies with the Lurk protocol version 2.3. This guide outlines the functionality of a general Lurk server. For specific gameplay information, visit [here](./enders-game.md).

## Configuration

For custom configuration options, a `Config.json` file may be present in the directory of the server's binary. This will overwrite the default settings as shown here:

```json
{
    "ServerPort": 5069
}
```

> NOTE: More configuration options may be added in the future.

## Client Interface

For a description on all message types denoted in brackets, (e.g. [CHARACTER]), visit the [Lurk Library](../../../pkg/lurk/README.md)

### Initial Connection

To connect to this Lurk server, create a TCP connection on the specified server port. The client will then be sent a [GAME] message. This will be a description of the game.

From there, the client must give a [CHARACTER] message describing the user's character. The server will then either send the same [CHARACTER] message back to the client with the _ready_ flag set or an [ERROR]. 

To start the game, the client must send a [START] message to the server. The client will be able to set their own _name_, _description_, _attack_, _defense_, _regen_ and _flags_ on creation.

The character will then be placed in a room and [gameplay](#gameplay) will begin.

### Gameplay

#### Movement

When entering a room or connecting to the server, the client will receive a [ROOM] message along with a [CHARACTER] message for each player and monster in the room. A [CHARACTER] message will also be sent to the client when another player or entity leaves the room. When initially moving to a room (or at any point), multiple [CONNECTION] messages may be sent to display each room connected to the player's current location.

To change room, the client should send a [CHANGEROOM] message. The server will then send a [ROOM], any number of [CHARACTER], and any number of [CONNECTION] messages.

#### Player / Entity Interaction
