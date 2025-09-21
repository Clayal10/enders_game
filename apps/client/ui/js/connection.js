const setupAPI = "/lurk-client/setup/"
const joinAPI = "/lurk-client/join"

// Sends:
// - Hostname
// - Port
// Receives:
// - lurk.Game object
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
            return response.json()
        }).then(data => {
            setGamePreview(data)
        });
    }catch(e){
        console.error("Could not send config: ", e)
        return
    }
}

// Set game-title to generic lurk
// Set game-desc to the description
function setGamePreview(data){
    const gameTitle = document.getElementById("game-title");
    const gameDesc = document.getElementById("game-text");

    gameTitle.innerHTML = "Lurk Server:";

    gameContent = data.GameDesc.replace(/\n/g, '<br>');
    gameDesc.innerHTML = gameContent;
}
