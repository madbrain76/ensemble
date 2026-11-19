'use strict';
// Synthetic real-browser component tests: no clipboard, microphone or provider.
const { launchBrowser, launchPersistentProfile } = require('./browser-launch.cjs');
const path = require('node:path');
const assert = require('node:assert/strict');
(async () => {
  const browser = await launchBrowser({ executablePath: '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome', headless: true });
  try {
    const page = await browser.newPage({ viewport: { width: 393, height: 851 } });
    await page.setContent('<textarea id="text"></textarea><div id="previews"></div><p id="status"></p><button id="discard"></button><input id="files" type="file" multiple>');
    await page.addScriptTag({ path: path.resolve(__dirname, '../gui/wwwroot/images.js') });
    await page.evaluate(() => {
      window.notices = []; window.frames = [];
      window.connection = { send: text => frames.push(JSON.parse(text)) };
      window.composer = new ImageComposer(document.querySelector('#text'), document.querySelector('#previews'), document.querySelector('#status'), document.querySelector('#discard'), () => {}, text => notices.push(text));
      document.querySelector('#files').onchange = () => composer.acceptPicker(document.querySelector('#files'), connection, 'a');
    });
    const png = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+aF1cAAAAASUVORK5CYII=', 'base64');
    const file = { name: 'synthetic.png', mimeType: 'image/png', buffer: png };
    async function choose(files, invalidate) {
      const chooser = page.waitForEvent('filechooser');
      await page.evaluate(() => composer.openPicker(document.querySelector('#files'), connection, 'a'));
      const dialog = await chooser;
      if (invalidate) await page.evaluate(invalidate);
      await dialog.setFiles(files);
    }
    await page.locator('#text').fill('caption');
    await choose(file);
    await page.waitForFunction(() => composer.hasImages && !composer.blocked);
    assert.equal(await page.locator('#previews img').count(), 1);
    await page.evaluate(() => { composer.send(connection, 'a'); composer.send(connection, 'a'); });
    const frames = await page.evaluate(() => frames);
    assert.equal(frames.length, 1); assert.equal(frames[0].text, 'caption');
    assert.equal(frames[0].agent_id, 'a');
    assert.equal(frames[0].images[0].data, png.toString('base64'));
    await page.evaluate(() => composer.acknowledge({ type: 'image_accepted', sequence: 1, agent_id: 'a', images: [] }, connection));
    assert.equal(await page.locator('#text').inputValue(), '');
    await choose(file);
    await page.waitForFunction(() => composer.hasImages && !composer.blocked);
    await page.locator('#previews button').click();
    assert.equal(await page.locator('#previews img').count(), 0);
    console.log('PASS picker preview/remove, exact payload, caption, send once and same-file reselection');
    await choose(file, () => { composer.selectionChanged('b'); composer.selectionChanged('a'); });
    assert.equal(await page.evaluate(() => composer.hasImages), false);
    await choose(file, () => composer.disconnect());
    assert.equal(await page.evaluate(() => composer.hasImages), false);
    await choose(file, () => { window.connection = { send: () => { throw Error('unexpected'); } }; });
    assert.equal(await page.evaluate(() => composer.hasImages), false);
    console.log('PASS stale picker after switch-away-and-back, disconnect and replacement connection');
    const noticeCount = await page.evaluate(() => notices.length);
    await choose(Array(5).fill(file));
    assert.equal(await page.evaluate(() => composer.hasImages), false);
    await choose({ name: 'not-image.txt', mimeType: 'text/plain', buffer: Buffer.from('not pixels') });
    assert.equal(await page.evaluate(() => composer.hasImages), false);
    await choose({ name: 'empty.png', mimeType: 'image/png', buffer: Buffer.alloc(0) });
    assert.equal(await page.evaluate(() => composer.hasImages), false);
    assert.equal(await page.evaluate(() => notices.length), noticeCount + 3);
    console.log('PASS whole-selection count/type/empty refusals');
    await page.evaluate(() => {
      window.connection = { send: text => frames.push(JSON.parse(text)) };
      composer.addFiles([{ size: 64 * 1024 * 1024 + 1, type: 'image/png' }], 'a');
    });
    assert.equal(await page.evaluate(() => composer.hasImages), false);
    assert.match(await page.evaluate(() => notices.at(-1)), /64 MiB/);
    // Padding exercises browser source-byte admission/read, not raster validity;
    // the native preparation test separately uses a real >32-MiB encoded image.
    await choose({ name: 'large-source.png', mimeType: 'image/png',
      buffer: Buffer.concat([png, Buffer.alloc(33 * 1024 * 1024)]) });
    await page.waitForFunction(() => composer.hasImages && !composer.blocked);
    await page.locator('#text').fill('preserve this caption');
    await page.evaluate(() => {
      composer.send(connection, 'a');
      composer.acknowledge({ type: 'error', image_sequence: connection.imageSequence,
        message: 'Checkpoint with micro_handoff before adding more images.' }, connection);
    });
    assert.equal(await page.evaluate(() => composer.hasImages && !composer.blocked), true);
    assert.equal(await page.locator('#text').inputValue(), 'preserve this caption');
    await page.locator('#previews button').click();
    console.log('PASS 64-MiB source ceiling, >32-MiB source read, and refusal preserves editable draft');
    // A document cannot be previewed as an <img> — pointing one at a PDF renders a
    // broken-image glyph — so the chip states the kind instead. Pages are the
    // server's business; the browser only has to ask for them.
    const html = require('node:fs').readFileSync(path.resolve(__dirname, '../gui/wwwroot/index.html'), 'utf8');
    assert.match(html, /accept="[^"]*application\/pdf/, 'The picker does not offer PDF files.');
    const document = Buffer.from('JVBERi0xLjQKMSAwIG9iago8PC9UeXBlL0NhdGFsb2cvUGFnZXMgMiAwIFI+PgplbmRvYmoKMiAwIG9iago8PC9UeXBlL1BhZ2VzL0tpZHNbNCAwIFJdL0NvdW50IDE+PgplbmRvYmoKMyAwIG9iago8PC9UeXBlL0ZvbnQvU3VidHlwZS9UeXBlMS9CYXNlRm9udC9IZWx2ZXRpY2E+PgplbmRvYmoKNCAwIG9iago8PC9UeXBlL1BhZ2UvUGFyZW50IDIgMCBSL01lZGlhQm94WzAgMCAzMDAgMTAwXS9SZXNvdXJjZXM8PC9Gb250PDwvRjEgMyAwIFI+Pj4+L0NvbnRlbnRzIDUgMCBSPj4KZW5kb2JqCjUgMCBvYmoKPDwvTGVuZ3RoIDQ1Pj5zdHJlYW0KQlQgL0YxIDIwIFRmIDEwIDMwIFRkIChWRVJNSUxJT04gUEFHRTEpIFRqIEVUCmVuZHN0cmVhbQplbmRvYmoKeHJlZgowIDYKMDAwMDAwMDAwMCA2NTUzNSBmIAowMDAwMDAwMDA5IDAwMDAwIG4gCjAwMDAwMDAwNTQgMDAwMDAgbiAKMDAwMDAwMDEwNSAwMDAwMCBuIAowMDAwMDAwMTY4IDAwMDAwIG4gCjAwMDAwMDAyODAgMDAwMDAgbiAKdHJhaWxlcgo8PC9TaXplIDYvUm9vdCAxIDAgUj4+CnN0YXJ0eHJlZgozNzIKJSVFT0YK', 'base64');
    await choose({ name: 'synthetic.pdf', mimeType: 'application/pdf', buffer: document });
    await page.waitForFunction(() => composer.hasImages && !composer.blocked);
    assert.equal(await page.locator('#previews img').count(), 0, 'A PDF was previewed as an image.');
    const chip = page.locator('#previews .attachment-kind');
    assert.equal(await chip.count(), 1);
    assert.equal(await chip.textContent(), 'PDF');
    assert.match(await page.locator('#previews figcaption').first().textContent(), /synthetic\.pdf/);
    await page.evaluate(() => composer.send(connection, 'a'));
    const sent = await page.evaluate(() => frames.at(-1));
    assert.equal(sent.images.length, 1);
    assert.equal(sent.images[0].mime, 'application/pdf', 'The document was sent under a rewritten MIME type.');
    assert.equal(sent.images[0].data, document.toString('base64'));
    await page.evaluate(() => composer.acknowledge({ type: 'image_accepted', sequence: connection.imageSequence,
      agent_id: 'a', images: [] }, connection));
    console.log('PASS document chip, picker acceptance and unrewritten document MIME');
  } finally { await browser.close(); }
})().catch(error => { console.error(error); process.exitCode = 1; });
