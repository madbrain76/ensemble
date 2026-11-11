// One browser launch for every GUI regression: the operator's Firefox.
//
// These suites used to require Playwright out of the frozen edition-2 tree and
// hard-launch macOS Google Chrome, so on Linux, and for anyone who does not run
// Chrome, the whole browser tier silently could not run. Bill reviews this GUI in
// Firefox, so verification happens there. No download is performed: PLAYWRIGHT_MODULE
// or the vendored module beside the GUI must already exist, and FIREFOX_BINARY may
// name an installed browser.
const path = require('node:path');

const candidates = [
    process.env.PLAYWRIGHT_MODULE,
    path.join(__dirname, '..', 'gui', 'node_modules', 'playwright'),
    'playwright',
];
let playwright = null;
const errors = [];
for (const candidate of candidates) {
    if (!candidate) continue;
    try { playwright = require(candidate); break; } catch (error) { errors.push(candidate + ': ' + error.message); }
}
if (!playwright) {
    throw new Error('No Playwright installation found. Install once with "npm install playwright" '
        + 'and "npx playwright install firefox", or set PLAYWRIGHT_MODULE. Tried:\n  '
        + errors.join('\n  '));
}
const { firefox } = playwright;

// Flags the Chrome build needed. Firefox ignores unknown switches rather than
// refusing to start, and the suites that depend on fake audio inject the web
// APIs themselves, so none of them belong in a Firefox launch.
const CHROME_FLAGS = /^(--(auto-accept-this-tab-capture|autoplay-policy|mute-audio|use-fake-|enable-features|disable-features|allow-insecure-localhost|unsafely-treat)|--headless)/;

function firefoxArgs(args) {
    return (args || []).filter((flag) => !CHROME_FLAGS.test(flag));
}

// The capture-handle identity API is Chromium-only, and agent/gui/wwwroot/
// page-capture.js will not offer the camera unless it can see all three of
// getDisplayMedia, setCaptureHandleConfig and MediaStreamTrack.prototype
// .getCaptureHandle. The capture suites fabricate the whole capture surface, so
// Firefox only needs the two probed members to exist: each of them installs its
// own getDisplayMedia and its own per-track getCaptureHandle afterwards.
async function installCaptureHandleShim(page) {
    await page.addInitScript(() => {
        if (!navigator.mediaDevices.setCaptureHandleConfig) {
            Object.defineProperty(navigator.mediaDevices, 'setCaptureHandleConfig',
                { configurable: true, value: () => {} });
        }
        if (!MediaStreamTrack.prototype.getCaptureHandle) {
            Object.defineProperty(MediaStreamTrack.prototype, 'getCaptureHandle',
                { configurable: true, writable: true, value: () => undefined });
        }
    });
}

async function launchBrowser(options = {}) {
    const { channel, executablePath, args, ...rest } = options;
    return firefox.launch({
        ...rest,
        headless: options.headless !== false,
        executablePath: process.env.FIREFOX_BINARY || undefined,
        args: firefoxArgs(args),
    });
}

// Persistent profile variant, for suites that need a stable profile directory.
async function launchPersistentProfile(userDir, options = {}) {
    const { channel, executablePath, args, ...rest } = options;
    return firefox.launchPersistentContext(userDir, {
        ...rest,
        headless: options.headless !== false,
        executablePath: process.env.FIREFOX_BINARY || undefined,
        args: firefoxArgs(args),
    });
}

module.exports = { launchBrowser, launchPersistentProfile, firefox, installCaptureHandleShim };
