process.env.NODE_ENV = 'production';

const Koa = require('koa');
const Router = require('@koa/router');
const bodyParser = require('koa-bodyparser');

const createWorker = require('throng');


createWorker(createWebServer) // multi-thread.

function createWebServer() {
    const app = new Koa();
    const router = new Router();
    router.post('/:id', function (ctx) {
        const id = parseInt(ctx.params.id);
        const inputs = ctx.request.body;
        ctx.body = {
            id: id,
            count: inputs.length,
        };
    });

    app.use(bodyParser({limit: '2mb'}));
    app.use(router.routes());

    app.listen(5000);
}