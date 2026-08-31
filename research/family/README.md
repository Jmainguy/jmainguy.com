# Mainguy.ca chart import

Retrieved 2026-08-31 from http://mainguy.ca/mainguy_website_new/.
The 11 original PDFs are retained here. `charts.json` records each source URL,
person's text, page, bounding box, and chart-local relationships. `identity-map.json`
maps all 778 source boxes to explorer IDs. The import adds 610 people to the
176-person hand-maintained tree: 786 connected people altogether.

## Reproduce

From the repository root, with Python + pdfplumber and the project's Node dependencies:

```sh
python3 -m pip install -r scripts/requirements.txt
python3 scripts/import-family.py
node scripts/build-family.mjs
node --test scripts/family.test.mjs
npm run build
go test ./...
```

Existing PDFs are reused; the importer downloads only missing files. The generated
`web/family-imports.json` is applied additively to the original `web/family-data.ts`
records. Existing IDs, user-supplied facts, and biographical notes are preserved.

## Sources

| Local PDF | Chart / branch |
| --- | --- |
| maingi.pdf | Guillaume Maingi, seven generations |
| pierre.pdf | Pierre Maingi, five generations; William Anstruther and James connections |
| nicholas.pdf | Nicholas Maingy, three generations |
| william-uk.pdf | William Maingay, United Kingdom |
| peter.pdf | Peter Maingy, United Kingdom / Australia |
| anthony-uk.pdf | Anthony Tyndall Maingay, United Kingdom |
| leslie-australia.pdf | Frederick Gother Maingay, Australia (source filename: Leslie) |
| john-naples.pdf | John Maingy, Naples |
| thomas-netherlands.pdf | Thomas Maingay, Netherlands |
| marcus-nz.pdf | Marcus Maingay, New Zealand |
| dan-western-canada.pdf | Daniel Wishart Mainguy, Western Canada |

The Peter PDF linked from Australia is byte-identical to the UK copy. The Pierre
PDFs linked from the Mainguy introduction and New Zealand pages are byte-identical
to the genealogy-page copy. These duplicate links were checked by SHA-256.

## Interpretation and limitations

- Relationships are transcriptions of these secondary-source charts, not verified
  civil or parish records. Names and dates retain source spelling, apart from
  display capitalization and removal of chart repetition markers such as `[1]`.
- Box geometry and actual line networks determine links; text reading order does
  not. In compact charts, vertically stacked people can be siblings, and double
  marriage lines can be vertical. Multi-spouse chains are normalized to the
  descendant's marriages, rather than marriages between adjacent spouses.
- A second parent is added only when the marriage partner is unambiguous. Empty
  parent lists mean no link was established, not that the person had no parents.
- Cross-chart identity matching uses matching names and birth years, or matching
  names and already-shared parents/spouses with compatible birth years. Names
  alone, surname similarity, or regional proximity do not establish a merge.
- Date conflicts are preserved and explicitly noted on affected people. For
  example, the chart gives Pierre Maingi's death as 1719 yet lists children born
  later, including Nicholas in 1730. Do not treat the resulting path as proof of
  genealogy until that discrepancy is resolved from better sources.
- The chart's James Maingy (1804–1883) and Charlotte Beckworth are parents of
  Daniel Wishart Mainguy (1842–1906), the root of the Western Canada chart.
- All 786 nodes are reachable from Holden. Tests cover reciprocity, cycles,
  source coverage, the older ancestry path, Western Canada, vertical siblings,
  and multiple marriages. Browser checks cover branch navigation and history.

No deployment or publication is performed by these scripts.
