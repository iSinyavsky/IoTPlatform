const mosca = require('mosca')
const express = require("express");
const app = express();
app.get("/health", (req, res) => {
    res.send({ success: true, message: "It is working" });
});
app.get("/", (req, res) => {
    // consider that this route crashes the entire application
});
const PORT = 9999;
app.listen(PORT, () => {
});

const settings = {
    port: 1883,
    http: {
        port: 3000,
        bundle: true,
        static: './'
    }
};
const server = new mosca.Server(settings, function (){

});
server.on('ready', setup);

server.on('clientConnected', function(client) {
    console.log('client connected', client.id);
});

// fired when a message is received
server.on('published', function(packet, client) {
    console.log('Published', packet.payload);
});

// fired when the mqtt broker is ready
function setup() {
    console.log('Mosca embedded MQTT broker running now')
}
