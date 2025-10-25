# LURK Protocol

The LURK protocol is used in text-based MMORPG-style games. Messages are sent little-endian.

## Message Types

### MESSAGE

Sent by the client to message other players. Can also be used to send presentable information to the client from the server.

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 1|
|1|2|Message Length (n)|
|3|32|Recipient Name|
|35|30|Sender Name|
|65|2|Null terminator and flag for narration|
|67+|n|Message content|

### CHANGEROOM

Sent by the client only, to change rooms. If the server changes the room a client is in, it should send an updated room, character, and connection message(s) to explain the new location. If not, for example because the client is not ready to start or specified an inappropriate choice, and error should be sent.

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 2|
|1|2|Room number to change to|

### FIGHT

Initiate a fight against monsters. This will start a fight in the current room against the monsters which are presently in the room. Players with the join battle flag set, who are in the same room, will automatically join in the fight. The server will allocate damage and rewards after the battle, and inform clients appropriately. Clients should expect a slew of messages after starting a fight, especially in a crowded room. This message is sent by the client. If a fight should ensue in the room the player is in, the server should notify the client, but not by use of this message. Instead, the players not initiating the fight should receive an updated CHARACTER message for each entity in the room. If the server wishes to send additional narrative text, this can be sent as a MESSAGE. Note that this is not the only way a fight against monsters can be initiated. The server can initiate a fight at any time.

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 3|

### PVPFIGHT

Initiate a fight against another player. The server will determine the results of the fight, and allocate damage and rewards appropriately. The server may include players with join battle in the fight, on either side. Monsters may or may not be involved in the fight as well. This message is sent by the client. If the server does not support PVP, it should send error 8 to the client.

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 4|
|1|32|Name of target player|

### LOOT

Loot gold from a dead player or monster. The server may automatically gift gold from dead monsters to the players who have killed them, or wait for a LOOT message. The server is responsible for communicating the results of the LOOT to the player, by sending an updated CHARACTER message. This message is sent by the client. 

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 5|
|1|32|Name of target player|

### START

Start playing the game. A client will send a CHARACTER message to the server to explain character stats, which the server may either accept or deny (by use of an ERROR message). If the stats are accepted, the server will not enter the player into the game world until it has received START. This is sent by the client. Generally, the server will reply with a ROOM, a CHARACTER message showing the updated room, and a CHARACTER message for each player in the initial room of the game.

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 6|

### ERROR

Notify the client of an error. This is used to indicate stat violations, inappropriate room connections, attempts to loot nonexistent or living players, attempts to attack players or monsters in different rooms, etc. 

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 7|
|1|1|Error code|
|2|2|Error message length (n)|
|4+|n|Error message|

#### Error Codes

|Code|Description|
|---|---|
|0|Other|
|1|Bad Room.|
|2|Player already exists|
|3|Bad monster|
|4|Stat error|
|5|Not Ready|
|6|No target|
|7|No fight|
|8|No pvp combat on the server|

### ACCEPT

Sent by the server to acknowledge a non-error-causing action which has no other direct result. This is not needed for actions which cause other results, such as changing rooms or beginning a fight. It should be sent in response to clients sending messages, setting character stats, etc. 

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 8|
|1|1|Type of action|


### ROOM

Sent by the server to describe the room that the player is in. This should be an expected response to CHANGEROOM or START. Can be re-sent at any time, for example if the player is teleported or falls through a floor. Outgoing connections will be specified with a series of CONNECTION messages. Monsters and players in the room should be listed using a series of CHARACTER messages.

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 9|
|1|2|Room number. Same number used for CHANGEROOM
|3|32|Room name|
|35|2|Room description length (n)|
|37+|n|Room description|


### CHARACTER

Sent by both the client and the server. The server will send this message to show the client changes to their player's status, such as in health or gold. The server will also use this message to show other players or monsters in the room the player is in or elsewhere. The client should expect to receive character messages at any time, which may be updates to the player or others. If the player is in a room with another player, and the other player leaves, a CHARACTER message should be sent to indicate this. In many cases, the appropriate room for the outgoing player is the room they have gone to. If the player goes to an unknown room, the room number may be set to a room that the player will not encounter (does not have to be part of the map). This could be accompanied by a narrative message (for example, "Glorfindel vanishes into a puff of smoke"), but this is not required.

The client will use this message to set the name, description, attack, defense, regen, and flags when the character is created. It can also be used to reprise an abandoned or deceased character.

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 10|
|1|32|Name of player|
|33|1|[Flags](#flag-bit-word). 8 bit word which uses the 5 most significant bits (left most)|
|34|2|Attack|
|36|2|Defense|
|38|2|Regen|
|40|2|Health (Signed int)|
|42|2|Gold|
|44|2|Current room number|
|46|2|Description length (n)|
|48+|n|Player description|

#### Flag bit word

|7|6|5|4|3|2|1|0
|---|---|---|---|---|---|---|---|
|Alive|Join Battle|Monster|Started|Ready|RESERVED|RESERVED|RESERVED|

### GAME

Used by the server to describe the game. The initial points is a combination of health, defense, and regen, and cannot be exceeded by the client when defining a new character. The stat limit is a hard limit for the combination for any player on the server regardless of experience. If unused, it should be set to 65535, the limit of the unsigned 16-bit integer. This message will be sent upon connecting to the server, and not re-sent. 

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 11|
|1|2|Initial points|
|3|2|Stat limit|
|5|2|Description length (n)|
|7+|n|Game Description|

### LEAVE

Used by the client to leave the game. This is a graceful way to disconnect. The server never terminates, so it doesn't send LEAVE. 

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 12|

### CONNECTION

Used by the server to describe rooms connected to the room the player is in. The client should expect a series of these when changing rooms, but they may be sent at any time. For example, after a fight, a secret staircase may extend out of the ceiling enabling another connection. Note that the room description may be an abbreviated version of the description sent when a room is actually entered. The server may also provide a different room description depending on which room the player is in. So a description on the connection could read "A strange whirr is heard through the solid oak door", and the description attached to the message once the player has entered could read "Servers line the walls, softly lighting the room in a cacophony of red, green, blue, and yellow flashes".

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 13|
|1|2|Room Number, same room number used for CHANGEROOM|
|3|32|Room name|
|35|2|Room description length (n)|
|37+|n|Room description|

### VERSION

 Sent by the server upon initial connection along with GAME. If no VERSION is received, the server can be assumed to support only LURK 2.0 or 2.1. 

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 14|
|1|1|LURK major revision|
|2|1|LURK minor revision|
|3|2|Size of the list of extensions (n)|
|5+|n|List of extensions|

#### List of Extensions Format

|Offset|Length (bytes)|Description|
|---|---|---|
|5|2|Length of the first extension (n)|
|7+|n|First extension|
...
