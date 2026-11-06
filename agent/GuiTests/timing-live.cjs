// The formatter suites prove what a duration would look like; neither could catch
// that the live stream never asked for one. `part_final` is emitted only by the
// replay projection, so the frames a running turn actually receives are deltas --
// and the timing rows appeared only after a reload. These drive the shipped frame
// handler, because the defect lived in the call, not the computation.
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
    const errors = [];
    try {
        await new Promise((done) => server.listen(0, '127.0.0.1', done));
        const base = `http://127.0.0.1:${server.address().port}/`;
        const page = await (await browser.newContext()).newPage();
        page.on('pageerror', (error) => errors.push(String(error)));
        await page.goto(base, { waitUntil: 'domcontentloaded' });
        await page.waitForFunction(() => typeof handle === 'function'
            && typeof beginTurn === 'function');

        const live = await page.evaluate(async () => {
            agentID = 'root';
            cards.clear();
            // A running turn only ever sends deltas; nothing marks a part finished.
            handle({ type: 'part_delta', agent_id: 'root', part_id: 'answer',
                chunk: 'working on it', kind: 'text' });
            const line = () => {
                const node = document.querySelector('article .card-timing');
                return node ? node.textContent.trim() : '';
            };
            const streaming = line();
            beginTurn();
            handle({ type: 'part_delta', agent_id: 'root', part_id: 'answer',
                chunk: ' some more', kind: 'text' });
            handle({ type: 'agent_state', agent_id: 'root', state: 'inflight' });
            await new Promise((wait) => setTimeout(wait, 260));
            handle({ type: 'agent_state', agent_id: 'root', state: 'idle' });
            const summary = document.querySelector('.turn-timing');
            return { streaming, settled: line(),
                summary: summary ? summary.textContent.replace(/\s+/g, ' ').trim() : '' };
        });

        // While the part is still arriving there is no honest duration to show, so
        // it carries the clock and an ellipsis; the measurement belongs to the close.
        check(/^\d\d:\d\d:\d\d · …$/.test(live.streaming),
            `a part in flight shows its arrival time and an honest ellipsis: '${live.streaming}'`);
        check(/^\d\d:\d\d:\d\d · \d\.\d+s$/.test(live.settled),
            `the same part carries a measured duration once the turn closes: '${live.settled}'`);
        check(/Total (?!0\.0s)\d/.test(live.summary) && /Model (?!0\.0s)\d/.test(live.summary),
            `the closed turn reports a measured split, not zeros: ${live.summary}`);

        // A sibling Agent's state must not close this Agent's turn.
        const isolation = await page.evaluate(async () => {
            agentID = 'root';
            cards.clear();
            beginTurn();
            handle({ type: 'part_delta', agent_id: 'root', part_id: 'x', chunk: 'hi', kind: 'text' });
            handle({ type: 'agent_state', agent_id: 'other%20agent', state: 'idle' });
            return { active: turn.active, summary: (document.querySelector('.turn-timing') || {}).textContent || '' };
        });
        check(isolation.active === true,
            "another Agent's idle frame does not close this Agent's turn");

        // A turn that ends before any live output must not claim a measurement.
        const early = await page.evaluate(async () => {
            agentID = 'root';
            cards.clear();
            document.querySelectorAll('.turn-timing').forEach((node) => node.remove());
            beginTurn();
            handle({ type: 'agent_state', agent_id: 'root', state: 'paused' });
            handle({ type: 'agent_state', agent_id: 'root', state: 'inputpending' });
            return { active: turn.active };
        });
        check(early.active === true,
            'a turn that never produced output is not reported as finished at zero');

        check(errors.length === 0, 'no page errors: ' + (errors.join('; ') || 'none'));
    } finally {
        await browser.close();
        server.close();
    }
})().catch((error) => { console.error('FAIL', error); process.exit(1); });
