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
    const painted = html.indexOf('ensemble.theme.v1');
    const stylesheet = html.indexOf('<link rel="stylesheet"');
    check(painted >= 0 && painted < stylesheet,
        'the theme is painted before the stylesheet is requested');
    check(html.indexOf('<script>') < stylesheet, 'the early paint is not deferred');

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
        if (stored !== null) {
            await page.evaluate((value) => localStorage.setItem('ensemble.theme.v1', value), stored);
        } else {
            await page.evaluate(() => localStorage.removeItem('ensemble.theme.v1'));
        }
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
    check(junk.pending.theme === 'dark',
        `an unrecognised stored value must not be trusted by the early paint: ${JSON.stringify(junk)}`);
} finally {
    await browser.close();
    server.close();
}
})();
