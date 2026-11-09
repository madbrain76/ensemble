// An endpoint that never asked for a stream answers in one buffered body, and from
// the page it looks like the model thought in silence and then dumped the answer.
// The formatter and timing suites could not see that: they hand the view a finished
// part. This drives the shipped frame handler the way a streaming endpoint does and
// insists the page can be observed mid-answer.
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
        await page.waitForFunction(() => typeof handle === 'function');

        const result = await page.evaluate(async () => {
            agentID = 'root';
            cards.clear();
            document.querySelectorAll('article[data-part-id]').forEach((node) => node.remove());
            const settle = () => new Promise((done) => setTimeout(done, 80));
            const read = (id) => {
                const node = document.querySelector(`article[data-part-id="${id}"] .message`);
                return node ? node.textContent.trim() : '';
            };
            const pieces = ['the ', 'page ', 'must ', 'grow ', 'in place'];
            const steps = [];
            const thinking = [];
            for (let index = 0; index < pieces.length; index++) {
                handle({ type: 'part_delta', agent_id: 'root', part_id: 'grow', chunk: pieces[index], kind: 'text' });
                if (index === 1) {
                    handle({ type: 'part_delta', agent_id: 'root', part_id: 'think', chunk: 'considering ', kind: 'thinking' });
                }
                handle({ type: 'part_delta', agent_id: 'root', part_id: 'think', chunk: 'more. ', kind: 'thinking' });
                await settle();
                steps.push(read('grow'));
                thinking.push(read('think'));
            }
            return { steps, thinking, whole: pieces.join('').trim() };
        });

        const partials = result.steps.filter((text) => text && text !== result.whole);
        check(result.steps[result.steps.length - 1] === result.whole,
            'the settled answer is complete and unduplicated: ' + JSON.stringify(result.steps[result.steps.length - 1]));
        check(partials.length >= 2,
            'the page was observable mid-answer, not filled in one go: ' + JSON.stringify(result.steps));
        check(new Set(result.steps.filter(Boolean)).size === result.steps.length,
            'each fragment advanced the visible text: ' + JSON.stringify(result.steps));
        check(result.steps.every((text, index) => index === 0 || text.startsWith(result.steps[index - 1])),
            'streamed text grew by appending rather than being rewritten: ' + JSON.stringify(result.steps));
        // The thought box is a separate card that also arrives piece by piece;
        // proving it grows is what distinguishes streaming from a final dump.
        check(result.thinking[0].length > 0
            && result.thinking[result.thinking.length - 1].startsWith(result.thinking[0])
            && result.thinking[result.thinking.length - 1].length > result.thinking[0].length,
            'thinking streamed into its own box: ' + JSON.stringify(result.thinking));
        check(errors.length === 0, 'no page errors: ' + (errors.join('; ') || 'none'));
    } finally {
        await browser.close();
        server.close();
    }
})().catch((error) => { console.error('FAIL', error); process.exit(1); });
