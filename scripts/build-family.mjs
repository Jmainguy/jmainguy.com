// Build the additive, source-linked graph without overwriting the hand-maintained family.
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const text = fs.readFileSync(`${root}/web/family-data.ts`, 'utf8');
const base = JSON.parse(text.split('export const family: Record<string, FamilyPerson> = ')[1].split(' as Record<string, FamilyPerson>;')[0]);
const charts = JSON.parse(fs.readFileSync(`${root}/research/family/charts.json`, 'utf8'));
const normalize = name => name.replace(/\[\d+\]\s*/g, '').toLowerCase().replace(/[^a-z0-9]/g, '');
const birth = p => p.facts.find(f => f.startsWith('Born'))?.match(/\b\d{4}\b/)?.[0];
const nodes = Object.entries(base).map(([id, p]) => ({id, ...p, chart:'existing'}));
const raw = [];
for (const [chart, data] of Object.entries(charts)) {
  const offset = raw.length;
  for (const [i, p] of data.people.entries()) raw.push({...p, chart, index:i, url:data.url,
    parents:p.parents.map(n=>n+offset), spouses:p.spouses.map(n=>n+offset), children:p.children.map(n=>n+offset)});
}
const canonical = [];
const compatible = (a,b) => !birth(a) || !birth(b) || birth(a) === birth(b);
for (const p of raw) {
  const matches = nodes.filter(n => normalize(n.name) === normalize(p.name) && birth(p) && birth(n) === birth(p) && n.chart !== p.chart);
  if(matches.length === 1) canonical.push(matches[0].id);
  else {
    const id = `${p.chart}-${p.index}`;
    nodes.push({id, ...p});
    canonical.push(id);
  }
}
// Undated spouses are merged only when a dated/shared relative establishes identity.
// Also handles the source's explicitly numbered repeated [1] Susanna box.
const redirects = new Map();
const resolve = id => redirects.has(id) ? resolve(redirects.get(id)) : id;
for(let pass=0;pass<4;pass++) for(const [i,p] of raw.entries()) {
  const current = resolve(canonical[i]);
  const relatives = [...p.parents, ...p.spouses].map(j=>resolve(canonical[j]));
  const candidates = raw.map((q,j)=>({q,j})).filter(({q,j}) => resolve(canonical[j]) !== current &&
    (q.chart !== p.chart || /^\[\d+\]/.test(p.name)) && normalize(q.name) === normalize(p.name) && compatible(p,q) &&
    [...q.parents,...q.spouses].some(k=>relatives.includes(resolve(canonical[k]))));
  const ids = [...new Set(candidates.map(({j})=>resolve(canonical[j])))];
  if(ids.length && new Set(candidates.map(({q})=>birth(q)).filter(Boolean)).size <= 1) {
    const ordered = [current,...ids].sort((a,b)=>nodes.findIndex(n=>n.id===a)-nodes.findIndex(n=>n.id===b));
    for(const id of ordered.slice(1)) redirects.set(id,ordered[0]);
  }
}
const output = {};
const unique = values => [...new Set(values)];
for(const [i,p] of raw.entries()) {
  const id = resolve(canonical[i]);
  const existing = base[id];
  const person = output[id] ??= existing ? structuredClone(existing) : {
    name:p.name.replace(/\[\d+\]\s*/g,'').replace(/\b[A-Z][A-Z]+\b/g,s=>s[0]+s.slice(1).toLowerCase()),
    facts:[], parents:[], spouses:[], children:[]
  };
  person.facts = unique([...person.facts, ...p.facts]);
  for(const relation of ['parents','spouses','children']) person[relation] = unique([...person[relation], ...p[relation].map(j=>resolve(canonical[j]))]);
  person.sources ??= [];
  if(!person.sources.some(s=>s.url===p.url)) person.sources.push({label:`Mainguy.ca — ${p.chart} chart, page ${p.page}`, url:p.url});
}
const all = {...base,...output};
// Make every link reciprocal, including links into the hand-maintained graph.
for(const [id,p] of Object.entries(all)) for(const [relation,reverse] of [['parents','children'],['children','parents'],['spouses','spouses']]) {
  for(const target of p[relation]) {
    if(!all[target] || id === target) throw Error(`Invalid ${relation}: ${id} -> ${target}`);
    if(!all[target][reverse].includes(id)) {
      output[target] ??= structuredClone(all[target]);
      all[target] = output[target];
      all[target][reverse].push(id);
    }
  }
}
const visiting = new Set(), done = new Set();
function visit(id) {
  if(visiting.has(id)) throw Error(`Ancestry cycle: ${id}`);
  if(done.has(id)) return;
  visiting.add(id); all[id].parents.forEach(visit); visiting.delete(id); done.add(id);
}
Object.keys(all).forEach(visit);
for(const [id,p] of Object.entries(output)) {
  if(p.parents.length>2) throw Error(`Too many parents for ${id}: ${p.parents}`);
  for(const parent of p.parents) {
    const death=all[parent].facts.find(f=>f.startsWith('Died'))?.match(/\b\d{4}\b/)?.[0];
    if(birth(p) && death && Number(birth(p))>Number(death)+1) {
      p.notes ??= [];
      const note=`Source discrepancy: this chart places ${p.name} (born ${birth(p)}) under ${all[parent].name}, whose recorded death is ${death}. The chart relationship is retained, not independently verified.`;
      if(!p.notes.includes(note)) p.notes.push(note);
    }
  }
}
fs.writeFileSync(`${root}/web/family-imports.json`, JSON.stringify(output,null,2)+'\n');
fs.writeFileSync(`${root}/research/family/identity-map.json`,JSON.stringify(Object.fromEntries(raw.map((p,i)=>[`${p.chart}:${p.index}`,resolve(canonical[i])])),null,2)+'\n');
console.log(`${raw.length} chart boxes; ${Object.keys(all).length} total people; ${Object.keys(output).filter(id=>!base[id]).length} added.`);
