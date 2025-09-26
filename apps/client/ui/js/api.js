const setupAPI = "/lurk-client/setup/"
// These require IDs
const startAPI = "/lurk-client/start"
const updateAPI = "/lurk-client/update"

class Client{
    constructor(id){
        this.id = id;
        this.startAPI = startAPI+"/"+id+"/"
        this.updateAPI = updateAPI+"/"+id+"/"
    };
};

var client;

// Sends:
// - Hostname
// - Port
// Receives:
// - Client Update object:
//  - info | general info about the game
//  - players | Already string formatted player / monster list
//  - 
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
    }catch(e){
        console.error("Could not send config: ", e);
        return
    }
}

function sendStart(){
    try{
        const start = {
            "start": ""
        };
        fetch(client.startAPI, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(start)
        }).then(response => {
            if(!response.ok){
                throw new Error("Bad Response");
            }
        })
    }catch(e){
        console.error("Could not send start: ", e);
        return
    }finally{
        pollUpdateEP();
    }
}

async function pollUpdateEP(){
    try{
        fetch(client.updateAPI).then(response => {
            if(!response.ok){
                throw new Error("Bad response in poll");
            }
            return response.json()
        }).then(data =>{
            updateGame(data);
        });
    }catch(e){
        console.error("Could not poll endpoint: ", e);
    }finally{
        setTimeout(pollUpdateEP, 2000);
    }
}

function updateGame(data){
    const gameDesc = document.getElementById("game-text");
    const gamePlayers = document.getElementById("game-players");
    const gameRooms = document.getElementById("game-rooms");

    gameDesc.innerHTML += data.info.replace(/\n/g, '<br>');
    gamePlayers.innerHTML += data.players.replace(/\n/g, '<br>');
    gameRooms.innerHTML += data.rooms.replace(/\n/g, '<br>');
}
