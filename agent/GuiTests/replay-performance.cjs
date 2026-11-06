// Selecting an Agent with a large history used to freeze the page: replaying
// thousands of parts read scrollHeight once per card, which is quadratic, and then
// attached them all in one task. Firefox offered to stop the page. These drive the
// shipped replay path and bound the longest task the browser is asked to run.
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
const PARTS = 1200;

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
        const page = await (await browser.newContext()).newPage();
        page.on('pageerror', (error) => errors.push(String(error)));
        await page.goto(`http://127.0.0.1:${server.address().port}/`, { waitUntil: 'domcontentloaded' });
        await page.waitForFunction(() => typeof handle === 'function');

        const replay = await page.evaluate(async (count) => {
            const events = [];
            for (let index = 0; index < count; index += 1) {
                events.push({ type: 'part_final', agent_id: 'perf%20agent',
                    part_id: 'p' + index,
                    part: { type: 'text', role: 'assistant',
                        text: 'part ' + String(index).padStart(5, '0') + ' of the transcript' } });
            }
            const gaps = [];
            let last = performance.now();
            const timer = setInterval(() => {
                const now = performance.now();
                gaps.push(now - last);
                last = now;
            }, 10);
            const started = performance.now();
            handle({ type: 'replay', agent_id: 'perf%20agent', data: { events } });
            const attachedAt = performance.now() - started;
            // The tail must be readable immediately; history may keep arriving.
            let visible = 0;
            for (let attempt = 0; attempt < 400; attempt += 1) {
                visible = document.querySelectorAll('article').length;
                if (visible >= count) break;
                await new Promise((wait) => requestAnimationFrame(wait));
            }
            clearInterval(timer);
            const nodes = [...document.querySelectorAll('article .message, article')].filter((node) =>
                node.classList.contains('part'));
            const texts = nodes.map((node) => (node.textContent.match(/part (\d{5})/) || [])[1] || '');
            return {
                attachedAt: Math.round(attachedAt), totalAt: Math.round(performance.now() - started),
                visible, maxGap: Math.round(Math.max(...gaps, 0)), gapSamples: gaps.length,
                duplicates: texts.length - new Set(texts).size,
                first: texts[0], last: texts[texts.length - 1],
                ordered: texts.every((value, index) => index === 0 || texts[index - 1] <= value),
            };
        }, PARTS);

        check(replay.visible === PARTS,
            `every part is eventually attached: ${replay.visible} of ${PARTS}`);
        check(replay.ordered && replay.duplicates === 0,
            `history keeps arrival order without duplicates (first ${replay.first}, last ${replay.last}, `
                + `dupes ${replay.duplicates})`);
        // The claim being protected: the pane is usable quickly and no single task
        // holds the browser for seconds. Measured before the fix: a 58.0s task for
        // 3,624 cards, i.e. roughly 15ms of layout per card.
        // An interval cannot fire inside a blocking task, so a transcript that
        // stalls the thread for a minute records no gaps at all. Require evidence
        // that the watchdog actually ran before trusting its reading.
        check(replay.gapSamples > 5,
            `the long-task watchdog observed the replay: ${replay.gapSamples} samples`);
        check(replay.attachedAt < 500,
            `the synchronous replay task stays short enough to feel instant: ${replay.attachedAt}ms`);
        check(replay.maxGap < 1000,
            `no task blocks the main thread long enough to look frozen: ${replay.maxGap}ms`);
        check(errors.length === 0, 'no page errors: ' + (errors.join('; ') || 'none'));
    } finally {
        await browser.close();
        server.close();
    }
})().catch((error) => { console.error('FAIL', error); process.exit(1); });
