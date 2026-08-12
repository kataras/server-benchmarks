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
    const Koa = require('koa');
    const Router = require('@koa/router');
    const { bodyParser } = require('@koa/bodyparser');

    const app = new Koa();
    const router = new Router();

    router.post('/:id', (ctx) => {
        const id = Number.parseInt(ctx.params.id, 10);
        if (Number.isNaN(id)) {
            // * Koa does not support parameter type-based routing.
            ctx.status = 404;
            return;
        }

        const inputs = ctx.request.body;
        ctx.body = {
            id: id,
            count: inputs.length,
            first_id: inputs[0].id
        };
    });

    app.use(bodyParser({ jsonLimit: '2mb' }));
    app.use(router.routes());

    app.listen(5000);
}
