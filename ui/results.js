let activeScanId = new URLSearchParams(window.location.search).get("scan_id")
  || window.localStorage.getItem("bastion.latestScanId");

function showMessage(message) {
  const element = document.querySelector("#action-message");
  if (element) {
    element.textContent = message;
  }
}

async function requestJSON(path, options) {
  const response = await fetch(path, options);
  const body = await response.json();
  if (!response.ok) {
    throw new Error(body.error || "Bastion request failed");
  }
  return body;
}

function makeElement(tag, className, text) {
  const element = document.createElement(tag);
  if (className) {
    element.className = className;
  }
  if (text) {
    element.textContent = text;
  }
  return element;
}

function renderFindings(findings) {
  const main = document.querySelector("main");
  const cards = Array.from(main.querySelectorAll(":scope > section.card"));
  cards.forEach((card) => card.remove());

  if (findings.length === 0) {
    const empty = makeElement("section", "card", "No findings were reported by this scan.");
    main.appendChild(empty);
    return;
  }

  for (const finding of findings) {
    const card = makeElement("section", "card");
    const header = makeElement("div", "head");
    const title = makeElement("h2", "", finding.title);
    const confidence = makeElement("span", `pill ${confidenceClass(finding.confidence)}`,
      `${finding.severity} / ${finding.confidence}`);
    header.append(title, confidence);
    card.appendChild(header);

    card.appendChild(makeElement("p", "mt",
      `Risk score ${finding.risk.score}: ${finding.risk.explanation}`));
    card.appendChild(makeElement("p", "k",
      `${finding.source_engine || "unknown engine"} - ${finding.asset_id || "unknown asset"} - ID ${finding.id} - ${finding.status} - first seen ${formatDate(finding.first_seen)} - last seen ${formatDate(finding.last_seen)}`));

    const evidence = makeElement("pre", "mt2");
    evidence.textContent = (finding.evidence || [])
      .map((item) => `${item.type}: ${item.summary}`)
      .join("\n") || "No evidence attached.";
    card.appendChild(evidence);
    main.appendChild(card);
  }
}

function confidenceClass(confidence) {
  if (confidence === "confirmed") {
    return "c";
  }
  if (confidence === "likely") {
    return "l";
  }
  return "u";
}

function formatDate(value) {
  if (!value) {
    return "unknown";
  }
  const date = new Date(value);
  return Number.isNaN(date.valueOf()) ? "unknown" : date.toLocaleString();
}

async function loadScanResults() {
  if (!activeScanId) {
    return;
  }
  const scan = await requestJSON(`/api/v1/scans/${encodeURIComponent(activeScanId)}`);
  const targets = await requestJSON("/api/v1/targets");
  const target = targets.find((item) => item.id === scan.target_id);
  const subtitle = document.querySelector("main .head .k");
  if (subtitle) {
    subtitle.textContent = `${target?.value || scan.target_id}, ${scan.profile} profile - ${scan.status}`;
  }

  const findings = await requestJSON(`/api/v1/scans/${encodeURIComponent(activeScanId)}/findings`);
  const pill = document.querySelector(".head .pill");
  if (pill) {
    pill.textContent = "Scan results";
  }
  renderFindings(findings);
  const analystNote = document.querySelector("aside .msg");
  if (analystNote) {
    analystNote.textContent = "Review scan evidence and verify observations before taking action.";
  }
}

function renderComparison(comparison) {
  const main = document.querySelector("main");
  main.querySelector("#comparison-summary")?.remove();

  const section = makeElement("section", "card");
  section.id = "comparison-summary";
  section.appendChild(makeElement("h2", "", "Scan comparison"));
  section.appendChild(makeElement("p", "k mt",
    `Baseline ${comparison.baseline_scan_id} → current ${comparison.current_scan_id}`));

  for (const [label, findings] of [
    ["New", comparison.new],
    ["Fixed", comparison.fixed],
    ["Unchanged", comparison.unchanged],
  ]) {
    const group = makeElement("div", "row");
    group.append(
      makeElement("b", "", `${label} (${findings.length})`),
      makeElement("span", "grow", findings.map((finding) => finding.title).join(", ") || "None"),
    );
    section.appendChild(group);
  }
  main.appendChild(section);
}

document.querySelector("#compare-scans")?.addEventListener("click", async () => {
  const baselineScanId = window.localStorage.getItem("bastion.previousScanId");
  if (!activeScanId || !baselineScanId) {
    showMessage("Run at least two scans to compare results.");
    return;
  }
  try {
    const comparison = await requestJSON(
      `/api/v1/scans/${encodeURIComponent(baselineScanId)}/compare/${encodeURIComponent(activeScanId)}`,
    );
    renderComparison(comparison);
    showMessage("");
  } catch (error) {
    showMessage(error.message);
  }
});

document.querySelector("#retest-scan")?.addEventListener("click", async (event) => {
  if (!activeScanId) {
    showMessage("Choose a completed scan before retesting.");
    return;
  }
  if (!window.confirm("Run a new safe-profile scan of the same authorized target?")) {
    return;
  }

  const button = event.currentTarget;
  button.disabled = true;
  showMessage("Retest running...");
  try {
    const scan = await requestJSON(
      `/api/v1/scans/${encodeURIComponent(activeScanId)}/retest`,
      { method: "POST" },
    );
    window.localStorage.setItem("bastion.previousScanId", activeScanId);
    window.localStorage.setItem("bastion.latestScanId", scan.id);
    window.location.assign(`results.html?scan_id=${encodeURIComponent(scan.id)}`);
  } catch (error) {
    showMessage(error.message);
    button.disabled = false;
  }
});

loadScanResults().catch((error) => {
  showMessage(error.message || "Could not load scan results.");
});
