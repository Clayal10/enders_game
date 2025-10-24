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

function sendStart(){
    try{
        const start = getCharacterInput();
        fetch(client.startAPI, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(start)
        }).then(response => {
            if(!response.ok){
                handleCharacterError();
                throw new Error("Bad Response");
            }
            shouldPoll = true;
        })
        hide("input-button");
        hide("game-input");
        revealGameInput();
    }catch(e){
        console.error("Could not send start: ", e);
        return
    }finally{
        pollUpdateEP();
    }
}

function sendChangeRoom(){
    try{
        changeRoom = {
            roomNumber: document.getElementById("game-input-change-room").value
        };
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
        loot = {
            target: document.getElementById("game-input-loot").value
        }
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
        pvp = {
            target: document.getElementById("game-input-pvp-fight").value
        }
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
        msg = {
            recipient: document.getElementById("game-input-message-recipient").value,
            text: document.getElementById("game-input-message").value
        }
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

