// The formatter suites proved what a duration would look like; none could catch
// that the live stream never asked for one. `part_final` is emitted only by the
// replay projection, so the frames a running turn actually receives are deltas,
// and a card's production span is therefore painted from what this tab watched.
// These drive the shipped frame handler and the shipped close, because the defect
// lived in the call, not the computation. The replay half pins the opposite rule:
// a recorded clock and a recorded span must outrank anything the tab could read.
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

// Read one card's two timing elements. Inlined in each page.evaluate so the
// page never has to eval a string assembled in Node.
const READ = `(id) => {
    const node = document.querySelector('article[data-part-id="' + id + '"]');
    const stamp = node && node.querySelector(':scope > .recorded-time');
    const duration = node && node.querySelector(':scope > .observation-duration');
    return { stamp: stamp ? stamp.textContent.trim() : null,
        dateTime: stamp ? stamp.getAttribute('datetime') : null,
        stampTitle: stamp ? stamp.getAttribute('title') : null,
        duration: duration ? duration.textContent.trim() : null,
        durationTitle: duration ? duration.getAttribute('title') : null };
}`;

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

        // A running turn sends only deltas, so while a part is arriving there is
        // no honest clock and no honest duration: the measurement belongs to the
        // close. Claiming either early would be invented data.
        const live = await page.evaluate(async (read) => {
            agentID = 'root';
            cards.clear();
            document.querySelectorAll('.turn-timing').forEach((node) => node.remove());
            handle({ type: 'part_delta', agent_id: 'root', part_id: 'answer',
                chunk: 'working on it', kind: 'text' });
            const part = eval(read);
            const streaming = part('answer');
            beginTurn();
            handle({ type: 'part_delta', agent_id: 'root', part_id: 'answer',
                chunk: ' some more', kind: 'text' });
            handle({ type: 'agent_state', agent_id: 'root', state: 'inflight' });
            await new Promise((wait) => setTimeout(wait, 260));
            handle({ type: 'agent_state', agent_id: 'root', state: 'idle' });
            const summary = document.querySelector('.turn-timing');
            return { streaming, settled: part('answer'),
                summary: summary ? summary.textContent.replace(/\s+/g, ' ').trim() : '' };
        }, READ);

        check(live.streaming.stamp === null && live.streaming.duration === null,
            `a part still in flight claims neither a clock nor a duration: ${JSON.stringify(live.streaming)}`);
        // The closed card must carry a measured duration without a reload. This is
        // the regression this suite exists for: the span was computed and never
        // painted, so the row only appeared after the page was reopened.
        check(live.settled.duration !== null && /[\d.]+ s$/.test(live.settled.duration),
            `the closed part carries a production duration in the live turn: ${JSON.stringify(live.settled)}`);
        check(/measured in this tab/.test(live.settled.durationTitle || ''),
            `a tab-measured span names its own source: ${live.settled.durationTitle}`);
        check(live.settled.stamp === null,
            'a live part invents no recorded instant, because no frame supplied one');
        check(/Total (?!0\.0s)\d/.test(live.summary) && /Model (?!0\.0s)\d/.test(live.summary),
            `the closed turn reports a measured split, not zeros: ${live.summary}`);

        // A sibling Agent's state must not close this Agent's turn.
        const isolation = await page.evaluate(() => {
            agentID = 'root';
            cards.clear();
            beginTurn();
            handle({ type: 'part_delta', agent_id: 'root', part_id: 'x', chunk: 'hi', kind: 'text' });
            handle({ type: 'agent_state', agent_id: 'other%20agent', state: 'idle' });
            return { active: turn.active };
        });
        check(isolation.active === true,
            "another Agent's idle frame does not close this Agent's turn");

        // A turn that ends before any live output must not claim a measurement.
        const early = await page.evaluate(() => {
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

        // The complaint that produced the clock rule: a replayed transcript stamped
        // itself with the moment the tab opened. Every reading the page could take
        // for itself is pushed 400 days away here, so a correct clock can only have
        // come from the frame the server sent.
        const recorded = await page.evaluate((read) => {
            agentID = 'root';
            cards.clear();
            document.querySelectorAll('article[data-part-id]').forEach((node) => node.remove());
            const realNow = Date.now;
            Date.now = () => realNow() + 400 * 24 * 3600 * 1000;
            const closed = new Date(realNow() - 3 * 3600 * 1000 - 90 * 1000);
            const iso = closed.toISOString();
            const began = closed.getTime() - 2000;
            handle({ type: 'part_final', agent_id: 'root', part_id: 'rec-reply', recorded_at: iso,
                part: { type: 'text', role: 'assistant', text: 'recorded reply',
                    started: began, elapsed: 4200 } });
            // A tool result's instant is its frame time; its span is the server's
            // dispatch-to-return observation, which is not background Job runtime.
            handle({ type: 'tool_finished', agent_id: 'root', part_id: 'rec-result',
                recorded_at: iso, observation_ms: 1234, is_error: false,
                parts: [{ type: 'text', text: 'tool output' }] });
            // A finalized part the provider never streamed has no span to report.
            // An invented 0.0s would be as false as a page-open timestamp.
            handle({ type: 'part_final', agent_id: 'root', part_id: 'rec-plain', recorded_at: iso,
                part: { type: 'text', role: 'assistant', text: 'never streamed' } });
            const part = eval(read);
            const out = { reply: part('rec-reply'), result: part('rec-result'),
                plain: part('rec-plain'), iso,
                want: new Date(iso).toLocaleTimeString() };
            Date.now = realNow;
            return out;
        }, READ);

        check(recorded.reply.stamp === recorded.want
                && recorded.reply.dateTime === recorded.iso,
            `a replayed reply keeps the recorded clock, not the tab's: `
            + `'${recorded.reply.stamp}' want '${recorded.want}'`);
        check(/·\s*4\.2 s$/.test(recorded.reply.duration || ''),
            `a recorded span is shown even though the tab measured nothing: ${recorded.reply.duration}`);
        check(/recorded by the server/.test(recorded.reply.durationTitle || ''),
            `a server-vouched span names its source: ${recorded.reply.durationTitle}`);
        check(/·\s*1\.234 s$/.test(recorded.result.duration || ''),
            `a tool result is timed dispatch to return: ${recorded.result.duration}`);
        check(/not background Job runtime/.test(recorded.result.durationTitle || ''),
            `an observation span says what it is not: ${recorded.result.durationTitle}`);
        check(recorded.plain.stamp === recorded.want && recorded.plain.duration === null,
            `an unmeasured finalized part shows a clock alone: ${JSON.stringify(recorded.plain)}`);

        check(errors.length === 0, 'no page errors: ' + (errors.join('; ') || 'none'));
    } finally {
        await browser.close();
        server.close();
    }
})().catch((error) => { console.error('FAIL', error); process.exit(1); });
