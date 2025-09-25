const setupAPI = "/lurk-client/setup/"
const startAPI = "/lurk-client/start"

// Sends:
// - Hostname
// - Port
// Receives:
// - Client Update object:
//  - info | general info about the game
//  - players | Already string formatted player / monster list
//  - 
async function sendConfig(){
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
            updateGame(data);
        });
    }catch(e){
        console.error("Could not send config: ", e);
        return
    }
}

async function sendStart(){
    try{
        const start = {
            "start": ""
        };
        fetch(startAPI, {
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
            return response.json();
        }).then(data => {
            updateGame(data);
        });
    }catch(e){
        console.error("Could not send start: ", e);
        return
    }
}

function updateGame(data){
    const gameDesc = document.getElementById("game-text");
    const gamePlayers = document.getElementById("game-players");
    const gameRooms = document.getElementById("game-rooms");

    console.log(data.id);

    gameDesc.innerHTML += data.info.replace(/\n/g, '<br>');
    gamePlayers.innerHTML += data.players.replace(/\n/g, '<br>');
    gameRooms.innerHTML += data.rooms.replace(/\n/g, '<br>');
}
