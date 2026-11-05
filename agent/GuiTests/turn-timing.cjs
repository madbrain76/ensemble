// The per-turn split, using the shipped formatter and the shipped state machine.
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');
const source = fs.readFileSync(path.join(__dirname, '../gui/wwwroot/app.js'), 'utf8');
const start = source.indexOf('const turn = {');
const end = source.indexOf('function handle(frame) {');
assert(start >= 0 && end > start, 'the timing block must ship in the page');
let body = source.slice(start, end);
const now = () => Date.parse('2026-11-03T12:00:00.000Z') + clock.offset;
let clock = { offset: 0 };
const document = {
    querySelectorAll: () => [],
    querySelector: () => null,
    createElement: () => ({ className: '', textContent: '', remove() {} }),
    body: { appendChild() {} },
};
const cards = new Map();
const context = { Date: { now: () => now() }, document, cards, console };
vm.createContext(context);
vm.runInContext(body + `
this.turn = turn; this.beginTurn = beginTurn; this.closeTurn = closeTurn;
this.recordTurnTiming = recordTurnTiming; this.turnTimingText = turnTimingText;`, context);
const { turn, beginTurn, recordTurnTiming, turnTimingText } = context;
const check = (ok, label) => { assert.ok(ok, label); console.log('PASS ' + label); };

check(turnTimingText({ started: 0, ended: 0 }) === '', 'an unfinished turn shows no summary');

beginTurn();                                     // t = 0, prompt submitted
recordTurnTiming({ type: 'agent_state', state: 'inflight' });
clock.offset += 10900; recordTurnTiming({ type: 'tool_dispatched' });
recordTurnTiming({ type: 'agent_state', state: 'toolspending' });
clock.offset += 1000; recordTurnTiming({ type: 'agent_state', state: 'idle' });
const line = turnTimingText(turn);
check(line.startsWith('Total 11.9s'), 'total runs from submission to idle: ' + line);
check(line.includes('Model 10.9s'), 'model wait counts only inflight spans: ' + line);
check(line.includes('Tools (1) 1.0s'), 'tool wait counts dispatches and their blocking span: ' + line);
check(!line.includes('Other'), 'no always-zero bucket is offered: ' + line);

beginTurn();
clock.offset += 500; recordTurnTiming({ type: 'part_delta', part_id: 1, chunk: 'x' });
clock.offset += 200; recordTurnTiming({ type: 'agent_state', state: 'inflight' });
clock.offset += 700; recordTurnTiming({ type: 'agent_state', state: 'idle' });
const quiet = turnTimingText(turn);
check(quiet.includes('Tools (0) 0.0s'), 'a tool-free turn still shows the split: ' + quiet);
check(quiet.startsWith('Total 1.4s'), 'a short turn keeps one decimal: ' + quiet);
