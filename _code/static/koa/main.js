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

    const app = new Koa();
    const router = new Router();

    router.get('/', (ctx) => {
        ctx.body = 'Index';
    });

    app.use(router.routes());
    app.listen(5000);
}
