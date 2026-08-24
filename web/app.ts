import { family, familyStart, type FamilyPerson } from "./family-data";

const root = document.documentElement;
const themeButton = document.querySelector<HTMLButtonElement>("[data-theme-toggle]");
const storedTheme = localStorage.getItem("theme");
const prefersDark = matchMedia("(prefers-color-scheme: dark)").matches;
const themeLabel = document.querySelector<HTMLElement>("[data-theme-label]");
const themeIcon = document.querySelector<HTMLElement>("[data-theme-icon]");

if (storedTheme === "dark" || (!storedTheme && prefersDark)) root.classList.add("dark");

const updateThemeButton = () => {
  const dark = root.classList.contains("dark");
  if (themeButton) themeButton.ariaLabel = `Switch to ${dark ? "light" : "dark"} theme`;
  if (themeLabel) themeLabel.textContent = dark ? "Light" : "Dark";
  if (themeIcon) themeIcon.textContent = dark ? "☀" : "☾";
};

updateThemeButton();

themeButton?.addEventListener("click", () => {
  root.classList.toggle("dark");
  localStorage.setItem("theme", root.classList.contains("dark") ? "dark" : "light");
  updateThemeButton();
});

const search = document.querySelector<HTMLInputElement>("[data-post-search]");
const cards = [...document.querySelectorAll<HTMLElement>("[data-post-card]")];
const empty = document.querySelector<HTMLElement>("[data-search-empty]");

search?.addEventListener("input", () => {
  const query = search.value.trim().toLowerCase();
  let visible = 0;
  cards.forEach((card) => {
    const match = !query || (card.dataset.search ?? "").includes(query);
    card.hidden = !match;
    if (match) visible++;
  });
  document.querySelectorAll<HTMLElement>("[data-year-group]").forEach((group) => {
    group.hidden = !group.querySelector("[data-post-card]:not([hidden])");
  });
  empty?.classList.toggle("hidden", visible !== 0);
});

const familyExplorer = document.querySelector<HTMLElement>("[data-family-explorer]");

if (familyExplorer) {
  const focus = familyExplorer.querySelector<HTMLElement>("[data-family-focus]")!;
  const children = familyExplorer.querySelector<HTMLElement>("[data-family-children]")!;
  const parents = familyExplorer.querySelector<HTMLElement>("[data-family-parents]")!;
  const connector = familyExplorer.querySelector<HTMLElement>(".family-connector")!;
  const parentConnector = familyExplorer.querySelector<HTMLElement>(".family-parent-connector")!;

  const orderedFacts = (person: FamilyPerson) => [...person.facts].sort((left, right) => {
    const rank = (fact: string) => fact.startsWith("Born") ? 0 : fact.startsWith("Married") ? 1 : fact.startsWith("Died") ? 2 : 3;
    return rank(left) - rank(right);
  });

  const facts = (person: FamilyPerson) => person.facts.length
    ? `<ul class="family-facts">${orderedFacts(person).map((fact) => `<li>${fact}</li>`).join("")}</ul>`
    : "";

  const notes = (person: FamilyPerson) => person.notes?.length
    ? `<div class="family-notes">${person.notes.map((note) => `<p>${note}</p>`).join("")}${person.source ? `<a href="${person.source.url}">${person.source.label} ↗</a>` : ""}</div>`
    : "";

  const familySummary = (node: FamilyPerson) => {
    const parts: string[] = [];
    if (node.spouses.length) parts.push(`${node.spouses.length} ${node.spouses.length === 1 ? "spouse" : "spouses"}`);
    if (node.children.length) parts.push(`${node.children.length} ${node.children.length === 1 ? "child" : "children"}`);
    return parts.join(" / ");
  };

  const personCard = (id: string, relation: string, extraClass = "") => {
    const person = family[id];
    const summary = familySummary(person);
    return `<button type="button" class="family-person ${extraClass}" data-family-id="${id}"><span>${relation}</span><strong>${person.name}</strong>${person.facts.length ? `<ul class="family-card-facts">${orderedFacts(person).map((fact) => `<li>${fact}</li>`).join("")}</ul>` : ""}${summary ? `<em>${summary}</em>` : ""}</button>`;
  };

  const renderFamily = (id: string, updateHash = true) => {
    const node = family[id];
    if (!node) return;
    const spouseMarkup = node.spouses.map((spouseID) => `<div class="family-marriage"><span aria-hidden="true"></span></div>${personCard(spouseID, "Spouse", "family-member family-member--spouse")}`).join("");
    focus.innerHTML = `<div class="family-couple"><div class="family-member family-member--primary"><span>Selected</span><h2>${node.name}</h2>${facts(node)}</div>${spouseMarkup}</div>${notes(node)}${node.children.length ? `<p class="family-count">${node.children.length} ${node.children.length === 1 ? "child" : "children"}</p>` : ""}`;
    parents.innerHTML = node.parents.map((parentID) => personCard(parentID, "Parent", "family-person--parent")).join("");
    children.innerHTML = node.children.map((childID) => personCard(childID, "Child")).join("");
    connector.hidden = node.children.length === 0;
    parentConnector.hidden = node.parents.length === 0;
    if (updateHash) {
      history.replaceState(null, "", `#${id}`);
      focus.scrollIntoView({ behavior: "smooth", block: "center" });
    }
  };

  familyExplorer.querySelector(".family-stage")?.addEventListener("click", (event) => {
    const target = (event.target as HTMLElement).closest<HTMLButtonElement>("[data-family-id]");
    if (!target?.dataset.familyId) return;
    renderFamily(target.dataset.familyId);
  });

  const initialID = location.hash.slice(1);
  renderFamily(family[initialID] ? initialID : familyStart, false);
}

const articleImages = [...document.querySelectorAll<HTMLImageElement>(".prose-custom img")];

if (articleImages.length) {
  const dialog = document.createElement("dialog");
  dialog.className = "image-lightbox";
  dialog.setAttribute("aria-label", "Expanded image viewer");
  dialog.innerHTML = `
    <div class="image-lightbox__bar">
      <span data-lightbox-count aria-live="polite"></span>
      <div class="flex items-center gap-2">
        <button class="image-lightbox__button" type="button" data-lightbox-previous aria-label="Previous image">←</button>
        <button class="image-lightbox__button" type="button" data-lightbox-zoom>Actual size</button>
        <button class="image-lightbox__button" type="button" data-lightbox-next aria-label="Next image">→</button>
        <button class="image-lightbox__button" type="button" data-lightbox-close>Close</button>
      </div>
    </div>
    <div class="image-lightbox__stage" data-lightbox-stage>
      <img class="image-lightbox__image" data-lightbox-image alt="">
    </div>`;
  document.body.append(dialog);

  const expanded = dialog.querySelector<HTMLImageElement>("[data-lightbox-image]")!;
  const count = dialog.querySelector<HTMLElement>("[data-lightbox-count]")!;
  const zoom = dialog.querySelector<HTMLButtonElement>("[data-lightbox-zoom]")!;
  let active = 0;

  const show = (index: number) => {
    active = (index + articleImages.length) % articleImages.length;
    const source = articleImages[active];
    expanded.src = source.currentSrc || source.src;
    expanded.alt = source.alt;
    expanded.dataset.zoomed = "false";
    zoom.textContent = "Actual size";
    count.textContent = `${active + 1} / ${articleImages.length}${source.alt ? ` · ${source.alt}` : ""}`;
  };

  articleImages.forEach((image, index) => {
    image.tabIndex = 0;
    image.setAttribute("role", "button");
    image.setAttribute("aria-label", `${image.alt || "Article image"}. Open larger view.`);
    const open = (event: Event) => {
      event.preventDefault();
      show(index);
      dialog.showModal();
    };
    image.addEventListener("click", open);
    image.addEventListener("keydown", (event) => {
      if (event.key === "Enter" || event.key === " ") open(event);
    });
  });

  dialog.querySelector("[data-lightbox-close]")?.addEventListener("click", () => dialog.close());
  dialog.querySelector("[data-lightbox-previous]")?.addEventListener("click", () => show(active - 1));
  dialog.querySelector("[data-lightbox-next]")?.addEventListener("click", () => show(active + 1));
  dialog.querySelector("[data-lightbox-stage]")?.addEventListener("click", (event) => {
    if (event.target !== expanded) dialog.close();
  });
  expanded.addEventListener("click", () => {
    const zoomed = expanded.dataset.zoomed === "true";
    expanded.dataset.zoomed = String(!zoomed);
    zoom.textContent = zoomed ? "Actual size" : "Fit screen";
  });
  zoom.addEventListener("click", () => expanded.click());
  dialog.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      event.preventDefault();
      dialog.close();
    }
    if (event.key === "ArrowLeft") show(active - 1);
    if (event.key === "ArrowRight") show(active + 1);
  });
}
