process.env.NODE_ENV = 'production';

const express = require('express');
const createWorker = require('throng');


createWorker(createWebServer) // multi-thread.

function createWebServer() {
    const app = express();
    app.use(express.json({limit: '2mb'})); // express v4.16+.
    app.use(express.urlencoded({limit: '2mb'}));

    app.post('/:id', function (req, res) {
        const id = parseInt(req.params.id);
        const inputs = req.body;
        res.json({
            id: id,
            count: inputs.length,
            first_id: inputs[0].id
        });
    });

    app.listen(5000);
}