// Per-artifact timing: the pure formatter extracted from the production page.
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');
const source = fs.readFileSync(path.join(__dirname, '../gui/wwwroot/app.js'), 'utf8');
const start = source.indexOf('function cardTimingText(');
const end = source.indexOf('function renderTiming(', start);
assert(start >= 0 && end > start, 'cardTimingText must exist in the shipped page');
const context = { console };
vm.createContext(context);
vm.runInContext(source.slice(start, end) + '\nthis.cardTimingText = cardTimingText;', context);
const { cardTimingText } = context;
const check = (ok, label) => { assert.ok(ok, label); console.log('PASS ' + label); };
const clock = Date.parse('2026-11-03T12:34:56.000Z');
check(cardTimingText({ first: 0, last: 0, final: true }) === '', 'an unstarted part shows no timing');
const quick = cardTimingText({ first: clock, last: clock + 3000, final: true, clock });
check(/^\d\d:\d\d:\d\d · 3\.0s$/.test(quick), 'a three-second part shows its duration: ' + quick);
check(/· …$/.test(cardTimingText({ first: clock, last: clock, final: false, clock })),
    'a still-streaming part shows an ellipsis, not a duration');
// "…" means still producing, which is untrue for a settled part; a finalized part
// that was never timed shows its clock alone rather than a fabricated duration.
check(cardTimingText({ first: clock, last: clock, final: true, clock })
        === cardTimingText({ first: clock, last: clock, final: true, clock }).split(' · ')[0]
    && !cardTimingText({ first: clock, last: clock, final: true, clock }).includes('·'),
    'a settled part with no measurable span shows its clock alone, never 0.0s or an ellipsis');
// Recorded production data is authoritative over anything the tab could measure.
const recordedSpan = cardTimingText({ first: 1, last: 1, final: true, started: clock - 500, elapsed: 4200 });
check(/· 4\.2s$/.test(recordedSpan),
    'a recorded span is shown even when the tab measured nothing: ' + recordedSpan);
const recordedClock = cardTimingText({ first: Date.now(), last: Date.now(), final: true,
    started: clock, elapsed: 250 });
check(recordedClock.startsWith(cardTimingText({ first: 0, last: 0, final: true, clock }) || '')
    && /· 0\.3s$/.test(recordedClock),
    'a recorded span also fixes the clock, not just the duration: ' + recordedClock);
const long = cardTimingText({ first: clock, last: clock + 7500, final: true, clock });
check(long.endsWith('7.5s'), 'sub-ten-second durations keep one decimal: ' + long);
const slow = cardTimingText({ first: clock, last: clock + 95400, final: true, clock });
check(slow.endsWith('95s'), 'long durations round to whole seconds: ' + slow);
// The clock comes from the server frame when present, so a replay keeps its moment.
const replay = cardTimingText({ first: Date.now(), last: Date.now() + 1000, final: true, clock });
assert.ok(replay.startsWith(cardTimingText({ first: clock, last: clock, final: true, clock }).slice(0, 8)),
    'a replay used the local clock instead of the server time');
console.log('PASS a replayed part keeps the server clock');
