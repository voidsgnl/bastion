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
    const savedTargetId = window.localStorage.getItem("bastion.demoTargetId");
    const targets = await api("/api/v1/targets");
    let target = targets.find((item) =>
      item.id === savedTargetId && item.type === "cidr" && item.value === "127.0.0.1/32",
    );
    if (!target) {
      target = await api("/api/v1/targets", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: "Demo lab",
          type: "cidr",
          value: "127.0.0.1/32",
        }),
      });
      window.localStorage.setItem("bastion.demoTargetId", target.id);
    }
    await api(`/api/v1/targets/${encodeURIComponent(target.id)}/attestation`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ statement: authorizationStatement }),
    });
    const previousScanId = window.localStorage.getItem("bastion.latestScanId");
    const scan = await api("/api/v1/scans", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ target_id: target.id, profile: "safe" }),
    });
    if (previousScanId) {
      window.localStorage.setItem("bastion.previousScanId", previousScanId);
    }
    window.localStorage.setItem("bastion.latestScanId", scan.id);
    window.location.assign(`results.html?scan_id=${encodeURIComponent(scan.id)}`);
  } catch (error) {
    newScanButton.textContent = "Scan failed";
    window.alert(error.message);
  }
});
