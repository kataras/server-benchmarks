process.env.NODE_ENV = 'production';

const express = require('express');
const createWorker = require('throng');
const session = require('express-session');
const MemoryStore = require('session-memory-store')(session);
const uuidv4 = require('uuid').v4;

createWorker(createWebServer) // multi-thread.

function createWebServer() {
    const app = express();
    app.use(session({
        store: new MemoryStore({
            expires: 1 * 60,
            checkperiod: 1 * 60
        }),
        secret: 'session',
        resave: true,
        saveUninitialized: false,
        cookie: {
            secure: true,
            maxAge: 60000
        }
    }));

    app.get('/sessions', function (req, res) {
        req.session.id = uuidv4();
        req.session.name = 'John Doe';
        req.session.save();

        var name = req.session.name;
        res.send(name);
    });

    app.listen(5000);
}