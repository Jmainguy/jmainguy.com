import assert from 'node:assert/strict';
import fs from 'node:fs';
import { test } from 'node:test';
import { build } from 'esbuild';

const result = await build({entryPoints:['web/family-data.ts'],bundle:true,write:false,platform:'node',format:'esm'});
const { family, familyStart } = await import(`data:text/javascript;base64,${Buffer.from(result.outputFiles[0].text).toString('base64')}`);
const charts = JSON.parse(fs.readFileSync('research/family/charts.json','utf8'));
const identities = JSON.parse(fs.readFileSync('research/family/identity-map.json','utf8'));
const layoutBuild = await build({entryPoints:['web/family-parent-layout.ts'],bundle:true,write:false,platform:'node',format:'esm'});
const { orderedParents } = await import(`data:text/javascript;base64,${Buffer.from(layoutBuild.outputFiles[0].text).toString('base64')}`);

test('parent display puts the father left without changing source order', () => {
  const william = family['william-anstruther-maingy'];
  const original = [...william.parents];
  assert.deepEqual(orderedParents(william), ['pierre-32', 'pierre-33']);
  assert.deepEqual(william.parents, original);
  assert.deepEqual(orderedParents(family['pierre-146']), ['pierre-45', 'pierre-46']);
  assert.deepEqual(orderedParents({parents:['euphemia-mary-maingy','charles-pope']}), ['charles-pope','euphemia-mary-maingy']);
  assert.deepEqual(orderedParents({parents:['marcus-nz-62','marcus-nz-63']}), ['marcus-nz-63','marcus-nz-62']);
  assert.deepEqual(orderedParents({parents:['unknown-a','unknown-b']}), ['unknown-a','unknown-b']);
  assert.deepEqual(orderedParents({parents:['maingi-0']}), ['maingi-0']);
  assert.deepEqual(orderedParents({parents:[]}), []);
});

test('all 11 charts and 778 boxes are accounted for', () => {
  assert.equal(Object.keys(charts).length,11);
  assert.equal(Object.values(charts).reduce((sum,c)=>sum+c.people.length,0),778);
  for(const [chart,data] of Object.entries(charts)) data.people.forEach((p,i)=>{
    const person=family[identities[`${chart}:${i}`]];
    assert.ok(person);
    assert.ok(person.sources.some(s=>s.url===data.url));
  });
});

test('all relationships exist, are reciprocal, have no self links or ancestry cycles', () => {
  const done = new Set(), visiting = new Set();
  function visit(id) {
    assert.ok(!visiting.has(id),`cycle at ${id}`);
    if(done.has(id))return;
    visiting.add(id); family[id].parents.forEach(visit); visiting.delete(id); done.add(id);
  }
  for(const [id,person] of Object.entries(family)) {
    assert.ok(person.parents.length<=2,`too many parents for ${id}`);
    for(const [relation,reverse] of [['parents','children'],['children','parents'],['spouses','spouses']]) {
      assert.equal(new Set(person[relation]).size,person[relation].length);
      person[relation].forEach(other=>{
        assert.notEqual(id,other);
        assert.ok(family[other]?.[reverse].includes(id),`${id} ${relation} ${other}`);
      });
    }
    visit(id);
  }
});

test('every person is reachable from the existing Holden branch', () => {
  const seen = new Set(), todo = [familyStart];
  while(todo.length) {
    const id=todo.pop(); if(seen.has(id))continue; seen.add(id);
    todo.push(...family[id].parents,...family[id].spouses,...family[id].children);
  }
  assert.equal(seen.size,Object.keys(family).length);
  assert.equal(seen.size,786);
});

test('William ascends to early Guillaume; James connects to the Western Canada chart', () => {
  const lineage=['william-anstruther-maingy','pierre-32','pierre-11','pierre-2','maingi-29','maingi-27','maingi-20','maingi-11','maingi-4','maingi-1','maingi-0'];
  lineage.slice(1).forEach((parent,i)=>assert.ok(family[lineage[i]].parents.includes(parent)));
  assert.equal(identities['dan-western-canada:0'],'pierre-146');
  assert.ok(family['pierre-146'].parents.includes('pierre-45'));
  assert.equal(family['pierre-45'].children.length,13);
  assert.equal(family['pierre-146'].children.length,5);
});

test('compact vertical siblings and multiple marriages are interpreted correctly', () => {
  const p=charts.maingi.people;
  assert.deepEqual(p[20].spouses,[18,19,21,22]); // Pierre's four spouses, not spouse-to-spouse links.
  assert.deepEqual(new Set(p[24].parents),new Set([18,20]));
  assert.deepEqual(new Set(p[27].parents),new Set([21,20]));
  assert.ok(family['pierre-107'].parents.includes('pierre-32')); // Edwin is William's sibling.
  assert.ok(!family['pierre-107'].parents.includes('william-anstruther-maingy'));
  assert.equal(identities['pierre:25'],identities['peter:2']); // Undated Anne de la Combe.
  assert.ok(family['pierre-8'].notes.some(n=>n.startsWith('Source discrepancy:')));
});
