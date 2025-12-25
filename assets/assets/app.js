let results = [];

function getLines(text) {
  return text
    .split("\n")
    .map(x => x.trim())
    .filter(Boolean);
}

function sortByQuality(list) {
  // Option B: "accepted" should appear above "uncertain"
  // We'll sort accepted first, then uncertain. (rejected is filtered out anyway)
  const order = { accepted: 1, uncertain: 2, rejected: 3 };
  return list.slice().sort((a, b) => (order[a.status] || 9) - (order[b.status] || 9));
}

function renderTable() {
  const tbody = document.querySelector("#table tbody");
  tbody.innerHTML = "";
  results.forEach(r => {
    const tr = document.createElement("tr");
    tr.innerHTML = `<td>${r.email}</td><td>${r.status}</td><td>${r.reason}</td>`;
    tbody.appendChild(tr);
  });
}

function renderTxt() {
  const mode = document.getElementById("mode").value;

  let filtered = results.filter(r => {
    if (r.status === "accepted") return true;
    if (mode === "accepted_uncertain" && r.status === "uncertain") return true;
    return false;
  });

  filtered = sortByQuality(filtered);

  // No headers, no empty lines, emails only
  document.getElementById("output").value = filtered.map(r => r.email).join("\n");
}

async function run() {
  const input = getLines(document.getElementById("input").value);
  const includeDuplicates = document.getElementById("dupes").checked;

  const res = await fetch("/validate", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ emails: input, includeDuplicates })
  });

  results = await res.json();
  renderTable();
  renderTxt();
}

async function copyTxt() {
  const txt = document.getElementById("output").value;
  await navigator.clipboard.writeText(txt);
  // simple feedback (no fancy UI yet)
  const btn = document.getElementById("btnCopy");
  const old = btn.textContent;
  btn.textContent = "Copied ✓";
  setTimeout(() => (btn.textContent = old), 900);
}

function downloadTxt() {
  const txt = document.getElementById("output").value;
  const blob = new Blob([txt], { type: "text/plain;charset=utf-8" });
  const a = document.createElement("a");
  a.href = URL.createObjectURL(blob);

  const mode = document.getElementById("mode").value;
  a.download = (mode === "accepted") ? "accepted.txt" : "accepted_and_uncertain.txt";

  a.click();
  URL.revokeObjectURL(a.href);
}

document.getElementById("btnRun").addEventListener("click", run);
document.getElementById("mode").addEventListener("change", renderTxt);
document.getElementById("btnCopy").addEventListener("click", copyTxt);
document.getElementById("btnDownload").addEventListener("click", downloadTxt);
