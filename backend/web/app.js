const health = document.querySelector('#health');
const unit = document.querySelector('#unit');
const rows = document.querySelector('#rows');
const alertRows = document.querySelector('#alert-rows');
const readingRows = document.querySelector('#reading-rows');
const message = document.querySelector('#message');

async function load() {
  const h = await (await fetch('/healthz')).json();
  health.textContent = `${h.service}: ${h.status}`;
  health.className = 'ok';
  const data = await (await fetch('/api/capture-units')).json();
  unit.replaceChildren(...data.items.map(x => new Option(x.id, x.id)));
  rows.replaceChildren(...data.items.map(x => { const row = document.createElement('tr'); row.innerHTML = `<td>${x.id}</td><td>${x.facility}</td><td>${x.capture_rate_pct}%</td><td>${x.pressure_kpa} kPa</td><td>${x.solvent_level_pct}%</td><td>${x.status}</td>`; return row; }));
  await loadAlerts();
  await loadReadings();
}

async function loadAlerts() {
  const data = await (await fetch('/api/alerts?page=1&page_size=10')).json();
  alertRows.replaceChildren(...data.items.map(x => { const row = document.createElement('tr'); row.innerHTML = `<td>${x.id}</td><td>${x.subject}</td><td>${x.priority}</td><td>${x.status}</td><td>${x.owner}</td><td>${x.updated_at}</td>`; return row; }));
}

async function loadReadings() {
  const units = await (await fetch('/api/capture-units')).json();
  const fragment = document.createDocumentFragment();
  for (const item of units.items) {
    const summary = await (await fetch(`/api/readings/summary?unit=${item.id}&window_minutes=60`)).json();
    const row = document.createElement('tr');
    row.innerHTML = `<td>${item.id}</td><td>${summary.stats.count}</td><td>${summary.stats.avg_capture_rate_pct.toFixed(1)}%</td><td>${summary.stats.avg_pressure_kpa.toFixed(1)} kPa</td><td>${summary.stats.avg_solvent_level_pct.toFixed(1)}%</td>`;
    fragment.appendChild(row);
  }
  readingRows.replaceChildren(fragment);
}

document.querySelector('#status-form').addEventListener('submit', async event => {
  event.preventDefault();
  const response = await fetch('/api/capture-units/status', {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({id: unit.value, status: document.querySelector('#status').value})});
  const data = await response.json();
  message.textContent = response.ok ? `已更新 ${data.id}` : data.error;
  await load();
});

document.querySelector('#refresh-alerts').addEventListener('click', loadAlerts);
document.querySelector('#refresh-readings').addEventListener('click', loadReadings);
load().catch(error => { health.textContent = error.message; });
