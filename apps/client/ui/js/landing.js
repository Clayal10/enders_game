const setupAPI = "/lurk-client/setup/"

async function sendConfig(){
    try{
        let hostname = document.getElementById("input-hostname");
        let port = document.getElementById("input-port");

        const cfg = {
            "Hostname": hostname.value,
            "Port": port.value
        };

        const response = await fetch(setupAPI, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(cfg)
        });

        if(!response.ok){
            throw new Error("Bad Response");
        }

    }catch(e){
        console.error("Could not send config: ", e)
        return
    }
}
