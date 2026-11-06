// The header must distinguish "turn over" from "work finished": a background Job
// keeps running after its turn ends, and idle alone hid that for half an hour.
const assert = require('node:assert/strict');
const fs = require('node:fs');
const http = require('node:http');
const path = require('node:path');
const playwrightPath = process.env.PLAYWRIGHT_MODULE
    || path.join(__dirname, '..', 'gui', 'node_modules', 'playwright');
const { firefox } = require(playwrightPath);
const root = path.join(__dirname, '..', 'gui', 'wwwroot');
const types = { '.html': 'text/html', '.js': 'text/javascript', '.css': 'text/css',
    '.svg': 'image/svg+xml', '.png': 'image/png', '.woff2': 'font/woff2' };
const check = (ok, label) => { assert.ok(ok, label); console.log('PASS ' + label); };

const row = (jobs, state) => ({ key: 'build%20agent', id: 'build agent', parent: null,
    state: state || 'busy', paused: false, submitted: false, can_delete: true,
    goal: '', phase: 'idle', tool: '', detail: '', jobs });

const server = http.createServer((request, response) => {
    const name = request.url === '/' ? 'index.html' : request.url.replace(/^\.\//, '').split('?')[0];
    const file = path.join(root, name);
    if (!file.startsWith(root + path.sep) || !fs.existsSync(file)) { response.writeHead(404); response.end(); return; }
    response.writeHead(200, { 'content-type': types[path.extname(file)] || 'application/octet-stream' });
    response.end(fs.readFileSync(file));
});

(async () => {
    const browser = await firefox.launch({ headless: true,
        executablePath: process.env.FIREFOX_BINARY || undefined });
    try {
        await new Promise((done) => server.listen(0, '127.0.0.1', done));
        const base = `http://127.0.0.1:${server.address().port}/`;
        const context = await browser.newContext();
        const page = await context.newPage();
        const errors = [];
        page.on('pageerror', error => errors.push(String(error)));
        await page.goto(base, { waitUntil: 'domcontentloaded' });

        const label = (data) => page.evaluate((input) => {
            socket = { readyState: 1 };
            connected = true;
            selectionPending = false;
            agentClosed = false;
            agentID = input.key;
            agentRows = input.rows;
            acknowledgedState = input.state;
            renderControls();
            return document.querySelector('#state').textContent;
        }, data);

        check(await label({ key: 'build%20agent', rows: [row(0)], state: 'idle' }) === 'idle',
            'an Agent with nothing running reads idle');
        check(await label({ key: 'build%20agent', rows: [row(1)], state: 'idle' }) === '1 job running',
            `one live Job is reported after the turn ends: ${await label({ key: 'build%20agent', rows: [row(1)], state: 'idle' })}`);
        check(await label({ key: 'build%20agent', rows: [row(3)], state: 'idle' }) === '3 jobs running',
            'several live Jobs are counted, not just flagged');
        check(await label({ key: 'build%20agent', rows: [row(2, 'idle')], state: 'paused' }) === '2 jobs running',
            'a paused actor does not hide work that is still running');
        check(await label({ key: 'build%20agent', rows: [row(1)], state: 'inflight' }) === 'inflight',
            'a turn in flight keeps the more specific label');
        check(await label({ key: 'build%20agent', rows: [row(1)], state: 'toolspending' }) === 'toolspending',
            'a tool still being observed keeps its own label');
        check(await label({ key: 'other%20agent', rows: [row(4)], state: 'idle' }) === 'idle',
            'another Agent Jobs do not leak into this Agent label');
        check(await label({ key: 'build%20agent', rows: [], state: 'idle' }) === 'idle',
            'a missing tree row is no work, not an error');

        // The tree poll is the only thing that carries the count, so the frame
        // handler must repaint the header; without it the label freezes.
        const viaFrame = await page.evaluate((input) => {
            socket = { readyState: 1 };
            connected = true;
            selectionPending = false;
            agentClosed = false;
            agentID = input.key;
            agentRows = [];
            acknowledgedState = 'idle';
            renderControls();
            const before = document.querySelector('#state').textContent;
            handle({ type: 'agent_tree', agents: input.rows });
            return { before, after: document.querySelector('#state').textContent };
        }, { key: 'build%20agent', rows: [row(2)] });
        check(viaFrame.before === 'idle' && viaFrame.after === '2 jobs running',
            `an agent_tree frame repaints the header: ${JSON.stringify(viaFrame)}`);

        check(errors.length === 0, `no page errors: ${JSON.stringify(errors.slice(0, 3))}`);
        console.log('JOB STATUS SUITE PASSED');
    } finally {
        await browser.close();
        server.close();
    }
})().catch(error => { console.error('FAIL ' + (error && error.stack || error)); process.exit(1); });
