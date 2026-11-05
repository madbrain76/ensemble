// Tree arrangement: drag to reorder siblings, one action set for the selected
// Agent, a confirmation before deletion, and hotkeys that do not fight the browser.
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
const rows = [
    { key: null, id: 'fleet', parent: null, state: 'idle', can_delete: false },
    { key: 'a1', id: 'alpha', parent: null, state: 'idle', can_delete: true },
    { key: 'a2', id: 'beta', parent: null, state: 'idle', can_delete: true },
    { key: 'a3', id: 'gamma', parent: null, state: 'idle', can_delete: true },
    { key: 'b1', id: 'child', parent: 'a1', state: 'idle', can_delete: true },
];
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
        await page.evaluate((data) => {
            window.confirm = () => false;
            agentRows = data;
            agentID = 'a2';
            renderAgentTree();
        }, rows);

        check(await page.locator('#agent-tree .tree-delete').count() === 0,
            'no delete button is repeated on every row');
        const bar = page.locator('#agent-actions');
        check(await bar.count() === 1, 'one action toolbar exists for the whole tree');
        check(await bar.locator('button[data-delete-agent="a2"]').count() === 1
            && await bar.locator('button[data-delete-agent="a2"]').isDisabled() === false,
            'the action set acts on the selected Agent');

        // Confirmation first: a declined dialog must delete nothing.
        const asked = await page.evaluate(async () => {
            let prompt = null;
            window.confirm = (message) => { prompt = message; return false; };
            const sent = [];
            const original = WebSocket.prototype.send;
            WebSocket.prototype.send = (body) => sent.push(body);
            document.querySelector('#agent-actions button[data-delete-agent]').click();
            await new Promise((done) => setTimeout(done, 60));
            WebSocket.prototype.send = original;
            return { prompt, sent };
        });
        check(asked.prompt && asked.prompt.includes('a2'),
            `delete asks for confirmation naming the Agent: ${asked.prompt}`);
        check(asked.sent.every(body => !body.includes('delete')),
            'a declined confirmation deleted nothing');

        // Drag to reorder, within the same parent only.
        const dragged = await page.evaluate(async () => {
            const store = new Map();
            const transfer = {
                setData: (kind, value) => store.set(kind, value),
                getData: (kind) => store.get(kind) || '',
            };
            const row = (key) => document.querySelector(`[data-drag-agent="${key}"]`);
            const fire = (node, type, extra = {}) => node.dispatchEvent(Object.assign(
                new Event(type, { bubbles: true, cancelable: true }), { dataTransfer: transfer }, extra));
            store.set('text/ensemble-agent', 'a1');
            fire(row('a3'), 'dragover');
            fire(row('a3'), 'drop');
            await new Promise((done) => setTimeout(done, 30));
            const order = [...document.querySelectorAll('#agent-tree .tree-select')]
                .map(button => button.dataset.agent).filter(Boolean);
            // A child must not be adoptable into the root list.
            store.set('text/ensemble-agent', 'b1');
            fire(row('a2'), 'drop');
            await new Promise((done) => setTimeout(done, 30));
            const after = [...document.querySelectorAll('#agent-tree .tree-select')]
                .map(button => button.dataset.agent).filter(Boolean);
            return { order, after, saved: localStorage.getItem('ensemble.agentOrder.v1') };
        });
        // Dropping on a row takes that row's place, pushing it down.
        check(dragged.order.join(',') === 'fleet,a2,a1,b1,a3',
            `siblings reorder by dropping: ${dragged.order.join(',')}`);
        check(dragged.after.join(',') === dragged.order.join(','),
            'a child cannot be dragged into another parent');
        check(dragged.saved && dragged.saved.includes('a2'),
            'the arrangement is remembered for this browser');

        // Reload: the order survives, because it is a local preference.
        await page.reload({ waitUntil: 'domcontentloaded' });
        const persisted = await page.evaluate(async (data) => {
            agentRows = data;
            agentID = 'a2';
            renderAgentTree();
            await new Promise((done) => setTimeout(done, 30));
            return [...document.querySelectorAll('#agent-tree .tree-select')]
                .map(button => button.dataset.agent).filter(Boolean);
        }, rows);
        check(persisted.join(',') === 'fleet,a2,a1,b1,a3',
            `the arrangement survives reload: ${persisted.join(',')}`);

        // Hotkeys: Alt+Shift+Arrow walks the visible rows and moves focus. Plain
        // arrows, and any Alt+Shift press while typing, must be left alone.
        const keys = await page.evaluate(async () => {
            const press = (node, key, init) => node.dispatchEvent(new KeyboardEvent('keydown',
                Object.assign({ key, bubbles: true, cancelable: true }, init)));
            const tree = document.querySelector('#agent-tree');
            const results = {};
            // Re-query each time: selecting re-renders the tree, and a node captured
            // before that render can no longer receive focus.
            const rows = () => [...tree.querySelectorAll('.tree-select:not(:disabled)')];
            rows()[0].focus();
            press(document, 'ArrowDown', { altKey: true, shiftKey: true });
            await new Promise((done) => setTimeout(done, 20));
            results.afterFirst = document.activeElement.dataset.agent;
            press(document, 'ArrowDown', { altKey: true, shiftKey: true });
            await new Promise((done) => setTimeout(done, 20));
            results.afterSecond = document.activeElement.dataset.agent;
            // Two presses from the same starting row advance two rows.
            rows()[0].focus();
            press(document, 'ArrowDown', { altKey: true, shiftKey: true });
            await new Promise((done) => setTimeout(done, 20));
            results.twoSteps = document.activeElement.dataset.agent;
            const input = document.querySelector('#prompt');
            input.focus();
            press(input, 'ArrowDown', { altKey: true, shiftKey: true });
            await new Promise((done) => setTimeout(done, 20));
            results.stillTyping = document.activeElement === input;
            return results;
        });
        const visibleOrder = ['a2', 'a1', 'b1', 'a3'];
        // Two presses from the first row land two rows down; one press lands one down.
        check(keys.afterFirst === visibleOrder[1] && keys.afterSecond === visibleOrder[2]
            && keys.twoSteps === visibleOrder[1],
            `Alt+Shift+Arrow walks the visible rows: ${JSON.stringify(keys)}`);
        check(keys.stillTyping, 'typing in the prompt keeps its arrow keys');
        check(errors.length === 0, 'the tree script stayed error-free: ' + errors.join('; '));
    } finally {
        await browser.close();
        server.close();
    }
})();
