process.env.NODE_ENV = 'production';

const express = require('express');
const createWorker = require('throng');


createWorker(createWebServer) // to make it multi-threading/faster.

function createWebServer() {
    const app = express();
    app.use(express.json()); // express v4.16+.

    app.post('/:id', function (req, res) {
        const id = parseInt(req.params.id);
        const input = req.body;
        res.json({
            id: id,
            name: input.email,
        });
    });

    app.listen(5000);
}