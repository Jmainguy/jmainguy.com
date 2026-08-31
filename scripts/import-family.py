"""Import the actual boxes and connector lines in Mainguy.ca's FTM charts.

Run with Python + pdfplumber. Original PDFs and extracted graph are retained.
"""
import json
import re
from pathlib import Path
import subprocess
from urllib.parse import quote
import pdfplumber

ROOT = Path(__file__).resolve().parents[1]
DEST = ROOT / 'research/family'
SOURCES = {
    'maingi': 'gy_genealogy/g maingi 7 gen.pdf',
    'pierre': 'gy_genealogy/pierre gi 5 gen.pdf',
    'nicholas': 'gay_intro/nicholas 1730 3 gen.pdf',
    'william-uk': 'gay_United_Kingdom/william 7 gen UK.pdf',
    'peter': 'gay_United_Kingdom/peter gy 4 gen.pdf',
    'anthony-uk': 'gay_United_Kingdom/anthony gay 4 gen uk2.pdf',
    'leslie-australia': 'gay_Australia/leslie 5 gen aust.pdf',
    'john-naples': 'gay_Naples/john gay 4gen naples.pdf',
    'thomas-netherlands': 'gay_Netherlands/thomas 7 gen NL.pdf',
    'marcus-nz': 'gay_New_Zealand/marcus gay 7gen new zealand.pdf',
    'dan-western-canada': 'guy_w_canada/dan wishart 6gen.pdf',
}

def near(a, b):
    return abs(a-b) < 0.9

def parse(path):
    people = []
    with pdfplumber.open(path) as pdf:
        for page in pdf.pages:
            boxes = [c for c in page.curves if c['width'] > 5 and c['height'] > 5]
            for box in boxes:
                text = page.crop((box['x0'], box['top'], box['x1'], box['bottom'])).extract_text(x_tolerance=1) or ''
                lines = text.splitlines()
                split = next((i for i, l in enumerate(lines) if re.match(r'[bmd]:', l)), len(lines))
                name = ' '.join(lines[:split])
                if not name:
                    raise ValueError(f'Empty box: {path} {box}')
                facts = []
                for line in lines[split:]:
                    if re.match(r'[bmd]:', line):
                        facts.append({'b': 'Born', 'm': 'Married', 'd': 'Died'}[line[0]] + line[2:])
                    else:
                        facts[-1] += ' ' + line
                box['person'] = len(people)
                people.append(dict(name=name, facts=facts, parents=[], spouses=[], children=[], page=page.page_number,
                                   box=[round(box[k], 2) for k in ('x0', 'top', 'x1', 'bottom')]))
            def add(a, relation, b):
                if b not in people[a][relation]:
                    people[a][relation].append(b)
            def touches(line, box):
                result = set()
                for x, y in line['pts']:
                    if box['x0']-.9 <= x <= box['x1']+.9:
                        if near(y, box['top']): result.add('top')
                        if near(y, box['bottom']) or near(y, box['bottom']+3): result.add('bottom')
                    if box['top'] <= y <= box['bottom']:
                        if near(x, box['x0']): result.add('left')
                        if near(x, box['x1']) or near(x, box['x1']+3): result.add('right')
                return result
            # Double-line marriages can be horizontal OR vertical in FTM's compact layout.
            marriage = set()
            for i, line in enumerate(page.lines):
                touching = [b['person'] for b in boxes if touches(line, b)]
                if len(touching) == 2:
                    for j, other in enumerate(page.lines[:i]):
                        if [b['person'] for b in boxes if touches(other, b)] == touching:
                            marriage.update((i, j))
                            a, b = touching
                            add(a, 'spouses', b)
                            add(b, 'spouses', a)
            lines = [l for i, l in enumerate(page.lines) if i not in marriage]
            groups = list(range(len(lines)))
            def root(i):
                while groups[i] != i:
                    i = groups[i]
                return i
            for i, a in enumerate(lines):
                for j, b in enumerate(lines[:i]):
                    # Connected buses include their T-junctions, not merely endpoints.
                    if max(a['x0'], b['x0']) <= min(a['x1'], b['x1'])+.2 and max(a['top'], b['top']) <= min(a['bottom'], b['bottom'])+.2:
                        groups[root(i)] = root(j)
            for group in {root(i) for i in range(len(lines))}:
                network = [l for i, l in enumerate(lines) if root(i) == group]
                ancestors, children = [], []
                for box in boxes:
                    sides = set().union(*(touches(l, box) for l in network))
                    if 'bottom' in sides: ancestors.append(box['person'])
                    elif sides & {'top', 'left'}: children.append(box['person'])
                if children and len(ancestors) != 1:
                    raise ValueError(f'Ambiguous connector in {path}: {ancestors} -> {children}')
                for a in ancestors:
                    for c in children:
                        add(c, 'parents', a)
                        add(a, 'children', c)
            # Multiple marriages are drawn as a row/column of adjacent boxes,
            # not a sequence of marriages between the adjacent spouses.
            visited = set()
            for box in boxes:
                start = box['person']
                if start in visited: continue
                component, todo = set(), [start]
                while todo:
                    a = todo.pop()
                    if a in component: continue
                    component.add(a)
                    todo.extend(people[a]['spouses'])
                visited.update(component)
                if len(component) > 2:
                    centers = [a for a in component if people[a]['parents']]
                    if not centers:
                        title = (page.extract_text() or '').splitlines()[0].removeprefix('Descendants of ').lower()
                        centers = [a for a in component if people[a]['name'].lower() == title]
                    if len(centers) != 1:
                        raise ValueError(f'Ambiguous multiple marriage: {path} {component}')
                    center = centers[0]
                    for a in component:
                        people[a]['spouses'] = sorted(component-{a}) if a == center else [center]
            # The line leaves one member of the marriage. Include the other only
            # where that partner is unambiguous (never all of several spouses).
            for box in boxes:
                child = people[box['person']]
                for a in child['parents'][:]:
                    partners = people[a]['spouses']
                    if len(partners) == 1:
                        add(box['person'], 'parents', partners[0])
                        add(partners[0], 'children', box['person'])
    return people

def main():
    DEST.mkdir(parents=True, exist_ok=True)
    records = {}
    for key, relative in SOURCES.items():
        url = 'http://mainguy.ca/mainguy_website_new/' + quote(relative)
        path = DEST / (key + '.pdf')
        if not path.exists():
            subprocess.run(['curl', '-fsSL', '--max-time', '45', url, '-o', str(path)], check=True)
        people = parse(path)
        records[key] = dict(url=url, people=people)
        print(key, len(people), 'people;', sum(len(p['parents']) for p in people), 'parent links')
    (DEST / 'charts.json').write_text(json.dumps(records, indent=2) + '\n')

if __name__ == '__main__':
    main()
