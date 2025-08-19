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

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 2|
|1|2|Room number to change to|

### FIGHT

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 3|

### PVPFIGHT

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 4|
|1|32|Name of target player|

### LOOT

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 5|
|1|32|Name of target player|

### START

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 6|

### ERROR

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

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 8|
|1|1|Type of action|


### ROOM

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 9|
|1|2|Room number. Same number used for CHANGEROOM
|3|32|Room name|
|35|2|Room description length (n)|
|37+|n|Room description|


### CHARACTER

Can be used by the client to describe a new 

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 10|
|1|32|Name of player|
|33|1|Flags. 8 bit word which uses the 5 most significant bits (right most in little endian)|
|34|2|Attack|
|36|2|Defense|
|38|2|Regen|
|40|2|Health (Signed int)|
|42|2|Gold|
|44|2|Current room number|
|46|2|Description length (n)|
|48+|n|Player description|

#### Flag bit word

|0|1|2|3|4|5|6|7|
|---|---|---|---|---|---|---|---|
|RESERVED|RESERVED|RESERVED|Ready|Started|Monster|Join Battle|Alive|

### GAME

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 11|
|1|2|Initial points|
|3|2|Stat limit|
|5|2|Description length (n)|
|7+|n|Game Description|

### LEAVE

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 12|

### CONNECTION

|Offset|Length (bytes)|Description|
|---|---|---|
|0|1|Type 13|
|1|2|Room Number, same room number used for CHANGEROOM|
|3|32|Room name|
|35|2|Room description length (n)|
|37+|n|Room description|

### VERSION

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
