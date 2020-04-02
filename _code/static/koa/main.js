process.env.NODE_ENV = 'production';

const Koa = require('koa');
const Router = require('@koa/router');
const createWorker = require('throng');


createWorker(createWebServer) // multi-thread.

function createWebServer() {
    const app = new Koa();
    const router = new Router();

    router.get('/', function (ctx) {
        ctx.body = 'Index';
    });

    app.use(router.routes());
    app.listen(5000);
}