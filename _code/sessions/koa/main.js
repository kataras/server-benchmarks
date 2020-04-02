process.env.NODE_ENV = 'production';

const Koa = require('koa');
const Router = require('@koa/router');
const createWorker = require('throng');
const session = require('koa-session');
const uuidv4 = require('uuid').v4;

createWorker(createWebServer) // multi-thread.

function createWebServer() {
    const app = new Koa();
    app.keys = ['some secret hurr'];
    app.use(session({
        key: 'session',
    }, app));
    const router = new Router();

    router.get('/sessions', function (ctx) {
        ctx.session.id = uuidv4();
        ctx.session.name = 'John Doe';
        var name = ctx.session.name;
        ctx.body = name;
    });

    app.use(router.routes());
    app.listen(5000);
}