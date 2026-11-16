#!/usr/bin/env node
// Run the GUI browser tier in the operator's Firefox and report the whole tier.
// Exit nonzero if anything failed.
//
//   node agent/GuiTests/run-browser-tests.cjs [name ...]
//
// Every suite gets its own GuiTests.dll --crossover-host in its own temporary
// directory. One host shared across the tier cannot work: each suite replays the
// mounted Agent's durable history, so an earlier suite's observations become a
// later suite's unexpected transcript, and a host that dies takes ten unrelated
// suites down with it as refused connections. A private host per suite is the
// isolation the suites themselves document.
//
// Suites in PLAN are not self-contained at all: their real driver also passes a
// fixture host, arguments and an environment. Reproduce that contract here
// instead of reporting a missing argument as a Firefox failure.
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
    // A book-figure capture, not a tier assertion: authorized per run by Bill and
    // it overwrites the approved PNGs in book/figures/advanced-artifacts. Its
    // renderers are covered here by artifact-renderers, compact-artifacts,
    // file-renderers and search-renderers. See that directory's README.
    'advanced-artifact-figures.cjs': 'human-authorized book figure capture',
};
const RUNNER = 'run-browser-tests.cjs';
const HELPER = 'browser-launch.cjs';
const PORT = Number(process.env.GUI_FIXTURE_PORT || 18202);
const URL = `http://127.0.0.1:${PORT}/assistant/`;
const DOTNET = process.env.DOTNET || 'dotnet';
const GUI_TESTS_DLL = path.join(here, 'bin', 'Debug', 'net10.0', 'GuiTests.dll');
const VOICE_TESTS_DLL = path.join(here, '..', 'VoiceLeaseTests', 'bin', 'Debug',
    'net10.0', 'VoiceLeaseTests.dll');
const GUI_BINARY = path.join(here, '..', 'Examples', 'Gui', 'bin', 'Debug', 'net10.0', 'Gui.dll');
// Renderers are verified against a copied distribution, as their drivers do, so
// a repository-only source tree cannot quietly stand in for an installed one.
const SKILLS = path.join(here, '..', 'Examples', 'Gui', 'bin', 'Debug', 'net10.0', 'skills');
const SHARED = path.join(SKILLS, 'shared');
const GOALS_SKILL = path.join(SKILLS, 'goals', 'SKILL.md');
const CROSSOVER_STOP = path.join('.crossover', 'runtime', 'stop');

// host: GuiTests arguments for a suite-owned fixture in place of the crossover
// one; waitFor: the file that marks it ready, or "READY" on its console;
// stopFile: the graceful shutdown marker inside that directory.
const PLAN = {
    'voice-server-page.cjs': {
        // Its driver is the voice-lease fixture, which alone owns the separately
        // funded transcription provider, the model provider and /fixture/status.
        // The crossover host advertises no transcription payer, so the browser
        // would correctly keep its hold button disabled forever.
        host: ['--browser-host'], dll: VOICE_TESTS_DLL, waitFor: 'READY', stopFile: 'stop',
        requires: () => (fs.existsSync(VOICE_TESTS_DLL) ? null : 'built VoiceLeaseTests host'),
    },
    'goal-idle.cjs': {
        host: ['--goal-idle-host'], waitFor: 'fixture.json', stopFile: 'stop',
        requires: () => (fs.existsSync(GOALS_SKILL) ? null : 'built distribution goals skill'),
        hostEnv: { GOAL_SKILL_SOURCE: GOALS_SKILL },
        inputs: (context) => ({ env: { GOAL_FIXTURE: path.join(context.dir, 'fixture.json') } }),
    },
    'working-context.cjs': {
        host: ['--working-context-host'], waitFor: 'READY', stopFile: CROSSOVER_STOP,
        // A working-context save is durable in the Agent's own settings, which is
        // the file whose failure the suite must provoke; the GUI appearance file
        // is no longer that store. The fixture names its one Agent "captain".
        inputs: (context) => ({
            env: {
                GUI_CONTEXT_SETTINGS: path.join(context.dir, '.crossover', 'runtime',
                    'data', 'agents', 'captain', 'settings.json'),
            },
        }),
    },
    'configured-models.cjs': {
        host: (context) => ['--configured-model-fixture', context.dir], waitFor: 'ready.json',
        stopFile: 'stop', inputs: (context) => ({ args: [context.dir] }),
    },
    'advanced-artifact-figures.cjs': {
        host: (context) => ['--artifact-renderer-fixture', context.dir], waitFor: 'url',
        stopFile: 'stop',
        inputs: (context) => ({
            args: [path.join(context.dir, 'figures'), SHARED],
            env: { GUI_FIXTURE_URL: fs.readFileSync(path.join(context.dir, 'url'), 'utf8').trim() },
        }),
    },
    // Both read their origin from the fixture's url file and stop it by writing
    // next to it, so they need the fixture that serves the shipped page from a
    // private directory: the crossover host cannot stand in, because these suites
    // inject frames into the page rather than driving a mounted Agent.
    // Both read their origin from the fixture's url file and stop it by writing
    // next to it, so they need the fixture that serves the shipped page from a
    // private directory: the crossover host cannot stand in, because these suites
    // inject frames into the page rather than driving a mounted Agent.
    'history-window.cjs': {
        host: (context) => ['--artifact-renderer-fixture', context.dir], waitFor: 'url',
        stopFile: 'stop', inputs: (context) => ({ args: [context.dir] }),
    },
    'error-cards-browser.cjs': {
        host: (context) => ['--artifact-renderer-fixture', context.dir], waitFor: 'url',
        stopFile: 'stop', inputs: (context) => ({ args: [context.dir] }),
    },
    'recall-browser.cjs': {
        host: (context) => ['--artifact-renderer-fixture', context.dir], waitFor: 'url',
        stopFile: 'stop', inputs: (context) => ({ args: [context.dir] }),
    },
    'workflow-job-sessions.cjs': {
        host: (context) => ['--artifact-renderer-fixture', context.dir], waitFor: 'url',
        stopFile: 'stop', inputs: (context) => ({ args: [context.dir] }),
    },
    'inline-dice.cjs': {
        requires: () => (fs.existsSync(path.join(SHARED, 'inline-dice.js')) ? null : 'built distribution skills/shared'),
        inputs: () => ({ args: [SHARED] }),
    },
    'settings-dialog.cjs': {
        // The suite writes screenshots into a caller-owned directory; it has no
        // default that exists on this platform.
        inputs: (context) => ({ env: { GUI_ARTIFACT_DIR: path.join(context.dir, 'browser') } }),
    },
    'voice-funding.cjs': {
        // A real-host suite: it launches the shipped GUI binary, so name ours.
        requires: () => (fs.existsSync(GUI_BINARY) ? null : 'built GUI binary'),
        inputs: () => ({ env: { GUI_TEST_BINARY: GUI_BINARY } }),
    },
    'renderer-trials.cjs': {
        requires: () => (fs.existsSync(GUI_TESTS_DLL) ? null : 'built GuiTests host'),
        inputs: () => ({ args: [DOTNET, GUI_TESTS_DLL] }),
    },
};

function suites(filter) {
    return fs.readdirSync(here).filter((name) => name.endsWith('.cjs')
        && name !== RUNNER && name !== HELPER
        && (filter.length === 0 || filter.some((value) => name.includes(value)))).sort();
}

// The crossover fixture dispatches real shell jobs, and this build's launcher is
// the macOS system PTY, so a host that cannot come up is a platform gap rather
// than a Firefox failure. Report the tier that could not be attempted.
function needsFixture(name) {
    const text = fs.readFileSync(path.join(here, name), 'utf8');
    return text.includes('GUI_FIXTURE_URL') || text.includes(String(PORT));
}

function startHost(argumentList, dir, waitFor, extraEnv, assembly) {
    return new Promise((resolve, reject) => {
        const child = spawn(DOTNET, [assembly || GUI_TESTS_DLL, ...argumentList], {
            cwd: dir,
            env: { ...process.env, DOTNET_ROOT: process.env.DOTNET_ROOT || '', ...(extraEnv || {}) },
            stdio: ['ignore', 'pipe', 'pipe'],
        });
        let output = '';
        let settled = false;
        let polling = 0;
        const finish = (error) => {
            if (settled) return;
            settled = true;
            clearTimeout(arming);
            clearInterval(polling);
            if (error) { try { child.kill('SIGKILL'); } catch (ignored) { } reject(error); }
            else resolve(child);
        };
        const arming = setTimeout(
            () => finish(new Error('host did not become ready: ' + argumentList.join(' ') + '\n' + output)), 90000);
        const collect = (chunk) => {
            output += chunk.toString();
            if (waitFor === 'READY' && output.includes('READY')) finish();
        };
        child.stdout.on('data', collect);
        child.stderr.on('data', collect);
        child.on('exit', (code) => finish(new Error(
            'host exited (' + code + '): ' + argumentList.join(' ') + '\n' + output)));
        if (waitFor !== 'READY') {
            polling = setInterval(() => {
                if (fs.existsSync(path.join(dir, waitFor))) finish();
            }, 25);
        }
    });
}

// These hosts serve real sockets, so a tier that leaves one running steals its
// port from the next suite. The graceful exit is the fixture's own stop file; a
// signal is only the backstop, because this build can ignore SIGTERM while
// Kestrel is draining.
function stopHost(child, dir, stopFile) {
    if (!child || child.exitCode !== null) return Promise.resolve();
    if (stopFile && dir) {
        try { fs.writeFileSync(path.join(dir, stopFile), 'stop'); } catch (ignored) { }
    }
    const exited = new Promise((resolve) => child.once('exit', resolve));
    return Promise.race([exited, new Promise((resolve) => setTimeout(resolve, 8000))]).then(() => {
        if (child.exitCode === null) child.kill('SIGKILL');
        return exited;
    });
}

function runSuite(name, env, argv, cwd) {
    return new Promise((resolve) => {
        const started = Date.now();
        const child = spawn(process.execPath, [path.join(here, name), ...argv], {
            cwd, env, stdio: ['ignore', 'pipe', 'pipe'],
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
    const root = fs.mkdtempSync(path.join(os.tmpdir(), 'gui-browser-tier-'));
    const results = [];
    try {
        for (const name of selected) {
            const plan = PLAN[name];
            if (SKIP[name] && filter.length === 0) {
                results.push({ name, skipped: SKIP[name] });
                console.log(`SKIP ${name} — ${SKIP[name]}`);
                continue;
            }
            const missing = plan && plan.requires ? plan.requires() : null;
            if (missing && filter.length === 0) {
                results.push({ name, skipped: 'needs ' + missing });
                console.log(`SKIP ${name} — needs ${missing}`);
                continue;
            }
            // A suite that never dials the fixture serves its own page, so it must
            // not wait on a host it will not use.
            const wantsHost = Boolean(plan && plan.host) || needsFixture(name);
            const ownHost = plan && plan.host ? plan.host : ['--crossover-host'];
            const dir = path.join(root, name.replace(/\.cjs$/, ''));
            fs.mkdirSync(dir, { recursive: true });
            let child = null;
            let outcome;
            try {
                if (wantsHost) {
                    const argumentsForHost = typeof ownHost === 'function' ? ownHost({ dir }) : ownHost;
                    child = await startHost(argumentsForHost, dir,
                        (plan && plan.waitFor) || 'READY', plan && plan.hostEnv, plan && plan.dll);
                }
                const given = plan && plan.inputs ? plan.inputs({ dir }) : {};
                const env = { ...process.env, GUI_FIXTURE_URL: URL, GUI_URL: URL, ...(given.env || {}) };
                outcome = await runSuite(name, env, given.args || [], dir);
            } catch (error) {
                const text = error.message || String(error);
                if (!plan && needsFixture(name) && filter.length === 0) {
                    results.push({ name, skipped: 'fixture host unavailable' });
                    console.log(`SKIP ${name} — fixture host unavailable`);
                    console.log('    | ' + text.split('\n').slice(0, 3).join('\n    | '));
                    continue;
                }
                outcome = { name, code: 1, seconds: 0, tail: 'fixture setup: ' + text };
            } finally {
                if (child) await stopHost(child, dir, plan && plan.stopFile);
            }
            // The next suite must not inherit this port while it is still held.
            await new Promise((resolve) => setTimeout(resolve, 250));
            results.push(outcome);
            console.log(`${outcome.code === 0 ? 'PASS' : 'FAIL'} ${name} (${outcome.seconds.toFixed(1)}s)`);
            if (outcome.code !== 0) {
                console.log(outcome.tail.split('\n').slice(-30).map((line) => '    | ' + line).join('\n'));
            }
        }
    } finally {
        if (!process.env.BROWSER_TEST_KEEP) {
            try { fs.rmSync(root, { recursive: true, force: true }); } catch (ignored) { }
        } else {
            console.log('kept fixture workspaces under ' + root);
        }
    }
    const passed = results.filter((value) => value.code === 0).length;
    const skipped = results.filter((value) => value.skipped).length;
    const failed = results.length - passed - skipped;
    console.log(`\nbrowser tier: ${passed} passed, ${failed} failed, ${skipped} skipped of ${results.length}`);
    if (failed > 0) { process.exitCode = 1; }
})();
