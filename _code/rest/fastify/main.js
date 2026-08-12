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
    const app = require('fastify')({ logger: false, bodyLimit: 2 * 1024 * 1024 });

    // The response schema enables fastify's fast-json-stringify path,
    // which is its idiomatic production setup.
    app.post('/:id', {
        schema: {
            response: {
                200: {
                    type: 'object',
                    properties: {
                        id: { type: 'integer' },
                        count: { type: 'integer' },
                        first_id: { type: 'string' }
                    }
                }
            }
        }
    }, (req, reply) => {
        const id = Number.parseInt(req.params.id, 10);
        if (Number.isNaN(id)) {
            // * Fastify does not support parameter type-based routing.
            reply.code(404).send();
            return;
        }

        const inputs = req.body;
        reply.send({
            id: id,
            count: inputs.length,
            first_id: inputs[0].id
        });
    });

    app.listen({ port: 5000 });
}
