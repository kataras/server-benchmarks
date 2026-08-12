process.env.NODE_ENV = 'production';

const cluster = require('node:cluster');
const { availableParallelism } = require('node:os');

if (cluster.isPrimary) {
    // One worker per logical CPU: a Node process is single-threaded.
    for (let i = 0; i < availableParallelism(); i++) {
        cluster.fork();
    }
} else {
    createWebServer();
}

function createWebServer() {
    const express = require('express');
    const app = express();
    app.use(express.json({ limit: '2mb' }));

    app.post('/:id', (req, res) => {
        const id = Number.parseInt(req.params.id, 10);
        if (Number.isNaN(id)) {
            // * Express does not support parameter type-based routing.
            return res.sendStatus(404);
        }

        const inputs = req.body;
        res.json({
            id: id,
            count: inputs.length,
            first_id: inputs[0].id
        });
    });

    app.listen(5000);
}
