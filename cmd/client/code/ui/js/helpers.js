function updateGame(data) {
    const gameDesc = document.getElementById("game-text");
    const gamePlayers = document.getElementById("game-players");
    const gameRooms = document.getElementById("game-rooms");
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
    reveal("input-label-name")
    reveal("input-text-name")
    reveal("input-label-attack")
    reveal("input-text-attack")
    reveal("input-label-defense")
    reveal("input-text-defense")
    reveal("input-label-regen")
    reveal("input-text-regen")
    reveal("input-label-join")
    reveal("input-text-join")
    reveal("input-label-description")
    reveal("input-text-description")
    reveal("input-submit-button")
    reveal("input-submit-button-auto")
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

function hideConfig(){
    hide("input-label-name")
    hide("input-text-name")
    hide("input-label-attack")
    hide("input-text-attack")
    hide("input-label-defense")
    hide("input-text-defense")
    hide("input-label-regen")
    hide("input-text-regen")
    hide("input-label-join")
    hide("input-text-join")
    hide("input-label-description")
    hide("input-text-description")
    hide("input-submit-button")
    hide("input-submit-button-auto")
}

function cleanup() {
    shouldPoll = false;
    hide("terminate-button");
    hideGameInput()
    hide("submit-button");
    hide("input-label-name")
    hide("input-text-name")
    hide("input-label-attack")
    hide("input-text-attack")
    hide("input-label-defense")
    hide("input-text-defense")
    hide("input-label-regen")
    hide("input-text-regen")
    hide("input-label-join")
    hide("input-text-join")
    hide("input-label-description")
    hide("input-text-description")
    hide("input-submit-button")
    hide("input-submit-button-auto")
    reveal("submit-button")
    document.getElementById("game-text").innerHTML = "";
    document.getElementById("game-players").innerHTML = "";
    document.getElementById("game-rooms").innerHTML = "";
}

var userCharacter = {};

function addName() {
    userCharacter.name = document.getElementById("input-text").value;
    clearText("input-text");
}
function addAttack() {
    userCharacter.attack = document.getElementById("input-text").value;
    clearText("input-text");
}
function addDefense() {
    userCharacter.defense = document.getElementById("input-text").value;
    clearText("input-text");
}
function addRegen() {
    userCharacter.regen = document.getElementById("input-text").value;
    clearText("input-text");
}
function addDescription() {
    userCharacter.description = document.getElementById("input-text").value;
    clearText("input-text");
}

function getCharacterInput() {
    userCharacter.name = document.getElementById("input-text-name").value;
    clearText("input-text-name");
    userCharacter.attack = document.getElementById("input-text-attack").value;
    clearText("input-text-attack");
    userCharacter.defense = document.getElementById("input-text-defense").value;
    clearText("input-text-defense");
    userCharacter.regen = document.getElementById("input-text-regen").value;
    clearText("input-text-regen");
    userCharacter.join = document.getElementById("input-text-join").value;
    clearText("input-text-join");
    userCharacter.description = document.getElementById("input-text-description").value;
    clearText("input-text-description");
    hide("game-input");
    return userCharacter;
}

function generateCharacter(){
    userCharacter.name = document.getElementById("input-text-name").value;
    clearText("input-text-name");
    userCharacter.attack = "nil"
    return userCharacter
}