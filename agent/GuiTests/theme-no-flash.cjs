// The confirmed theme must be painted before the stylesheet, not after first paint.
const assert = require('node:assert/strict');
const fs = require('node:fs');
const http = require('node:http');
const path = require('node:path');
// Inlined so this suite stands on its own branch; the Firefox tier centralises the
// same resolution in browser-launch.cjs.
const playwrightPath = process.env.PLAYWRIGHT_MODULE
    || path.join(__dirname, '..', 'gui', 'node_modules', 'playwright');
const { firefox } = require(playwrightPath);
const launchBrowser = () => firefox.launch({ headless: true,
    executablePath: process.env.FIREFOX_BINARY || undefined });
const root = path.join(__dirname, '..', 'gui', 'wwwroot');
const types = { '.html': 'text/html', '.js': 'text/javascript', '.css': 'text/css',
    '.svg': 'image/svg+xml', '.png': 'image/png', '.woff2': 'font/woff2' };
const check = (ok, label) => { assert.ok(ok, label); console.log('PASS ' + label); };

(async () => {
const server = http.createServer((request, response) => {
    const name = request.url === '/' ? 'index.html' : request.url.replace(/^\.\//, '').split('?')[0];
    const file = path.join(root, name);
    if (!file.startsWith(root + path.sep) || !fs.existsSync(file)) { response.writeHead(404); response.end(); return; }
    response.writeHead(200, { 'content-type': types[path.extname(file)] || 'application/octet-stream' });
    response.end(fs.readFileSync(file));
});
const browser = await launchBrowser();
try {
    await new Promise((done) => server.listen(0, '127.0.0.1', done));
    const base = `http://127.0.0.1:${server.address().port}/`;
    const html = fs.readFileSync(path.join(root, 'index.html'), 'utf8');
    const painted = html.indexOf('ensemble.theme.v2');
    const stylesheet = html.indexOf('<link rel="stylesheet"');
    check(painted >= 0 && painted < stylesheet,
        'the theme is painted before the stylesheet is requested');
    check(html.indexOf('<script>') < stylesheet, 'the early paint is not deferred');
    check(/<option value="system" selected>/.test(html)
        && !/<option value="dark" selected>/.test(html),
        'the theme select preselects System, matching the resolved default');

    // The flash window is the gap between the first frame and the stylesheet
    // arriving. Holding the stylesheet back makes that window observable: whatever
    // the theme attribute says while CSS is pending is what would have been painted.
    const probe = async (stored, colorScheme) => {
        const context = await browser.newContext({ colorScheme });
        const page = await context.newPage();
        await page.route('**/vendor/bootstrap-5.3.8.min.css', async (route) => {
            await new Promise((done) => setTimeout(done, 500));
            await route.continue();
        });
        await page.goto(base, { waitUntil: 'domcontentloaded' });
        await page.evaluate((value) => {
            localStorage.clear();
            if (value !== null) localStorage.setItem('ensemble.theme.v2', value);
        }, stored);
        await page.reload({ waitUntil: 'commit' });
        const pending = await page.evaluate(() => new Promise((done) => setTimeout(
            () => done({ theme: document.documentElement.dataset.bsTheme,
                sheets: document.styleSheets.length }), 200)));
        const settled = await page.evaluate(async () => {
            await new Promise((done) => setTimeout(done, 800));
            return { theme: document.documentElement.dataset.bsTheme,
                sheets: document.styleSheets.length };
        });
        await context.close();
        return { pending, settled };
    };

    const light = await probe('light', 'dark');
    // sheets < 2 proves Bootstrap had not applied yet when the theme already read light.
    check(light.pending.theme === 'light' && light.pending.sheets < 2
        && light.settled.theme === 'light' && light.settled.sheets >= 2,
        `a stored light profile must not paint dark while the stylesheet is pending: ${JSON.stringify(light)}`);
    const dark = await probe('dark', 'light');
    check(dark.pending.theme === 'dark' && dark.settled.theme === 'dark',
        `a stored dark profile stays dark: ${JSON.stringify(dark)}`);
    const system = await probe('system', 'light');
    check(system.pending.theme === 'light' && system.settled.theme === 'light',
        `system follows the browser's own preference: ${JSON.stringify(system)}`);
    const junk = await probe('translucent-sepia', 'light');
    // An unrecognised value is not trusted as a choice, but it is also not a reason
    // to paint dark: theme.js resolves the same value to system, so matching it here
    // is what removes the flash rather than hiding it.
    check(junk.pending.theme === 'light' && junk.settled.theme === 'light',
        `an unrecognised stored value follows the operating system: ${JSON.stringify(junk)}`);

    // The first visit on a browser profile: nothing stored at all. This is the case
    // that shipped dark for OS-light users.
    const firstRun = await probe(null, 'light');
    check(firstRun.pending.theme === 'light' && firstRun.pending.sheets < 2
        && firstRun.settled.theme === 'light' && firstRun.settled.sheets >= 2,
        `a first visit on an OS-light profile must not paint dark: ${JSON.stringify(firstRun)}`);
    const firstRunDark = await probe(null, 'dark');
    check(firstRunDark.pending.theme === 'dark' && firstRunDark.settled.theme === 'dark',
        `a first visit on an OS-dark profile paints dark: ${JSON.stringify(firstRunDark)}`);

    // Bill's profile: the old seed path recorded "dark" under v1 while the server
    // default was dark, so a real OS-light user was locked to dark. The stored v1
    // value must neither paint dark nor survive the visit.
    const poisoned = await browser.newContext({ colorScheme: 'light' });
    const poisonedPage = await poisoned.newPage();
    await poisonedPage.goto(base, { waitUntil: 'domcontentloaded' });
    await poisonedPage.evaluate(() => localStorage.setItem('ensemble.theme.v1', 'dark'));
    await poisonedPage.reload({ waitUntil: 'commit' });
    const poisonedPaint = await poisonedPage.evaluate(() => new Promise((done) => setTimeout(
        () => done({ theme: document.documentElement.dataset.bsTheme,
            v1: localStorage.getItem('ensemble.theme.v1'),
            v2: localStorage.getItem('ensemble.theme.v2') }), 200)));
    await poisoned.close();
    check(poisonedPaint.theme === 'light' && poisonedPaint.v1 === null && poisonedPaint.v2 === null,
        `a profile poisoned by the old dark seed follows the operating system and loses v1: ${JSON.stringify(poisonedPaint)}`);
} finally {
    await browser.close();
    server.close();
}
})();
