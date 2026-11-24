const newScanButton = document.querySelector(".head .btn");

async function api(path, options) {
  const response = await fetch(path, options);
  const body = await response.json();
  if (!response.ok) {
    throw new Error(body.error || "Bastion request failed");
  }
  return body;
}

newScanButton?.addEventListener("click", async (event) => {
  event.preventDefault();
  const authorizationStatement =
    "I confirm this local lab scan is limited to 127.0.0.1/32.";
  if (!window.confirm(`${authorizationStatement}\n\nContinue and record this attestation?`)) {
    return;
  }
  newScanButton.textContent = "Starting...";

  try {
    const target = await api("/api/v1/targets", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: "Demo lab",
        type: "cidr",
        value: "127.0.0.1/32",
      }),
    });
    await api(`/api/v1/targets/${encodeURIComponent(target.id)}/attestation`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ statement: authorizationStatement }),
    });
    const scan = await api("/api/v1/scans", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ target_id: target.id, profile: "safe" }),
    });
    const findings = await api(`/api/v1/findings?scan_id=${encodeURIComponent(scan.id)}`);

    window.localStorage.setItem("bastion.latestScanId", scan.id);

    const recentSection = Array.from(document.querySelectorAll("section.card")).find(
      (section) => section.querySelector("h2")?.textContent === "Recent scans",
    );
    const recent = recentSection?.querySelector("h2");
    if (recent) {
      recent.textContent = `${findings.length} finding from the latest scan`;
    }

    const analyst = document.querySelector("aside .msg");
    if (analyst) {
      analyst.textContent = findings[0]
        ? `${findings[0].title}. Risk score: ${findings[0].risk.score}.`
        : "The scan completed without findings.";
    }
    const samplePill = document.querySelector(".head .pill");
    if (samplePill) {
      samplePill.textContent = "Sample dashboard - scan complete";
    }
    newScanButton.textContent = "Scan complete";
  } catch (error) {
    newScanButton.textContent = "Scan failed";
    window.alert(error.message);
  }
});
