function updateGame(data){
    const gameDesc = document.getElementById("game-text");
    const gamePlayers = document.getElementById("game-players");
    const gameRooms = document.getElementById("game-rooms");

    gameDesc.innerHTML = data.info.replace(/\n/g, '<br>');
    gamePlayers.innerHTML = data.players.replace(/\n/g, '<br>');
    gameRooms.innerHTML = data.rooms.replace(/\n/g, '<br>');

    gameDesc.scrollTop = gameDesc.scrollHeight;
}

function setupDisplay(){
    hide("submit-button");
    reveal("terminate-button");
    let label = document.getElementById("input-label");
    label.innerHTML = "Character Name:";
    reveal("input-button-name");
    reveal("input-label");
    reveal("input-text");
}

function hide(id){
    document.getElementById(id).classList.add("hidden");
}

function reveal(id){
    document.getElementById(id).classList.remove("hidden");
}

function cleanup(){
    shouldPoll = false;
    hide("input-text");
    hide("input-label");
    hide("input-text");
    hide("input-button");
    hide("terminate-button");
    reveal("submit-button");
    document.getElementById("game-text").innerHTML = "";
    document.getElementById("game-players").innerHTML = "";
    document.getElementById("game-rooms").innerHTML = "";
}

var userCharacter = {};

function addName(){
    userCharacter.name = document.getElementById("input-text").value;
    hide("input-button-name");
    document.getElementById("input-label").innerHTML = "Enter Attack: "
    reveal("input-button-attack");
}
function addAttack(){
    userCharacter.attack = document.getElementById("input-text").value;
    hide("input-button-attack");
    document.getElementById("input-label").innerHTML = "Enter Defense: "
    reveal("input-button-defense");
}
function addDefense(){
    userCharacter.defense = document.getElementById("input-text").value;
    hide("input-button-defense");
    document.getElementById("input-label").innerHTML = "Enter Regen: "
    reveal("input-button-regen");
}
function addRegen(){
    userCharacter.regen = document.getElementById("input-text").value;
    hide("input-button-regen");
    document.getElementById("input-label").innerHTML = "Enter Description"
    reveal("input-button-description");
}
function addDescription(){
    userCharacter.description = document.getElementById("input-text").value;
    hide("input-button-description");
    hide("input-label");
    hide("input-text");
    reveal("input-button");
}

function getCharacterInput(){
    return userCharacter;
}