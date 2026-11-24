async function loadLatestFindings() {
  const scanId = new URLSearchParams(window.location.search).get("scan_id")
    || window.localStorage.getItem("bastion.latestScanId");
  if (!scanId) {
    return;
  }

  const scanResponse = await fetch(`/api/v1/scans/${encodeURIComponent(scanId)}`);
  if (!scanResponse.ok) {
    throw new Error("Could not load the scan.");
  }
  const scan = await scanResponse.json();
  const targetsResponse = await fetch("/api/v1/targets");
  if (targetsResponse.ok) {
    const targets = await targetsResponse.json();
    const target = targets.find((item) => item.id === scan.target_id);
    const subtitle = document.querySelector("main .head .k");
    if (subtitle) {
      subtitle.textContent = `${target?.value || scan.target_id}, ${scan.profile} profile - ${scan.status}`;
    }
  }

  const response = await fetch(`/api/v1/scans/${encodeURIComponent(scanId)}/findings`);
  if (!response.ok) {
    throw new Error("Could not load scan findings.");
  }
  const findings = await response.json();
  const cards = Array.from(document.querySelectorAll("main > section.card"));
  cards.forEach((card) => card.remove());

  const main = document.querySelector("main");
  const pill = document.querySelector(".head .pill");
  if (pill) {
    pill.textContent = "Scan results";
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

loadLatestFindings().catch(() => {
  const main = document.querySelector("main");
  if (!main) {
    return;
  }
  const message = document.createElement("p");
  message.className = "card";
  message.textContent = "Could not load the latest scan results. Check the server and retry.";
  main.appendChild(message);
});
