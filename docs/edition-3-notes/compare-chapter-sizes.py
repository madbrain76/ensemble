#!/usr/bin/env python3
"""Compare frozen Edition 1 against current Edition 2, with explicit chapter lineage."""
import hashlib
import json
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / 'docs/edition-3-notes'
# October 11 merge map composed with October 14 memory/fleet reorder.
MAP = {0:[0], 1:[1], 2:[2], 3:[3], 4:[4,5], 5:[6], 6:[7], 7:[8],
       8:[9], 9:[10], 10:[11], 11:[12], 12:[13], 13:[14], 14:[15],
       15:[18], 16:[19], 17:[20], 18:[21,22], 19:[23], 20:[24,25],
       21:[26], 22:[16], 23:[17], 24:[], 25:[], 26:[27]}
assert sorted(n for ns in MAP.values() for n in ns) == list(range(28))

def measure(path):
    raw=path.read_bytes(); text=raw.decode('utf-8'); outside=[]; fenced=[]; fence=None
    for line in text.splitlines():
        m=re.match(r'^ {0,3}(`{3,}|~{3,})(.*)$',line)
        if fence is None and m:
            fence=(m[1][0],len(m[1]));fenced.append(line)
        elif fence is not None:
            fenced.append(line)
            if m and m[1][0]==fence[0] and len(m[1])>=fence[1] and not m[2].strip():fence=None
        else:outside.append(line)
    assert fence is None, f'Unclosed fence: {path}'
    total=len(text.split());nonfenced=len('\n'.join(outside).split());in_fence=len('\n'.join(fenced).split())
    assert total == nonfenced+in_fence
    return {'path':str(path.relative_to(ROOT)), 'sha256':hashlib.sha256(raw).hexdigest(),
            'bytes':len(raw),'words':total,'outside_fences_words':nonfenced,
            'fenced_words':in_fence,'title':next(l[2:] for l in text.splitlines() if l.startswith('# '))}

rows=[]
for e2, e1s in MAP.items():
    new=measure(ROOT/f'book/chapter-{e2:02d}.md')
    old=[measure(ROOT/f'edition-1/book/chapter-{n:02d}.md') for n in e1s]
    row={'edition2_chapter':e2,'edition1_chapters':e1s,'edition1_files':old,'edition2_file':new}
    for metric in ('words','bytes','outside_fences_words','fenced_words'):
        before=sum(f[metric] for f in old);after=new[metric]
        row[metric]={'edition1':before,'edition2':after,'delta':after-before,
                     'growth_percent':100*(after-before)/before if before else None}
    rows.append(row)
matched=[r for r in rows if r['edition1_chapters']]
totals={}
for name,subset in [('matched_chapters',matched),('all_numbered_chapters',rows),('matched_excluding_intro',[r for r in matched if r['edition2_chapter']!=0])]:
    totals[name]={}
    for metric in ('words','bytes','outside_fences_words','fenced_words'):
        before=sum(r[metric]['edition1'] for r in subset);after=sum(r[metric]['edition2'] for r in subset)
        totals[name][metric]={'edition1':before,'edition2':after,'delta':after-before,'growth_percent':100*(after-before)/before}
result={'repository_head':subprocess.check_output(['git','rev-parse','HEAD'],cwd=ROOT,text=True).strip(),
        'mapping_source':'docs/edition-2-chapter-map.md',
        'method':'UTF-8 byte count and Python whitespace-delimited words. Outside-fence words exclude fenced blocks and delimiters, but include headings, tables, inline code and link destinations. These are not model tokens or rendered prose counts.',
        'scope':'Numbered manuscript chapters only, including Chapter 0. Excludes preface, epilogue, notes, outlines, review records, generated assembled books and Edition 3. Merged chapters aggregate their Edition 1 predecessors; distributed architectural lessons are not attributed sentence by sentence.',
        'totals':totals,'chapters':rows}
(OUT/'edition-1-vs-2-chapter-sizes.json').write_text(json.dumps(result,indent=2)+'\n')
lines=['# Edition 1 → Edition 2 chapter sizes','',result['method'],'',result['scope'],'',
       f"Source HEAD: `{result['repository_head']}`. File hashes are in the [data](edition-1-vs-2-chapter-sizes.json).",'',
       'Mapping: [edition map](../edition-2-chapter-map.md), October 11 mergers plus October 14 reorder. The three combined predecessor groups are 4+5→4, 21+22→18 and 24+25→20. Edition 2 Chapters 24 and 25 are new.', '',
       '| E1 → E2 | Topic | E1 words | E2 words | Added | Growth | Outside-fence growth |',
       '|---|---|---:|---:|---:|---:|---:|']
for r in rows:
    old='+'.join(map(str,r['edition1_chapters'])) or 'new';w=r['words'];p=r['outside_fences_words'];topic=re.sub(r'^Chapter \d+: ','',r['edition2_file']['title'])
    growth=f"{w['growth_percent']:+.1f}%" if w['growth_percent'] is not None else '—'
    pg=f"{p['growth_percent']:+.1f}%" if p['growth_percent'] is not None else '—'
    lines.append(f"| {old} → {r['edition2_chapter']} | {topic} | {w['edition1']:,} | {w['edition2']:,} | {w['delta']:+,} | {growth} | {pg} |")
lines+=['','## Totals','']
for name,t in totals.items():
    w=t['words'];p=t['outside_fences_words'];lines.append(f"- **{name.replace('_',' ')}:** {w['edition1']:,} → {w['edition2']:,} words ({w['growth_percent']:+.1f}%). Outside fences: {p['edition1']:,} → {p['edition2']:,} ({p['growth_percent']:+.1f}%).")
lines+=['','## Interpretation limits','',
        'Growth is a review priority, not a finding of padding. High percentage growth from a short original may reflect missing teaching; large absolute growth can matter even at a modest percentage. Compare added sections and their provenance before proposing cuts. Preserve meaning and avoid new ambiguity.',
        '', 'The prior Chapter 11 pruning pilot examined Edition 3 in isolation. This report measures Edition 2 against Edition 1; those are different questions.','']
(OUT/'edition-1-vs-2-chapter-sizes.md').write_text('\n'.join(lines))
print('\n'.join(lines))
print('\nLARGEST ABSOLUTE INCREASES:')
for r in sorted(matched,key=lambda r:r['words']['delta'],reverse=True)[:10]:
    print(r['edition2_chapter'],r['words'],r['outside_fences_words'])
