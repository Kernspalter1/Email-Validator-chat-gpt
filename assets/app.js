document.getElementById("btnRun").addEventListener("click", run);

async function run() {
  const text = document.getElementById("input").value;

  const res = await fetch("/validate", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ text })
  });

  const entries = await res.json();
  renderTable(entries);
}

function renderTable(entries) {
  const body = document.getElementById("tableBody");
  body.innerHTML = "";

  entries.forEach(e => {
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${e.email}</td>
      <td>${e.isDuplicate ? "duplicate" : "unique"}</td>
    `;
    body.appendChild(tr);
  });
}
