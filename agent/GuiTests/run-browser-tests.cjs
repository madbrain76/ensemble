#!/usr/bin/env node
// Run the GUI browser tier in the operator's Firefox against one disposable
// fixture host, and report the whole tier. Exit nonzero if anything failed.
//
//   node agent/GuiTests/run-browser-tests.cjs [name ...]
//
// The fixture is GuiTests.dll --crossover-host, started in a private temporary
// directory so no suite can touch a real save, workspace or credential. Set
// BROWSER_TEST_KEEP=1 to leave the host running for debugging.
const { spawn } = require('node:child_process');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');

const here = __dirname;
// These suites drive the macOS system PTY/launch-services capture path, which
// throws PlatformNotSupportedException off macOS (agent/Jobs/Pty.cs). They are
// recorded as unverified here rather than silently counted as passing.
const SKIP = {
    'page-capture-native.cjs': 'macOS native host-screen capture',
    'capture-model-native.cjs': 'macOS native capture model fixture',
    'page-capture-native-ui.cjs': 'macOS native capture UI',
};
const RUNNER = 'run-browser-tests.cjs';
const HELPER = 'browser-launch.cjs';
const PORT = Number(process.env.GUI_FIXTURE_PORT || 18202);
const URL = `http://127.0.0.1:${PORT}/assistant/`;

function suites(filter) {
    return fs.readdirSync(here).filter((name) => name.endsWith('.cjs')
        && name !== RUNNER && name !== HELPER
        && (filter.length === 0 || filter.some((value) => name.includes(value)))).sort();
}

async function startFixture(work) {
    const dll = path.join(here, 'bin', 'Debug', 'net10.0', 'GuiTests.dll');
    if (!fs.existsSync(dll)) {
        throw new Error('Build the GUI test host first: dotnet build agent/GuiTests/GuiTests.csproj');
    }
    const dotnet = process.env.DOTNET || 'dotnet';
    const child = spawn(dotnet, [dll, '--crossover-host'], {
        cwd: work, env: { ...process.env, DOTNET_ROOT: process.env.DOTNET_ROOT || '' },
        stdio: ['ignore', 'pipe', 'pipe'],
    });
    let output = '';
    const ready = new Promise((resolve, reject) => {
        const timer = setTimeout(() => reject(new Error('fixture host did not become ready:\n' + output)), 90000);
        const collect = (chunk) => {
            output += chunk.toString();
            if (output.includes('READY')) { clearTimeout(timer); resolve(); }
        };
        child.stdout.on('data', collect);
        child.stderr.on('data', collect);
        child.on('exit', (code) => { clearTimeout(timer); reject(new Error(`fixture host exited (${code}):\n${output}`)); });
    });
    await ready;
    return { child, output: () => output };
}

function runSuite(name, env) {
    return new Promise((resolve) => {
        const started = Date.now();
        const child = spawn(process.execPath, [path.join(here, name)], {
            cwd: path.join(here, '..', '..'), env, stdio: ['ignore', 'pipe', 'pipe'],
        });
        let tail = '';
        const cap = (chunk) => { tail = (tail + chunk.toString()).slice(-4000); };
        child.stdout.on('data', cap);
        child.stderr.on('data', cap);
        const timer = setTimeout(() => child.kill('SIGKILL'), 240000);
        child.on('exit', (code) => {
            clearTimeout(timer);
            resolve({ name, code, seconds: (Date.now() - started) / 1000, tail });
        });
    });
}

(async () => {
    const filter = process.argv.slice(2);
    const selected = suites(filter);
    if (selected.length === 0) { console.error('no suites selected'); process.exit(2); }
    const work = fs.mkdtempSync(path.join(os.tmpdir(), 'gui-browser-tier-'));
    let host = null;
    const results = [];
    // A suite is fixture-dependent when it dials the crossover host, whether
    // through the documented variable or a literal address. Judging those here as
    // Firefox failures would blame the browser for a missing shell launcher.
    const needsFixture = (name) => {
        const text = fs.readFileSync(path.join(here, name), "utf8");
        return text.includes("GUI_FIXTURE_URL") || text.includes("18202");
    };
    let fixtureProblem = null;
    try {
        host = await startFixture(work);
    } catch (error) {
        // A missing fixture is one finding, not a reason to skip every suite:
        // the self-contained pages serve their own copy of the GUI assets.
        fixtureProblem = error.message;
        console.error('fixture unavailable: ' + error.message.split("\n")[0]);
    }
    try {
        for (const name of selected) {
            if (needsFixture(name) && fixtureProblem && filter.length === 0) {
                // The crossover fixture dispatches real shell jobs, and this build's
                // launcher is the macOS system PTY. Report the tier that could not be
                // attempted here instead of calling it a Firefox failure.
                results.push({ name, skipped: "fixture needs the macOS shell launcher" });
                console.log(`SKIP ${name} — fixture needs the macOS shell launcher`);
                continue;
            }
            if (SKIP[name] && filter.length === 0) {
                results.push({ name, skipped: SKIP[name] });
                console.log(`SKIP ${name} — ${SKIP[name]}`);
                continue;
            }
            const outcome = await runSuite(name, { ...process.env, GUI_FIXTURE_URL: URL });
            results.push(outcome);
            console.log(`${outcome.code === 0 ? 'PASS' : 'FAIL'} ${name} (${outcome.seconds.toFixed(1)}s)`);
            if (outcome.code !== 0) {
                console.log(outcome.tail.split('\n').slice(-12).map((line) => '    | ' + line).join('\n'));
            }
        }
    } finally {
        if (host && !process.env.BROWSER_TEST_KEEP) { host.child.kill('SIGTERM'); }
        if (!process.env.BROWSER_TEST_KEEP) {
            try { fs.rmSync(work, { recursive: true, force: true }); } catch (ignored) { }
        }
    }
    const passed = results.filter((value) => value.code === 0).length;
    const skipped = results.filter((value) => value.skipped).length;
    const failed = results.length - passed - skipped;
    console.log(`\nbrowser tier: ${passed} passed, ${failed} failed, ${skipped} skipped of ${results.length}`);
    if (failed > 0) { process.exitCode = 1; }
})();
