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
    const app = require('fastify')({ logger: false });

    app.get('/', (req, reply) => {
        reply.send('Index');
    });

    app.listen({ port: 5000 });
}
