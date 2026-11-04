// Fenced blocks the operator is meant to take away, copied in real Firefox.
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
const markdown = fs.readFileSync(path.join(root, 'markdown.js'), 'utf8');
assert(markdown.includes('Copyable text'), 'the shipped page must offer copyable blocks');

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
        // Firefox grants no clipboard permission at all, so both copy paths are
    // Firefox grants no clipboard permission at all, so both copy paths are
    // observed rather than trusted: the Clipboard API call and the selection fallback.
    const context = await browser.newContext();
        const page = await context.newPage();
        await page.goto(base, { waitUntil: 'domcontentloaded' });
        const snippet = 'ssh-copy-id -i key.pub user@host\n\ntar -czf out.tgz "a dir" \\'
            + '\n  --exclude=cache .';
        const long = Array.from({ length: 40 }, (_, i) => `line-${i} = ${i}; // teaching code`).join('\n');
        await page.evaluate(() => {
            // Firefox grants no clipboard-read, so observe the real API call instead
            // of reading the system clipboard back.
            window.__copied = [];
            const original = navigator.clipboard.writeText.bind(navigator.clipboard);
            navigator.clipboard.writeText = async (text) => {
                window.__copied.push(text);
                return original(text);
            };
        });
        await page.evaluate((([code, big]) => {
            const host = document.createElement('div');
            document.body.append(host);
            markdown(host, 'Text before.\n\n```bash\n' + code + '\n```\n\n'
                + 'Long block:\n\n```javascript\n' + big + '\n```\n');
        }), [snippet, long]);
        const blocks = page.locator('.clipboard-block');
        check(await blocks.count() === 2, 'both fenced blocks became copyable boxes');
        const field = blocks.nth(0).getByRole('textbox', { name: 'Copyable text', exact: true });
        check(await field.inputValue() === snippet,
            'the copyable payload is literal, blank lines and escapes intact');
        check(await blocks.nth(0).getByRole('button', { name: 'Copy text', exact: true }).count() === 1,
            'each box offers exactly one copy control');
        await blocks.nth(0).getByRole('button', { name: 'Copy text' }).click();
        check(JSON.stringify(await page.evaluate(() => window.__copied))
            === JSON.stringify([snippet]),
            'the copy button wrote the exact payload through the Clipboard API');
        const confirmed = await blocks.nth(0).getByRole('button', { name: /copied/i }).count();
        check(confirmed === 1, 'the control confirms the copy');
        check(await blocks.nth(0).getByRole('button', { name: /copied/i })
            .evaluate(button => document.activeElement === button),
            'copying keeps keyboard focus on the control');
        // The fallback used when the Clipboard API is unavailable must select the
        // whole literal payload, or the operator silently copies a fragment.
        await page.evaluate(() => {
            window.__fallback = null;
            navigator.clipboard.writeText = async () => { throw new Error('denied'); };
            document.execCommand = (command) => {
                if (command !== 'copy') return false;
                const active = document.activeElement;
                window.__fallback = active && active.getAttribute('aria-label') === 'Copyable text'
                    ? { text: active.value, start: active.selectionStart, end: active.selectionEnd }
                    : null;
                return true;
            };
        });
        await blocks.nth(0).getByRole('button', { name: 'Copy text' }).click();
        const fallback = await page.evaluate(() => window.__fallback);
        check(fallback && fallback.text === snippet && fallback.start === 0
            && fallback.end === snippet.length,
            `the fallback selects the whole payload: ${JSON.stringify(fallback)}`);
        check(await blocks.nth(0).getByRole('button', { name: /copied/i }).count() === 1,
            'a successful fallback still confirms');
        // A collapsed long block must still hand over the complete payload.
        const big = blocks.nth(1).getByRole('textbox', { name: 'Copyable text', exact: true });
        check((await big.inputValue()).split('\n').length === 40,
            'a collapsed block keeps its complete payload for copying');
        await blocks.nth(1).getByRole('button', { name: 'Copy text' }).click();
        // The API is stubbed to fail here, so this proves the fallback carries the
        // whole block rather than the visible preview.
        const collapsed = await page.evaluate(() => window.__fallback);
        check(collapsed && collapsed.text === long && collapsed.start === 0
            && collapsed.end === long.length,
            `the collapsed block copied complete, not truncated: ${JSON.stringify(collapsed && collapsed.text.slice(-24))}`);
        check(await page.locator('.clipboard-block pre').count() === 2,
            'the readable block is still shown, not replaced');
    } finally {
        await browser.close();
        server.close();
    }
})();
