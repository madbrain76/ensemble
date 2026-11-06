// The Agent picker and the capture consent dialog show a name, not a wire address.
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

// Names as the server actually sends them: key is the address, id is the name.
const rows = [
    { key: '', id: 'root', parent: null, state: 'idle', can_delete: false },
    { key: 'ensemble%20debugger', id: 'ensemble debugger', parent: null, state: 'idle', can_delete: true },
    { key: 'PDF%20designer', id: 'PDF designer', parent: 'ensemble%20debugger', state: 'running', tool: 'run_command', can_delete: true },
    { key: null, id: 'fleet header', parent: null, state: 'idle', can_delete: false },
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
        const seen = await page.evaluate((data) => {
            agentRows = data;
            agentID = 'PDF%20designer';
            showAgents();
            return {
                options: [...document.querySelectorAll('#agent-choice option')]
                    .map(option => option.textContent),
                selectedText: document.querySelector('#agent-choice').selectedOptions[0]?.textContent || '',
                placeholder: document.querySelector('#prompt').placeholder
            };
        }, rows);
        console.log('options: ' + JSON.stringify(seen.options));
        check(seen.options.some(text => text.startsWith('ensemble debugger —')),
            'a spaced Agent name is shown as typed, not percent-encoded');
        check(seen.options.some(text => text.startsWith('PDF designer — running — run_command')),
            'activity still follows the name');
        check(!seen.options.some(text => text.includes('%20')),
            `no option leaks a wire address: ${JSON.stringify(seen.options)}`);
        check(!seen.selectedText.includes('%20'),
            `the selected text the consent dialog reads is a name: ${JSON.stringify(seen.selectedText)}`);
        check(seen.placeholder === 'Ask PDF designer anything…',
            `the composer addresses the Agent by name: ${JSON.stringify(seen.placeholder)}`);
        check(seen.options.length === rows.length, 'every row still appears in the picker');
        check(errors.length === 0, `no page errors (${JSON.stringify(errors).slice(0, 200)})`);
    } finally {
        await browser.close();
        server.close();
    }
})();
