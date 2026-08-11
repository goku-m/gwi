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
  <script src="https://cdn.tailwindcss.com"></script>
  <link href="https://cdnjs.cloudflare.com/ajax/libs/flowbite/2.5.2/flowbite.min.css" rel="stylesheet" />
  <script src="https://unpkg.com/htmx.org@2.0.4"></script>
  <script src="https://cdn.jsdelivr.net/npm/apexcharts"></script>
  <style>
    :root {
      --bg: #f3f7f6;
      --surface: #ffffff;
      --surface-soft: #f8fbfa;
      --text: #102127;
      --muted: #56707d;
      --accent: #0f766e;
      --accent-soft: #e8f6f4;
      --border: #d8e5e3;
      --shadow: 0 12px 28px rgba(16, 33, 39, 0.08);
      --radius: 16px;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: "Manrope", "Segoe UI", sans-serif;
      color: var(--text);
      background:
        radial-gradient(circle at 12% 12%, #dff2ef 0%, transparent 40%),
        linear-gradient(180deg, #f9fcfb 0%, var(--bg) 100%);
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
    .brand {
      margin: 6px 10px 18px;
    }
    .brand-logo {
      display: block;
      width: 100px;
      height: auto;
      max-width: 100%;
    }
    .zone-list {
      display: grid;
      gap: 8px;
    }
    .nav-link {
      display: block;
      width: 100%;
      border: 1px solid var(--border);
      background: var(--surface);
      border-radius: 12px;
      padding: 11px 12px;
      color: var(--text);
      text-decoration: none;
      font-size: 0.95rem;
      font-weight: 600;
      transition: background-color 120ms ease, border-color 120ms ease, transform 120ms ease;
    }
    .nav-link:hover {
      border-color: #a7cfc8;
      transform: translateY(-1px);
    }
    .nav-link.active {
      background: var(--accent-soft);
      border-color: #7dcfc5;
      color: #0a5d56;
    }
    .mobile-logs-link {
      display: none;
      margin-bottom: 10px;
    }
    .mobile-actions {
      display: none;
      gap: 8px;
      margin-bottom: 10px;
    }
    .mobile-actions .nav-link {
      width: auto;
      flex: 1 1 0;
    }
    .main-mobile-brand {
      display: none;
    }
    .main-mobile-actions {
      display: none;
    }
    .zone-mobile-picker {
      display: none;
      margin: 0 10px 10px;
    }
    .community-zone-picker {
      display: none;
      gap: 6px;
    }
    .zone-mobile-select {
      width: 100%;
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 10px 12px;
      background: var(--surface-soft);
      color: var(--text);
      font-size: 0.92rem;
      font-weight: 600;
    }
    .zone-btn {
      width: 100%;
      border: 1px solid var(--border);
      background: var(--surface);
      border-radius: 12px;
      padding: 11px 12px;
      text-align: left;
      color: var(--text);
      cursor: pointer;
      font-size: 0.95rem;
      font-weight: 600;
      transition: background-color 120ms ease, border-color 120ms ease, transform 120ms ease;
    }
    .zone-btn:hover { border-color: #a7cfc8; transform: translateY(-1px); }
    .zone-btn.active {
      background: var(--accent-soft);
      border-color: #7dcfc5;
      color: #0a5d56;
      box-shadow: 0 6px 14px rgba(15, 118, 110, 0.12);
    }
    .main {
      padding: 24px;
      max-width: 1280px;
      width: 100%;
    }
    .header {
      background: var(--surface);
      border: 1px solid var(--border);
      border-radius: var(--radius);
      padding: 8px 12px;
      box-shadow: var(--shadow);
      margin-bottom: 10px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 8px;
    }
    .header-left { min-width: 0; }
    .title {
      margin: 0;
      font-size: clamp(0.98rem, 0.88rem + 0.7vw, 1.35rem);
      line-height: 1.08;
      font-weight: 800;
      letter-spacing: -0.01em;
    }
    .subtitle {
      margin-top: 3px;
      color: var(--muted);
      font-size: 0.78rem;
      min-height: 0.9rem;
      display: flex;
      flex-wrap: wrap;
      align-items: baseline;
      gap: 4px;
    }
    .subtitle-syncs {
      color: #0a5d56;
      font-size: 0.76rem;
      font-weight: 700;
    }
    .community-picker {
      min-width: 220px;
      display: grid;
      gap: 6px;
    }
    .community-label {
      color: var(--muted);
      font-size: 0.8rem;
      font-weight: 700;
    }
    .community-input {
      width: 100%;
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 7px 9px;
      background: var(--surface-soft);
      color: var(--text);
      font-size: 0.86rem;
    }
    .community-input:focus {
      outline: none;
      border-color: #8ad2ca;
      box-shadow: 0 0 0 3px rgba(138, 210, 202, 0.25);
    }
    .charts {
      display: grid;
      gap: 14px;
      grid-template-columns: repeat(2, minmax(240px, 1fr));
    }
    .chart-card {
      background: linear-gradient(180deg, #ffffff 0%, #f9fcfb 100%);
      border: 1px solid var(--border);
      border-radius: var(--radius);
      box-shadow: var(--shadow);
      padding: 14px;
    }
    .chart-title {
      margin: 0 0 10px;
      font-size: 0.9rem;
      font-weight: 700;
      letter-spacing: 0.02em;
      color: var(--muted);
    }
    .chart-meta {
      margin: -2px 0 10px;
      font-size: 0.86rem;
      font-weight: 700;
      color: #0a5d56;
    }
    .loading-indicator {
      display: none;
      align-items: center;
      gap: 8px;
      margin: 10px 0 14px;
      color: var(--muted);
      font-size: 0.9rem;
      font-weight: 600;
    }
    .loading .loading-indicator {
      display: flex;
    }
    .spinner {
      width: 16px;
      height: 16px;
      border: 2px solid #c9dcda;
      border-top-color: var(--accent);
      border-radius: 50%;
      animation: spin 700ms linear infinite;
    }
    @keyframes spin {
      to { transform: rotate(360deg); }
    }
    @media (prefers-reduced-motion: reduce) {
      .zone-btn, .card { transition: none; }
    }
    @media (max-width: 960px) {
      .layout { grid-template-columns: 1fr; }
      .main {
        display: flex;
        flex-direction: column;
      }
      .sidebar {
        display: none;
      }
      .main-mobile-brand {
        display: flex;
        align-items: center;
        justify-content: flex-start;
        margin-bottom: 8px;
      }
      .main-mobile-brand .brand-logo {
        width: 92px;
      }
      .main-mobile-actions {
        display: flex;
        gap: 6px;
        margin-bottom: 8px;
      }
      .main-mobile-actions .nav-link {
        width: auto;
        flex: 1 1 0;
        padding: 7px 9px;
        font-size: 0.82rem;
      }
      .zone-list {
        display: none;
      }
      .header {
        display: flex;
        flex-direction: column;
        align-items: stretch;
        gap: 6px;
        padding: 8px 10px;
      }
      .header-left,
      .community-picker {
        width: 100%;
        min-width: 0;
      }
      .community-picker {
        margin-top: 0;
      }
      .community-zone-picker { display: grid; }
      .charts { grid-template-columns: 1fr; }
    }
    @media (max-width: 560px) {
      .main { padding: 14px; }
      .title { font-size: 0.98rem; }
      .main-mobile-brand .brand-logo { width: 84px; }
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
        <nav class="zone-list" id="zoneList">
          <button class="zone-btn active" data-zone="General">General</button>
          <button class="zone-btn" data-zone="Wa">Wa</button>
          <button class="zone-btn" data-zone="Yendi">Yendi</button>
          <button class="zone-btn" data-zone="Tamale">Tamale</button>
          <button class="zone-btn" data-zone="Sandema">Sandema</button>
          <button class="zone-btn" data-zone="Garu">Garu</button>
          <button class="zone-btn" data-zone="Langbinsi">Langbinsi</button>
          <button class="zone-btn" data-zone="Napkaduri">Napkanduri</button>
        </nav>
      </div>
      <div class="sidebar-bottom">
        <div class="sidebar-divider"></div>
        <div class="sidebar-actions">
          <a class="nav-link" href="/">Home</a>
          <a class="nav-link" href="/logs">Logs</a>
        </div>
      </div>
    </aside>
    <main class="main" id="analyticsMain">
      <div class="main-mobile-brand">
        <img class="brand-logo" src="/static/images/tk.png" alt="TEBMA KANDU logo" />
      </div>
      <div class="main-mobile-actions">
        <a class="nav-link" href="/">Home</a>
        <a class="nav-link" href="/logs">Daily Logs</a>
        <a class="nav-link active" href="/analytics">Analytics</a>
      </div>
      <section class="header">
        <div class="header-left">
          <h1 class="title" id="selectedZoneTitle">Analytics</h1>
          <div class="subtitle">
            <span id="subtitleText">Zone-level Stats</span>
            <span class="subtitle-syncs" id="subtitleSyncs">Syncs Today: <span id="dailySyncs">0</span></span>
          </div>
        </div>
        <div class="community-picker">
          <div class="community-zone-picker">
            <label class="community-label" for="zoneMobileSelect">Select Zone</label>
            <select id="zoneMobileSelect" class="zone-mobile-select" aria-label="Zone select">
              <option value="General">General</option>
              <option value="Wa">Wa</option>
              <option value="Yendi">Yendi</option>
              <option value="Tamale">Tamale</option>
              <option value="Sandema">Sandema</option>
              <option value="Garu">Garu</option>
              <option value="Langbinsi">Langbinsi</option>
              <option value="Napkaduri">Napkanduri</option>
            </select>
          </div>
        </div>
      </section>
      <div class="loading-indicator" role="status" aria-live="polite">
        <span class="spinner" aria-hidden="true"></span>
        <span>Fetching latest analytics...</span>
      </div>
      <section class="charts">
        <article class="chart-card">
          <p class="chart-title">Recovery Rate</p>
          <p class="chart-meta" id="recoveryRateText">Recovery: 0%</p>
          <div id="recoveryChart"></div>
        </article>
        <article class="chart-card">
          <p class="chart-title">Financial Breakdown</p>
          <div id="financeChart"></div>
        </article>
      </section>
    </main>
  </div>
  <script>
    (function () {
      const zoneButtons = Array.from(document.querySelectorAll(".zone-btn"));
      const selectedZoneTitle = document.getElementById("selectedZoneTitle");
      const subtitleText = document.getElementById("subtitleText");
      const subtitleSyncs = document.getElementById("subtitleSyncs");
      const dailySyncs = document.getElementById("dailySyncs");
      const recoveryRateText = document.getElementById("recoveryRateText");
      const analyticsMain = document.getElementById("analyticsMain");
      const zoneMobileSelect = document.getElementById("zoneMobileSelect");

      let selectedZone = "General";
      let recoveryChart = null;
      let financeChart = null;
      let activeStatsRequest = null;

      function formatNumber(value, maxFractionDigits) {
        const n = Number(value || 0);
        return n.toLocaleString(undefined, { maximumFractionDigits: maxFractionDigits });
      }

      function setLoading() {
        analyticsMain.classList.add("loading");
        recoveryRateText.textContent = "Recovery: ...";
        dailySyncs.textContent = "...";
      }

      function renderCharts(data) {
        const recoveryChartEl = document.querySelector("#recoveryChart");
        const financeChartEl = document.querySelector("#financeChart");

        if (!recoveryChartEl || !financeChartEl || !recoveryRateText) {
          return;
        }

        const amount = Number(data.totalAmount || 0);
        const prefinance = Number(data.totalPrefinance || 0);
        const balance = Number(data.totalBalance || 0);

        const unpaidPrefinance = Math.max(0, Math.min(balance, prefinance));
        const recoveredPrefinance = Math.max(0, prefinance - unpaidPrefinance);
        const recoveryPercent = prefinance > 0 ? (recoveredPrefinance / prefinance) * 100 : 0;
        recoveryRateText.textContent = "Recovery: " + recoveryPercent.toFixed(1) + "%";

        const recoveryOptions = {
          chart: { type: "pie", height: 280 },
          series: [recoveredPrefinance, unpaidPrefinance],
          labels: ["Recovered Prefinance", "Unpaid Prefinance"],
          colors: ["#0f766e", "#f59e0b"],
          legend: { position: "bottom" }
        };

        const financeOptions = {
          chart: { type: "donut", height: 280 },
          series: [amount, prefinance, balance],
          labels: ["Total Amount", "Prefinance", "Balance"],
          colors: ["#0f766e", "#0ea5e9", "#f59e0b"],
          legend: { position: "bottom" }
        };

        if (!recoveryChart) {
          recoveryChart = new ApexCharts(recoveryChartEl, recoveryOptions);
          recoveryChart.render();
        } else {
          recoveryChart.updateOptions(recoveryOptions);
        }

        if (!financeChart) {
          financeChart = new ApexCharts(financeChartEl, financeOptions);
          financeChart.render();
        } else {
          financeChart.updateOptions(financeOptions);
        }
      }

      function normalizeForCompare(value) {
        return (value || "").trim().toLowerCase();
      }

      async function fetchStats(route) {
        if (activeStatsRequest) {
          activeStatsRequest.abort();
        }

        activeStatsRequest = new AbortController();
        const response = await fetch(route, { signal: activeStatsRequest.signal });
        if (!response.ok) {
          throw new Error("Request failed with status " + response.status);
        }

        return response.json();
      }

      async function loadZoneAnalytics(zone) {
        selectedZoneTitle.textContent = zone;
        subtitleText.textContent = "Live zone-level analytics";
        subtitleSyncs.style.display = "inline-flex";
        setLoading();

        const route = zone === "General"
          ? "/api/farmers/stats"
          : "/api/zones/" + encodeURIComponent(zone) + "/farmers/stats";

        try {
          const data = await fetchStats(route);
          dailySyncs.textContent = formatNumber(data.dailySyncs, 0);
          analyticsMain.classList.remove("loading");
          renderCharts(data);
        } catch (err) {
          if (err.name === "AbortError") {
            return;
          }
          analyticsMain.classList.remove("loading");
        }
      }

      zoneButtons.forEach(function (button) {
        button.addEventListener("click", function () {
          zoneButtons.forEach(function (b) { b.classList.remove("active"); });
          button.classList.add("active");
          selectedZone = button.dataset.zone;
          zoneMobileSelect.value = selectedZone;
          loadZoneAnalytics(selectedZone);
        });
      });

      zoneMobileSelect.addEventListener("change", function () {
        const zone = zoneMobileSelect.value;
        selectedZone = zone;
        zoneButtons.forEach(function (b) {
          b.classList.toggle("active", b.dataset.zone === zone);
        });
        loadZoneAnalytics(selectedZone);
      });
      loadZoneAnalytics("General");
    })();
  </script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/flowbite/2.5.2/flowbite.min.js"></script>
</body>
</html>`)
		return err
	})
}
