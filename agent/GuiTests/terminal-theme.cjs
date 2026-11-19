'use strict';
// Real CSS in a real browser: the command log must follow the theme, and every ANSI
// colour must stay legible on whichever ground that theme picks. Contrast is measured
// here rather than assumed, so a palette tuned for one ground cannot silently rot.
const assert = require('node:assert/strict');
const path = require('node:path');
const { launchBrowser } = require('./browser-launch.cjs');
const root = path.resolve(__dirname, '../gui/wwwroot');
function luminance(color) {
    const channels = color.match(/[\d.]+/g).slice(0, 3).map(Number).map(value => {
        value /= 255;
        return value <= .04045 ? value / 12.92 : ((value + .055) / 1.055) ** 2.4;
    });
    return channels[0] * .2126 + channels[1] * .7152 + channels[2] * .0722;
}
function contrast(a, b) {
    const values = [luminance(a), luminance(b)].sort((x, y) => y - x);
    return (values[0] + .05) / (values[1] + .05);
}
const expected = {
    dark: { ground: 'rgb(20, 23, 27)', ink: 'rgb(231, 233, 238)' },
    light: { ground: 'rgb(239, 241, 245)', ink: 'rgb(23, 37, 59)' }
};
(async () => {
    const browser = await launchBrowser({ headless: true });
    try {
        const page = await browser.newPage();
        const cells = [];
        for (let index = 0; index < 16; index++) {
            cells.push(`<span class="terminal-fg-${index}" data-fg="${index}">f</span>`);
            cells.push(`<span class="terminal-bg-${index}" data-bg="${index}">b</span>`);
        }
        await page.setContent(`<main id="tools-pane"><section class="job-session part">`
            + `<pre class="terminal-output">${cells.join('')}</pre></section></main>`);
        await page.addStyleTag({ path: path.join(root, 'vendor/bootstrap-5.3.8.min.css') });
        await page.addStyleTag({ path: path.join(root, 'style.css') });
        const grounds = {};
        for (const theme of Object.keys(expected)) {
            await page.evaluate(value => { document.documentElement.dataset.bsTheme = value; }, theme);
            const measured = await page.evaluate(() => {
                const log = document.querySelector('.terminal-output');
                return {
                    ground: getComputedStyle(log).backgroundColor,
                    ink: getComputedStyle(log).color,
                    pane: getComputedStyle(document.querySelector('.job-session')).backgroundColor,
                    foreground: [...document.querySelectorAll('[data-fg]')].map(entry =>
                        [entry.dataset.fg, getComputedStyle(entry).color]),
                    background: [...document.querySelectorAll('[data-bg]')].map(entry =>
                        [entry.dataset.bg, getComputedStyle(entry).backgroundColor, getComputedStyle(entry).color])
                };
            });
            grounds[theme] = measured.ground;
            assert.equal(measured.ground, expected[theme].ground, `${theme}: command-log ground`);
            assert.equal(measured.ink, expected[theme].ink, `${theme}: command-log ink`);
            assert.ok(contrast(measured.ink, measured.ground) >= 4.5, `${theme}: unstyled log text`);
            // A block identical to the plain page would trade one theme fault for another.
            assert.notEqual(measured.ground, 'rgb(255, 255, 255)', `${theme}: log no longer reads as a block`);
            for (const [index, color] of measured.foreground) {
                const ratio = contrast(color, measured.ground);
                assert.ok(ratio >= 4.5, `${theme}: ANSI foreground ${index} at ${ratio.toFixed(2)}:1`);
            }
            for (const [index, color, ink] of measured.background) {
                const ratio = contrast(ink, color);
                assert.ok(ratio >= 4.5, `${theme}: ANSI background cell ${index} at ${ratio.toFixed(2)}:1`);
                assert.notEqual(color, measured.ground, `${theme}: cell ${index} invisible on the ground`);
            }
            const weakest = Math.min(...measured.foreground.map(([, color]) => contrast(color, measured.ground)));
            const box = await page.locator('.terminal-output').boundingBox();
            assert.ok(box && box.height > 40 && box.width > 200, `${theme}: command log has collapsed`);
            console.log(`PASS ${theme} command log: ground follows theme, weakest ANSI colour ${weakest.toFixed(2)}:1`);
        }
        assert.notEqual(grounds.dark, grounds.light, 'The command log ignores the theme.');
    } finally { await browser.close(); }
})().catch(error => { console.error(error); process.exitCode = 1; });
