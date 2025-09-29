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
    let button = document.getElementById("input-button");
    label.innerHTML = "Character Name:";
    button.innerHTML = "Next";
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
    hide("terminate-button");
    reveal("submit-button");
}

var userCharacter = {};

function addName(){
    userCharacter.name = document.getElementById("input-text").value;
    hide("input-button-name");
    reveal("input-button-attack");
}
function addAttack(){
    userCharacter.attack = document.getElementById("input-text").value;
    hide("input-button-attack");
    reveal("input-button-defense");
}
function addDefense(){
    userCharacter.defense = document.getElementById("input-text").value;
    hide("input-button-defense");
    reveal("input-button-regen");
}
function addRegen(){
    userCharacter.regen = document.getElementById("input-text").value;
    hide("input-button-regen");
    reveal("input-button-description");
}
function addDescription(){
    userCharacter.description = document.getElementById("input-text").value;
    hide("input-button-description");
    reveal("input-button");
}

function getCharacterInput(){
    return userCharacter;
}