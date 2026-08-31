import type { FamilyPerson } from './family-data';

// Explicit presentation preferences for the current tree, not inferred sex
// metadata. Keep this separate from the source transcription and never sort by
// surname, chart position, or a runtime guess about a person's first name.
// Unspecified couples retain their recorded order until reviewed.
const leftParents = new Set(`
pierre-32 william-anstruther-maingy william-mcleod-mainguy charles-pope
helenus-gilbert-mainguy lefeuvre-anstruther-mainguy philip-anstruther-wynder-mainguy
philip-neville-mainguy william-anstruther-mainguy charles-pope-2
herbert-layton-mainguy william-francis-mainguy robert-edgeworth-mainguy
charles-pope-3 harry-layton-bray julian-wimberly robert-samuel-anstruther-mainguy
richard-christian william-hasenauer holden-leonard-mainguy john-herbert-layton-mainguy
william-neville-kilbourn-mainguy donald-neville-mainguy henk-westra george-r-kirby-mainguy
michael-campbell david-holden-mainguy robert-mainguy russell-w-lammer
herbert-wayne-mainguy eric-crosier russell-kirk peter-william-mainguy
robert-thomas-mainguy david-gregory-mainguy brett-david-massey jonathan-mainguy
maingi-4 maingi-11 maingi-20 maingi-27 maingi-29
pierre-2 pierre-8 pierre-11 pierre-16 pierre-18 pierre-21 pierre-24 pierre-27
pierre-37 pierre-41 pierre-45 pierre-60 pierre-64 pierre-77 pierre-88 pierre-89
pierre-96 pierre-100 pierre-160
william-uk-18 william-uk-22 william-uk-24 william-uk-34 william-uk-36
william-uk-40 william-uk-44 pierre-176 peter-37 anthony-uk-5 anthony-uk-9 peter-29
leslie-australia-5 leslie-australia-7 leslie-australia-9 leslie-australia-11
leslie-australia-13 leslie-australia-22 leslie-australia-24 leslie-australia-27
leslie-australia-30 leslie-australia-32 leslie-australia-40 leslie-australia-44
leslie-australia-68 pierre-184 pierre-156
thomas-netherlands-10 thomas-netherlands-12 thomas-netherlands-15
thomas-netherlands-38 thomas-netherlands-43 thomas-netherlands-47
thomas-netherlands-49 thomas-netherlands-51 thomas-netherlands-53
thomas-netherlands-55 thomas-netherlands-57 thomas-netherlands-65
thomas-netherlands-67 thomas-netherlands-69 thomas-netherlands-71
pierre-120 marcus-nz-4 marcus-nz-11 marcus-nz-13 marcus-nz-21 marcus-nz-25
marcus-nz-32 marcus-nz-37 marcus-nz-39 marcus-nz-46 marcus-nz-48
marcus-nz-54 marcus-nz-60 marcus-nz-63 pierre-146
dan-western-canada-5 dan-western-canada-8 dan-western-canada-11
dan-western-canada-13 dan-western-canada-15 dan-western-canada-18
dan-western-canada-21 dan-western-canada-23 dan-western-canada-26
dan-western-canada-27 dan-western-canada-30 dan-western-canada-32
dan-western-canada-33 dan-western-canada-40 dan-western-canada-47 dan-western-canada-50
`.trim().split(/\s+/));

export function orderedParents(person: Pick<FamilyPerson, 'parents'>): string[] {
  const parents = [...person.parents];
  if (parents.length === 2 && leftParents.has(parents[1]) && !leftParents.has(parents[0])) {
    parents.reverse();
  }
  return parents;
}
