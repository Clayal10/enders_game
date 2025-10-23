function updateGame(data) {
    const gameDesc = document.getElementById("game-text");
    const gamePlayers = document.getElementById("game-players");
    const gameRooms = document.getElementById("game-rooms");
    console.log(data)
    gameDesc.innerHTML = data.info.replace(/\n/g, '<br>');
    gamePlayers.innerHTML = data.players.replace(/\n/g, '<br>');
    gameRooms.innerHTML = data.rooms.replace(/\n/g, '<br>');
    gameRooms.innerHTML += data.connections.replace(/\n/g, `<br>`);

    gameDesc.scrollTop = gameDesc.scrollHeight;
}

function setupDisplay() {
    hide("submit-button");
    reveal("terminate-button");
    reveal("game-input")
    let label = document.getElementById("input-label");
    label.innerHTML = "Character Name:";
    reveal("input-button-name");
    reveal("input-label");
    reveal("input-text");
}

const errorCharacter = '<span style="color: red;">Error</span>: Invalid character settings, try again.'
function handleCharacterError() {
    document.getElementById("game-text").innerHTML += errorCharacter
    setupDisplay()
}

function hide(id) {
    document.getElementById(id).classList.add("hidden");
}

function reveal(id) {
    document.getElementById(id).classList.remove("hidden");
}

function hideGameInput() {
    let el = document.getElementById("game-input-main")
    el.style.display = 'none';
}

function revealGameInput() {
    let el = document.getElementById("game-input-main")
    el.style.display = 'flex';
}

function clearText(id) {
    document.getElementById(id).value = "";
}

function cleanup() {
    shouldPoll = false;
    hide("input-text");
    hide("input-label");
    hide("input-text");
    hide("input-button");
    hide("input-button-attack");
    hide("input-button-defense");
    hide("input-button-regen");
    hide("input-button-description");
    hide("input-button-name");
    hide("terminate-button");
    hideGameInput()
    hide("game-input");
    reveal("submit-button");
    document.getElementById("game-text").innerHTML = "";
    document.getElementById("game-players").innerHTML = "";
    document.getElementById("game-rooms").innerHTML = "";
}

var userCharacter = {};

//This should be an overlay type with all elements shown at once.
function addName() {
    userCharacter.name = document.getElementById("input-text").value;
    hide("input-button-name");
    clearText("input-text");
    document.getElementById("input-label").innerHTML = "Enter Attack: "
    reveal("input-button-attack");
}
function addAttack() {
    userCharacter.attack = document.getElementById("input-text").value;
    hide("input-button-attack");
    clearText("input-text");
    document.getElementById("input-label").innerHTML = "Enter Defense: "
    reveal("input-button-defense");
}
function addDefense() {
    userCharacter.defense = document.getElementById("input-text").value;
    hide("input-button-defense");
    clearText("input-text");
    document.getElementById("input-label").innerHTML = "Enter Regen: "
    reveal("input-button-regen");
}
function addRegen() {
    userCharacter.regen = document.getElementById("input-text").value;
    hide("input-button-regen");
    clearText("input-text");
    document.getElementById("input-label").innerHTML = "Enter Description"
    reveal("input-button-description");
}
function addDescription() {
    userCharacter.description = document.getElementById("input-text").value;
    clearText("input-text");
    hide("input-button-description");
    hide("input-label");
    hide("input-text");
    reveal("input-button");
}

function getCharacterInput() {
    return userCharacter;
}