const setupAPI = "/lurk-client/setup/"
// These require IDs
const startAPI = "/lurk-client/start/"
const updateAPI = "/lurk-client/update/"
const terminateAPI = "/lurk-client/terminate/"
const changeRoomAPI = "/lurk-client/change-room/"
const fightAPI = "/lurk-client/fight/"
const lootAPI = "/lurk-client/loot/"
const pvpFightAPI = "/lurk-client/pvp/"
const messageAPI = "/lurk-client/message/"

class Client{
    constructor(id){
        this.id = id;
        this.startAPI = startAPI+id+"/";
        this.updateAPI = updateAPI+id+"/";
        this.terminateAPI = terminateAPI+id+"/";
        this.changeRoomAPI = changeRoomAPI+id+"/";
        this.fightAPI = fightAPI+id+"/";
        this.lootAPI = lootAPI+id+"/";
        this.pvpFightAPI = pvpFightAPI+id+"/";
        this.messageAPI = messageAPI+id+"/";
    };
};


var client;

window.addEventListener('beforeunload', (event) => {
  if (client.id === 0){
    return
  }
  navigator.sendBeacon(client.terminateAPI);
  shouldPoll = false;
});

// Sends:
// - Hostname
// - Port
// Receives:
// - Client Update object:
//  - info | general info about the game
//  - players | Already string formatted player / monster list
function sendConfig(){
    try{
        let hostname = document.getElementById("input-hostname");
        let port = document.getElementById("input-port");

        const cfg = {
            "Hostname": hostname.value,
            "Port": port.value
        };

        fetch(setupAPI, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(cfg)
        }).then(response => {
            if(!response.ok){
                throw new Error("Bad Response");
            }
            return response.json();
        }).then(data => {
            if(data.id === ""){
                throw new Error("No valid ID");
            }
            console.log("New client ID: ", data.id);
            client = new Client(data.id)
            updateGame(data);
        });
        setupDisplay(); // For character input.
    }catch(e){
        console.error("Could not send config: ", e);
        cleanup();
    }
}

function sendTerminate(){
    try{
        cleanup();
        fetch(client.terminateAPI, {
            method: 'POST',
            headers:{
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify({})
        })
    }catch(e){
        console.error("Could not send terminate: ", e)
    }
}

function autoSendStart(){
    const character = generateCharacter();
    sendCharacter(character);
}

function sendStart(){
    const character = getCharacterInput();
    sendCharacter(character);
}

function sendCharacter(character){
    try{
        fetch(client.startAPI, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(character)
        }).then(response => {
            if(!response.ok){
                handleCharacterError();
                throw new Error("Bad Response");
            }
            shouldPoll = true;
            hide("input-submit-button");
            hideConfig();
            revealGameInput();
            pollUpdateEP();
        })
    }catch(e){
        console.error("Could not send start: ", e);
        return
    }
}

function sendChangeRoom(){
    try{
        let changeRoomElement = document.getElementById("game-input-change-room")
        let changeRoom = {
            roomNumber: changeRoomElement.value
        };
        changeRoomElement.value = "";
        fetch(client.changeRoomAPI, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(changeRoom)
        }).then(response =>{
            if(!response.ok){
                throw new Error("Bad Response");
            }
        })
    }catch(e){
        console.error("Could not send change room: ", e);
    }
}

function sendFight(){
    try{
        fetch(client.fightAPI, {
            method: 'POST',
            headers:{
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify({})
        }).then(response => {
            if(!response.ok){
                throw new Error("Bad Response");
            }
        })
    }catch(e){
        console.error("Could not initiate fight: ", e);
    }
}

function sendLoot(){
    try{
        let lootElement = document.getElementById("game-input-loot")
        loot = {
            target: lootElement.value
        }
        lootElement.value = "";

        fetch(client.lootAPI, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(loot)
        }).then(response => {
            if(!response.ok){
                throw new Error("Bad Response");
            }
        })
    }catch(e){
        console.error("Could not loot target: ", e);
    }
}

function sendPVP(){
    try{
        let pvpElement = document.getElementById("game-input-pvp-fight");
        let pvp = {
            target: pvpElement.value
        }
        pvpElement.value = "";

        fetch(client.pvpFightAPI, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(pvp)
        }).then(response => {
            if(!response.ok){
                throw new Error("Bad Response");
            }
        })
    }catch(e){
        console.error("Could not pvp fight: ", e)
    }
}

function sendMessage(){
    try{
        let msgElement = document.getElementById("game-input-message-recipient")
        let msgMessage = document.getElementById("game-input-message")
        let msg = {
            recipient: msgElement.value,
            text: msgMessage.value
        }
        msgElement.value = "";
        msgMessage.value = "";

        fetch(client.messageAPI, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(msg)
        }).then(response => {
            if(!response.ok){
                throw new Error("Bad Response");
            }
        })
    }catch(e){
        console.error("Could not send message: ", e)
    }
}

var shouldPoll = true;
async function pollUpdateEP(){
    try{
        let response = await fetch(client.updateAPI)
        if(response.status === 200){
            data = await response.json();
            updateGame(data);
        }
    }catch(e){
        console.error(e);
    }finally{
        if(shouldPoll){
            await pollUpdateEP();
        }
    }
}

