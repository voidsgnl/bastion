async function loadLatestFindings() {
  const scanId = new URLSearchParams(window.location.search).get("scan_id")
    || window.localStorage.getItem("bastion.latestScanId");
  if (!scanId) {
    return;
  }

  const response = await fetch(`/api/v1/scans/${encodeURIComponent(scanId)}/findings`);
  if (!response.ok) {
    return;
  }
  const findings = await response.json();
  const cards = Array.from(document.querySelectorAll("main > section.card"));
  cards.forEach((card) => card.remove());

  const main = document.querySelector("main");
  const pill = document.querySelector(".head .pill");
  if (pill) {
    pill.textContent = "Live fixture";
  }

  findings.forEach((finding) => {
    const card = document.createElement("section");
    card.className = "card";
    card.innerHTML = `
      <div class="head">
        <h2></h2>
        <span class="pill l"></span>
      </div>
      <p class="mt"></p>
      <div class="acts">
        <a class="btn" href="finding.html">View fix</a>
        <a class="btn ghost" href="finding.html">Show evidence</a>
      </div>`;
    card.querySelector("h2").textContent = finding.title;
    card.querySelector(".pill").textContent =
      finding.confidence.charAt(0).toUpperCase() + finding.confidence.slice(1);
    card.querySelector(".mt").textContent =
      `Risk score ${finding.risk.score}: ${finding.risk.explanation}`;
    main.appendChild(card);
  });
}

loadLatestFindings();
