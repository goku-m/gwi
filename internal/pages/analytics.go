package pages

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func Analytics() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Analytics</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Manrope:wght@400;500;600;700;800&display=swap" rel="stylesheet">
  <style>
    :root {
      --bg: #eef8ee;
      --surface: #ffffff;
      --surface-soft: #f8fbfa;
      --text: #102127;
      --muted: #56707d;
      --accent: #166534;
      --accent-soft: #eef8ee;
      --border: #d7e0d8;
      --shadow: 0 12px 28px rgba(16, 33, 39, 0.08);
      --radius: 16px;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: "Manrope", "Segoe UI", sans-serif;
      color: var(--text);
      background: var(--bg);
    }
    .layout {
      min-height: 100vh;
      display: grid;
      grid-template-columns: 260px 1fr;
    }
    .sidebar {
      background: linear-gradient(180deg, #ffffff 0%, #f8fcfb 100%);
      border-right: 1px solid var(--border);
      padding: 22px 14px;
      position: sticky;
      top: 0;
      height: 100vh;
      overflow: auto;
      display: flex;
      flex-direction: column;
    }
    .sidebar-top {
      display: grid;
      gap: 10px;
    }
    .sidebar-bottom {
      margin-top: auto;
      padding-top: 14px;
    }
    .sidebar-actions {
      display: grid;
      gap: 8px;
    }
    .sidebar-divider {
      border-top: 1px solid var(--border);
      margin-bottom: 10px;
    }
    .brand { margin: 6px 10px 18px; }
    .brand-logo {
      display: block;
      width: 100px;
      height: auto;
      max-width: 100%;
    }
    .nav-link {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 100%;
      border: 1px solid var(--border);
      background: var(--surface);
      border-radius: 12px;
      padding: 11px 12px;
      color: var(--text);
      text-decoration: none;
      font-size: 0.95rem;
      font-weight: 600;
      text-align: center;
      transition: background-color 120ms ease, border-color 120ms ease, transform 120ms ease;
    }
    .nav-link:hover {
      border-color: var(--accent);
      transform: translateY(-1px);
    }
    .nav-link.active {
      background: var(--accent-soft);
      border-color: var(--accent);
      color: var(--accent);
      font-weight: 800;
    }
    .main {
      min-width: 0;
      padding: 24px;
      display: flex;
      flex-direction: column;
      gap: 14px;
    }
    .main-mobile-brand,
    .main-mobile-actions {
      display: none;
    }
    .header {
      background: var(--surface);
      border: 1px solid var(--border);
      border-radius: var(--radius);
      padding: 14px 16px;
      box-shadow: var(--shadow);
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 12px;
    }
    .header-left { min-width: 0; }
    .title {
      margin: 0;
      font-size: clamp(1.05rem, 0.92rem + 0.8vw, 1.55rem);
      line-height: 1.12;
      font-weight: 800;
      letter-spacing: -0.01em;
    }
    .subtitle {
      margin-top: 4px;
      display: flex;
      flex-wrap: wrap;
      gap: 4px;
      align-items: baseline;
    }
    .subtitle-text {
      color: var(--muted);
      font-size: 0.82rem;
      font-weight: 700;
    }
    .summary-pill {
      align-self: center;
      border: 1px solid var(--border);
      background: var(--surface-soft);
      border-radius: 999px;
      padding: 8px 12px;
      color: var(--accent);
      font-size: 0.84rem;
      font-weight: 700;
      white-space: nowrap;
    }
    .loading-indicator {
      display: none;
      align-items: center;
      gap: 10px;
      color: var(--muted);
      font-size: 0.9rem;
      font-weight: 600;
    }
    .spinner {
      width: 16px;
      height: 16px;
      border: 2px solid rgba(22, 101, 52, 0.18);
      border-top-color: var(--accent);
      border-radius: 50%;
      animation: spin 700ms linear infinite;
      flex-shrink: 0;
    }
    @keyframes spin { to { transform: rotate(360deg); } }
    .table-card {
      background: var(--surface);
      border: 1px solid var(--border);
      border-radius: var(--radius);
      box-shadow: 0 8px 18px rgba(16, 33, 39, 0.04);
      padding: 14px;
      overflow: hidden;
    }
    .table-head {
      display: grid;
      gap: 4px;
      margin-bottom: 12px;
    }
    .table-head h2 {
      margin: 0;
      font-size: 1rem;
      font-weight: 800;
      letter-spacing: -0.01em;
    }
    .table-head p {
      margin: 0;
      color: var(--muted);
      font-size: 0.86rem;
      line-height: 1.4;
    }
    .table-wrap {
      overflow-x: auto;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      min-width: 420px;
    }
    th, td {
      padding: 11px 10px;
      border-bottom: 1px solid rgba(215, 224, 216, 0.8);
      text-align: left;
      vertical-align: middle;
    }
    th {
      color: var(--muted);
      font-size: 0.72rem;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    td {
      font-size: 0.92rem;
      font-weight: 600;
      color: var(--text);
    }
    tbody tr:last-child td {
      border-bottom: 0;
    }
    .zone-pill {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      padding: 0.35rem 0.6rem;
      border-radius: 999px;
      background: rgba(22, 101, 52, 0.08);
      color: var(--accent);
      font-weight: 800;
      font-size: 0.84rem;
    }
    .rate-cell {
      display: grid;
      gap: 6px;
    }
    .rate-bar {
      width: 100%;
      height: 8px;
      border-radius: 999px;
      background: #e7efed;
      overflow: hidden;
    }
    .rate-bar span {
      display: block;
      height: 100%;
      width: 0%;
      border-radius: inherit;
      background: linear-gradient(90deg, #166534 0%, #22c55e 100%);
    }
    .error {
      border: 1px solid rgba(185, 28, 28, 0.22);
      background: rgba(185, 28, 28, 0.06);
      color: #9f1239;
      border-radius: 12px;
      padding: 10px 12px;
      font-size: 0.9rem;
      font-weight: 600;
      display: none;
    }
    @media (pointer: coarse) and (hover: none) {
      .layout {
        grid-template-columns: 1fr;
      }
      .sidebar {
        display: none;
      }
      .main {
        padding: 16px 12px 20px;
      }
      .main-mobile-brand,
      .main-mobile-actions {
        display: flex;
      }
      .main-mobile-brand {
        align-items: center;
        justify-content: flex-start;
      }
      .main-mobile-brand .brand-logo {
        width: 92px;
      }
      .main-mobile-actions {
        gap: 6px;
        margin-bottom: 8px;
      }
      .main-mobile-actions .nav-link {
        width: auto;
        flex: 1 1 0;
        padding: 7px 9px;
        font-size: 0.82rem;
      }
      .header {
        padding: 10px 12px;
        flex-direction: column;
        align-items: stretch;
      }
      .summary-pill {
        align-self: flex-start;
      }
    }
  </style>
</head>
<body>
  <div class="layout">
    <aside class="sidebar">
      <div class="sidebar-top">
        <div class="brand">
          <img class="brand-logo" src="/static/images/tk.png" alt="TEBMA KANDU logo" />
        </div>
        <nav class="sidebar-actions" aria-label="Primary">
          <a class="nav-link" href="/">Home</a>
          <a class="nav-link" href="/logs">Daily Logs</a>
          <a class="nav-link" href="/recovery">Recovery</a>
          <a class="nav-link active" href="/analytics">Analytics</a>
        </nav>
      </div>
    </aside>
    <main class="main" id="analyticsMain">
      <div class="main-mobile-brand" hidden>
        <img class="brand-logo" src="/static/images/tk.png" alt="TEBMA KANDU logo" />
      </div>
      <div class="main-mobile-actions" hidden>
        <a class="nav-link" href="/">Home</a>
        <a class="nav-link" href="/logs">Daily Logs</a>
        <a class="nav-link" href="/recovery">Recovery</a>
        <a class="nav-link active" href="/analytics">Analytics</a>
      </div>
      <div class="loading-indicator" role="status" aria-live="polite">
        <span class="spinner" aria-hidden="true"></span>
        <span>Loading zone table...</span>
      </div>
      <section class="header">
        <div class="header-left">
          <h1 class="title">Analytics</h1>
          <div class="subtitle">
            <span class="subtitle-text">Zone performance overview</span>
          </div>
        </div>
       
      </section>
      <div class="error" id="errorBox"></div>
      <section class="table-card" aria-label="Zone analytics table">
        <div class="table-head">
          <h2>Zone Recovery Table</h2>
          <p>Compare each zone by total nuts and recovery rate.</p>
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Zone</th>
                <th>Total Nuts</th>
                <th>Recovery Rate</th>
              </tr>
            </thead>
            <tbody id="zoneRows"></tbody>
          </table>
        </div>
      </section>
    </main>
  </div>
  <script>
    (() => {
      const zones = ["Wa", "Yendi", "Tamale", "Sandema", "Garu", "Langbinsi", "Napkaduri"];
      const loadingIndicator = document.querySelector(".loading-indicator");
      const errorBox = document.getElementById("errorBox");
      const zoneRows = document.getElementById("zoneRows");
      const zoneCount = document.getElementById("zoneCount");

      const numberFormatter = new Intl.NumberFormat(undefined, { maximumFractionDigits: 2 });

      function formatNumber(value) {
        return numberFormatter.format(Number(value || 0));
      }

      function recoveryRate(item) {
        const prefinance = Number(item.totalPrefinance || 0);
        const balance = Number(item.totalBalance || 0);
        if (prefinance <= 0) return 0;
        return Math.max(((prefinance - balance) / prefinance) * 100, 0);
      }

      function setLoading(isLoading) {
        loadingIndicator.style.display = isLoading ? "flex" : "none";
      }

      function setError(message) {
        if (!message) {
          errorBox.style.display = "none";
          errorBox.textContent = "";
          return;
        }
        errorBox.textContent = message;
        errorBox.style.display = "block";
      }

      function renderRows(rows) {
        zoneRows.innerHTML = rows.map(function (row) {
          const rate = recoveryRate(row);
          return [
            "<tr>",
            "<td><span class=\"zone-pill\">" + row.zoneName + "</span></td>",
            "<td><strong>" + formatNumber(row.totalKgBrought) + "</strong> kg</td>",
            "<td><div class=\"rate-cell\"><strong>" + rate.toFixed(1) + "%</strong>",
            "<div class=\"rate-bar\" aria-hidden=\"true\"><span style=\"width:" + Math.min(rate, 100) + "%\"></span></div>",
            "</div></td>",
            "</tr>"
          ].join("");
        }).join("");
      }

      async function fetchZoneOverview(zone) {
        const response = await fetch("/api/zones/" + encodeURIComponent(zone) + "/farmers/overview", {
          headers: { "Accept": "application/json" }
        });
        if (!response.ok) {
          throw new Error("Failed to load " + zone);
        }
        const data = await response.json();
        return { ...data, zoneName: zone };
      }

      async function loadTable() {
        try {
          setLoading(true);
          setError("");

          const settled = await Promise.allSettled(zones.map(fetchZoneOverview));
          const rows = settled
            .filter((item) => item.status === "fulfilled")
            .map((item) => item.value);

          if (!rows.length) {
            throw new Error("No zone analytics data could be loaded.");
          }

          const sortedRows = [...rows].sort((a, b) => {
            const rateDiff = recoveryRate(b) - recoveryRate(a);
            if (rateDiff !== 0) return rateDiff;
            return Number(b.totalKgBrought || 0) - Number(a.totalKgBrought || 0);
          });

          zoneCount.textContent = rows.length + " zones compared";
          renderRows(sortedRows);
        } catch (err) {
          setError(err instanceof Error ? err.message : "Failed to build zone table.");
        } finally {
          setLoading(false);
        }
      }

      loadTable();
    })();
  </script>
</body>
</html>`)
		return err
	})
}
