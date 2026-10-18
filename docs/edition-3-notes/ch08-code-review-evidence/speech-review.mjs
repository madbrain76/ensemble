import assert from 'node:assert/strict';
import {pathToFileURL} from 'node:url';
const {SpeechQueue,nativeSpeechRate}=await import(pathToFileURL(process.argv[2]).href);
let rate=1.2345;const spoken=[];let cancels=0;class U{constructor(text){this.text=text}}
const q=new SpeechQueue(()=>{}, {speak:u=>spoken.push(u),cancel:()=>cancels++},U,()=>rate);
q.enqueue('fraction one');q.enqueue('fraction two');assert.equal(spoken[0].rate,1.2345);
rate=1.8765;assert.equal(spoken[0].rate,1.2345,'active utterance must not be re-rated');spoken[0].onend({type:'end'});assert.equal(spoken[1].rate,1.8765);
rate=4.25;q.enqueue('extended approximate');spoken[1].onend({type:'end'});assert.equal(spoken[2].rate,2.75);
const stale=spoken[2];q.cancel();assert(!q.busy);assert.equal(cancels,1);rate=.7654;q.enqueue('fraction below one');stale.onend({type:'end'});assert.equal(spoken[3].rate,.7654);assert(q.busy);spoken[3].onerror({type:'error',error:'controlled'});assert(!q.busy);
assert.equal(nativeSpeechRate(2),2);assert(nativeSpeechRate(3.125)>nativeSpeechRate(3));
console.log(JSON.stringify({passed:true,rates:spoken.map(u=>u.rate),scope:'Actual production SpeechQueue with synthetic synthesis; NOT browser/native speech proof'}));
