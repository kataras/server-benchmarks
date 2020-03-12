process.env.NODE_ENV = 'production';

const express = require('express');
const createWorker = require('throng');


createWorker(createWebServer) // to make it multi-threading/faster.

function createWebServer() {
    const app = express();

    app.get('/', function (req, res) {
        res.send('Index');
    });

    app.listen(5000);
}