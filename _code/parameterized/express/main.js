process.env.NODE_ENV = 'production';

const express = require('express');
const createWorker = require('throng');


createWorker(createWebServer) // multi-thread.

function createWebServer() {
    const app = express();

    app.get('/hello/:name', function (req, res) {
        res.send('Hello ' + req.params.name);
    });

    app.listen(5000);
}