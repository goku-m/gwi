package pages

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func Home() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>TEBMA KANDU</title>
  <script>
    (function () {
      var ua = navigator.userAgent || "";
      if (/Mobi|Android|iPhone|iPad|iPod/i.test(ua)) {
        document.documentElement.classList.add("mobile-layout");
      }
    })();
  </script>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Manrope:wght@400;500;600;700;800&display=swap" rel="stylesheet">
  <script src="https://cdn.tailwindcss.com"></script>
  <link href="https://cdnjs.cloudflare.com/ajax/libs/flowbite/2.5.2/flowbite.min.css" rel="stylesheet" />
  <script src="https://unpkg.com/htmx.org@2.0.4"></script>
  <script src="https://cdn.jsdelivr.net/npm/apexcharts"></script>
  <style>
    :root {
      --bg: #f3f4f6;
      --surface: #ffffff;
      --surface-soft: #f8fafc;
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
      background: linear-gradient(180deg, #ffffff 0%, #f8fafc 100%);
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
    .mobile-community-picker {
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
      display: flex;
      align-items: center;
      justify-content: center;
      width: 100%;
      border: 1px solid var(--border);
      background: var(--surface);
      border-radius: 12px;
      padding: 11px 12px;
      text-align: center;
      color: var(--text);
      cursor: pointer;
      font-size: 0.95rem;
      font-weight: 600;
      transition: background-color 120ms ease, border-color 120ms ease, transform 120ms ease;
    }
    .zone-btn:hover { border-color: var(--accent); transform: translateY(-1px); }
    .zone-btn:focus-visible {
      outline: 2px solid rgba(22, 101, 52, 0.35);
      outline-offset: 1px;
    }
    .zone-btn.active {
      background: var(--accent-soft);
      border-color: var(--accent);
      color: var(--accent);
      font-weight: 800;
      box-shadow: 0 6px 14px rgba(22, 101, 52, 0.12);
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
      padding: 11px 14px;
      box-shadow: var(--shadow);
      margin-bottom: 12px;
      display: flex;
      align-items: center;
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
      color: var(--muted);
      font-size: 0.84rem;
      min-height: 1rem;
      display: flex;
      flex-wrap: wrap;
      align-items: baseline;
      gap: 6px;
    }
    .subtitle-syncs {
      color: var(--accent);
      font-size: 0.82rem;
      font-weight: 700;
    }
    .header-metrics {
      margin-top: 8px;
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }
    .header-metric {
      background: var(--surface-soft);
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 6px 8px;
      min-width: 122px;
    }
    .header-metric-label {
      display: block;
      color: var(--muted);
      font-size: 0.68rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.04em;
      margin-bottom: 3px;
    }
    .header-metric-value {
      margin: 0;
      font-size: 0.95rem;
      font-weight: 800;
      color: #0d3a36;
    }
    .metric-increment {
      color: var(--accent);
      font-weight: 800;
      margin-left: 6px;
    }
    .community-picker {
      min-width: 220px;
      display: grid;
      gap: 6px;
    }
    .filters-row {
      display: grid;
      gap: 8px;
      grid-template-columns: repeat(2, minmax(130px, 1fr));
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
      padding: 10px 12px;
      background: var(--surface-soft);
      color: var(--text);
      font-size: 0.92rem;
    }
    .community-input:focus {
      outline: none;
      border-color: var(--accent);
      box-shadow: 0 0 0 3px rgba(22, 101, 52, 0.18);
    }
    .community-select {
      display: none;
      width: 100%;
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 10px 12px;
      background: var(--surface-soft);
      color: var(--text);
      font-size: 0.92rem;
    }
    .date-input {
      width: 100%;
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 10px 12px;
      background: var(--surface-soft);
      color: var(--text);
      font-size: 0.9rem;
    }
    .cards {
      display: grid;
      gap: 12px;
      grid-template-columns: minmax(0, 1fr);
      margin-bottom: 14px;
    }
    .card {
      background: linear-gradient(180deg, #ffffff 0%, #f9fcfb 100%);
      border: 1px solid var(--border);
      border-radius: var(--radius);
      box-shadow: var(--shadow);
      padding: 12px;
      min-height: 108px;
      display: grid;
      gap: 6px;
    }
    .card-label {
      display: block;
      margin: 0;
      color: var(--muted);
      font-size: 0.76rem;
      font-weight: 800;
      line-height: 1.1;
      text-transform: uppercase;
      letter-spacing: 0.06em;
    }
    .card-value {
      margin: 0;
      font-size: 1.2rem;
      line-height: 1.05;
      font-weight: 800;
      letter-spacing: -0.02em;
      color: var(--text);
    }
    .currency-symbol {
      font-size: 0.78em;
      font-weight: 700;
      vertical-align: baseline;
      margin-right: 0.18em;
    }
    .currency-amount {
      font-size: 1.1em;
    }
    .card-subvalue {
      color: var(--muted);
      font-size: 0.82rem;
      line-height: 1.25;
    }
    #totalKgBrought,
    #totalAmount {
      color: var(--accent);
    }
    #totalKgBrought,
    #totalAmount {
      font-size: 1.22rem;
    }
    .error {
      margin-top: 12px;
      color: #8f1d1d;
      font-size: 0.9rem;
      display: none;
    }
    .charts {
      margin-top: 16px;
      display: grid;
      gap: 14px;
      grid-template-columns: repeat(2, minmax(220px, 1fr));
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
      color: var(--accent);
    }
    .loading .card-value {
      color: #8ca2ad;
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
      border: 2px solid rgba(22, 101, 52, 0.18);
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
    .main-mobile-brand,
    .main-mobile-actions,
    .mobile-community-picker {
      display: none;
    }
    html.mobile-layout .layout { grid-template-columns: 1fr; }
    html.mobile-layout .main {
      display: flex;
      flex-direction: column;
    }
    html.mobile-layout .sidebar {
      display: none;
    }
    html.mobile-layout .main-mobile-brand {
      display: flex;
      align-items: center;
      justify-content: flex-start;
      margin-bottom: 8px;
    }
    html.mobile-layout .main-mobile-brand .brand-logo {
      width: 92px;
    }
    html.mobile-layout .main-mobile-actions {
      display: flex;
      gap: 6px;
      margin-bottom: 8px;
    }
    html.mobile-layout .main-mobile-actions .nav-link {
      width: auto;
      flex: 1 1 0;
      padding: 7px 9px;
      font-size: 0.82rem;
    }
    html.mobile-layout .zone-list {
      display: none;
    }
    html.mobile-layout .header {
      display: flex;
      flex-direction: column;
      align-items: stretch;
      gap: 8px;
      padding: 10px 12px;
    }
    html.mobile-layout .header-left,
    html.mobile-layout .cards,
    html.mobile-layout .community-picker {
      width: 100%;
      min-width: 0;
    }
    html.mobile-layout .header-metrics {
      display: grid;
      grid-template-columns: repeat(2, minmax(130px, 1fr));
      gap: 8px;
      margin-top: 8px;
    }
    html.mobile-layout .header-metric {
      min-width: 0;
      padding: 5px 7px;
    }
    html.mobile-layout .cards {
      order: 2;
      margin-top: 8px;
      grid-template-columns: 1fr;
    }
    html.mobile-layout .community-picker {
      order: 3;
      margin-top: 0;
    }
    html.mobile-layout .header .community-picker {
      display: none;
    }
    html.mobile-layout .mobile-community-picker {
      display: grid;
      gap: 6px;
      margin-top: 8px;
      order: 4;
      background: var(--surface-soft);
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 10px 12px;
    }
    html.mobile-layout .community-zone-picker { display: grid; }
    html.mobile-layout .community-input { display: none; }
    html.mobile-layout .community-select { display: block; }
    html.mobile-layout .filters-row { grid-template-columns: 1fr 1fr; }
    @media (max-width: 560px) {
      .main { padding: 14px; }
      .title { font-size: 1.2rem; }
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
          <a class="nav-link" href="/logs">Logs</a>
          <a class="nav-link" href="/recovery">Recovery</a>
        </div>
      </div>
    </aside>
    <main class="main" id="dashboardMain">
      <div class="main-mobile-brand" hidden>
        <img class="brand-logo" src="/static/images/tk.png" alt="TEBMA KANDU logo" />
      </div>
      <div class="main-mobile-actions" hidden>
        <a class="nav-link" href="/logs">Daily Logs</a>
        <a class="nav-link" href="/recovery">Recovery</a>
      </div>
       <div class="loading-indicator" role="status" aria-live="polite">
        <span class="spinner" aria-hidden="true"></span>
        <span>Fetching latest data...</span>
      </div>
      <section class="header">
        <div class="header-left">
          <h1 class="title" id="selectedZoneTitle">General</h1>
          <div class="subtitle">
            <span id="subtitleText">Zone-level Stats</span>
            <span class="subtitle-syncs" id="subtitleSyncs">Syncs Today: <span id="dailySyncs">0</span></span>
          </div>
          <div class="header-metrics">
            <div class="header-metric">
              <span class="header-metric-label">Farmers</span>
              <p class="header-metric-value" id="totalFarmers">0</p>
            </div>
            <div class="header-metric">
              <span class="header-metric-label">Communities</span>
              <p class="header-metric-value" id="totalCommunities">0</p>
            </div>
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
          <label class="community-label" for="communitySearch">Community</label>
          <input id="communitySearch" class="community-input" list="communityOptions" placeholder="All communities" />
          <select id="communitySelect" class="community-select" aria-label="Community select"></select>
          <datalist id="communityOptions"></datalist>
          <div class="filters-row">
            <div>
              <label class="community-label" for="fromDate">From</label>
              <input id="fromDate" class="date-input" type="date" />
            </div>
            <div>
              <label class="community-label" for="toDate">To</label>
              <input id="toDate" class="date-input" type="date" />
            </div>
          </div>
        </div>
      </section>

      <section class="cards">
        <article class="card">
          <p class="card-label">Total Nuts</p>
          <p class="card-value" id="totalKgBrought">0</p>
          <span class="card-subvalue">Total kilograms brought in</span>
        </article>
        <article class="card">
          <p class="card-label">Nuts Value</p>
          <p class="card-value" id="totalAmount">0</p>
          <span class="card-subvalue">Gross value of nuts brought in</span>
        </article>
      </section>
      <div class="community-picker mobile-community-picker" hidden>
        <div class="community-zone-picker">
          <label class="community-label" for="zoneMobileSelectBottom">Select Zone</label>
          <select id="zoneMobileSelectBottom" class="zone-mobile-select" aria-label="Zone select">
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
        <label class="community-label" for="communitySearchBottom">Community</label>
        <input id="communitySearchBottom" class="community-input" list="communityOptionsBottom" placeholder="All communities" />
        <select id="communitySelectBottom" class="community-select" aria-label="Community select"></select>
        <datalist id="communityOptionsBottom"></datalist>
        <div class="filters-row">
          <div>
            <label class="community-label" for="fromDateBottom">From</label>
            <input id="fromDateBottom" class="date-input" type="date" />
          </div>
          <div>
            <label class="community-label" for="toDateBottom">To</label>
            <input id="toDateBottom" class="date-input" type="date" />
          </div>
        </div>
      </div>
      <div class="error" id="errorBox" role="status" aria-live="polite"></div>
    </main>
  </div>

  <script>
    (function () {
      const zoneButtons = Array.from(document.querySelectorAll(".zone-btn"));
      const selectedZoneTitle = document.getElementById("selectedZoneTitle");
      const subtitleText = document.getElementById("subtitleText");
      const subtitleSyncs = document.getElementById("subtitleSyncs");
      const errorBox = document.getElementById("errorBox");
      const communitySearch = document.getElementById("communitySearch");
      const communitySelect = document.getElementById("communitySelect");
      const zoneMobileSelect = document.getElementById("zoneMobileSelect");
      const communityOptions = document.getElementById("communityOptions");
      const fromDate = document.getElementById("fromDate");
      const toDate = document.getElementById("toDate");
      const zoneMobileSelectBottom = document.getElementById("zoneMobileSelectBottom");
      const communitySearchBottom = document.getElementById("communitySearchBottom");
      const communitySelectBottom = document.getElementById("communitySelectBottom");
      const communityOptionsBottom = document.getElementById("communityOptionsBottom");
      const fromDateBottom = document.getElementById("fromDateBottom");
      const toDateBottom = document.getElementById("toDateBottom");
      const dashboardMain = document.getElementById("dashboardMain");
      const mainMobileBrand = document.querySelector(".main-mobile-brand");
      const mainMobileActions = document.querySelector(".main-mobile-actions");
      const mobileCommunityPicker = document.querySelector(".mobile-community-picker");

      if (document.documentElement.classList.contains("mobile-layout")) {
        [mainMobileBrand, mainMobileActions, mobileCommunityPicker].forEach(function (element) {
          if (element) {
            element.hidden = false;
          }
        });
      }

      const totalFarmers = document.getElementById("totalFarmers");
      const totalCommunities = document.getElementById("totalCommunities");
      const dailySyncs = document.getElementById("dailySyncs");
      const totalKgBrought = document.getElementById("totalKgBrought");
      const totalAmount = document.getElementById("totalAmount");
      const recoveryRateText = document.getElementById("recoveryRateText");
      let selectedZone = "General";
      const communitiesByZone = {};
      let recoveryChart = null;
      let financeChart = null;
      let activeStatsRequest = null;
      let latestTotalFarmers = 0;
      let latestNewFarmers = 0;

      function formatNumber(value, maxFractionDigits) {
        const n = Number(value || 0);
        return n.toLocaleString(undefined, { maximumFractionDigits: maxFractionDigits });
      }

      function formatCurrency(value) {
        return "<span class=\"currency-symbol\">GH\u20B5</span><span class=\"currency-amount\">" + formatNumber(value, 2) + "</span>";
      }

      function setLoading() {
        dashboardMain.classList.add("loading");
        totalFarmers.textContent = "...";
        totalCommunities.textContent = "...";
        dailySyncs.textContent = "...";
        totalKgBrought.textContent = "...";
        totalAmount.textContent = "...";
      }

      function setStats(data) {
        dashboardMain.classList.remove("loading");
        latestTotalFarmers = Number(data.totalFarmers || 0);
        totalFarmers.innerHTML = formatNumber(latestTotalFarmers, 0) + " <span class=\"metric-increment\">↑" + formatNumber(latestNewFarmers, 0) + "</span>";
        totalCommunities.textContent = formatNumber(data.totalCommunities, 0);
        dailySyncs.textContent = formatNumber(data.dailySyncs, 0);
        totalKgBrought.textContent = formatNumber(data.totalKgBrought, 2) + " kg";
        totalAmount.innerHTML = formatCurrency(data.totalAmount);
        renderCharts(data);
      }

      function setNewFarmers(data) {
        const count = (data && typeof data.newFarmers !== "undefined") ? data.newFarmers : 0;
        latestNewFarmers = Number(count || 0);
        totalFarmers.innerHTML = formatNumber(latestTotalFarmers, 0) + " <span class=\"metric-increment\">↑" + formatNumber(latestNewFarmers, 0) + "</span>";
      }

      function renderCharts(data) {
        const recoveryChartEl = document.querySelector("#recoveryChart");
        const financeChartEl = document.querySelector("#financeChart");
        const recoveryRateNode = document.getElementById("recoveryRateText");

        if (!recoveryChartEl || !financeChartEl || !recoveryRateNode) {
          return;
        }

        const amount = Number(data.totalAmount || 0);
        const prefinance = Number(data.totalPrefinance || 0);
        const balance = Number(data.totalBalance || 0);

        const unpaidPrefinance = Math.max(0, Math.min(balance, prefinance));
        const recoveredPrefinance = Math.max(0, prefinance - unpaidPrefinance);
        const recoveryPercent = prefinance > 0 ? (recoveredPrefinance / prefinance) * 100 : 0;
        recoveryRateNode.textContent = "Recovery: " + recoveryPercent.toFixed(1) + "%";

        const recoveryOptions = {
          chart: { type: "pie", height: 280 },
          series: [recoveredPrefinance, unpaidPrefinance],
          labels: ["Recovered Prefinance", "Unpaid Prefinance"],
          colors: ["#166534", "#84cc16"],
          legend: { position: "bottom" }
        };

        const financeOptions = {
          chart: { type: "donut", height: 280 },
          series: [amount, prefinance, balance],
          labels: ["Total Amount", "Prefinance", "Balance"],
          colors: ["#166534", "#22c55e", "#f59e0b"],
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

      function getDateQuery() {
        const from = fromDate.value;
        const to = toDate.value;

        if (from && to && from > to) {
          errorBox.textContent = "From date cannot be after To date.";
          errorBox.style.display = "block";
          return null;
        }

        const params = new URLSearchParams();
        if (from) params.set("from", from);
        if (to) params.set("to", to);
        const query = params.toString();
        return query ? ("?" + query) : "";
      }

      async function loadCommunities(zone) {
        communitySearch.value = "";
        communitySearchBottom.value = "";
        communityOptions.innerHTML = "";
        communityOptionsBottom.innerHTML = "";

        if (zone === "General") {
          communitySearch.disabled = true;
          communitySearchBottom.disabled = true;
          communitySelect.disabled = true;
          communitySelectBottom.disabled = true;
          communitySearch.placeholder = "All communities";
          communitySearchBottom.placeholder = "All communities";
          communitySelect.innerHTML = "<option value=\"\">All communities</option>";
          communitySelectBottom.innerHTML = "<option value=\"\">All communities</option>";
          return;
        }

        communitySearch.disabled = false;
        communitySearchBottom.disabled = false;
        communitySelect.disabled = false;
        communitySelectBottom.disabled = false;
        communitySearch.placeholder = "Search communities";
        communitySearchBottom.placeholder = "Search communities";

        if (communitiesByZone[zone]) {
          renderCommunityOptions(communitiesByZone[zone]);
          return;
        }

        try {
          const response = await fetch("/api/zones/" + encodeURIComponent(zone) + "/communities");
          if (!response.ok) {
            throw new Error("Request failed with status " + response.status);
          }
          const payload = await response.json();
          const seenCommunities = new Set();
          const communities = (Array.isArray(payload.communities) ? payload.communities : []).reduce(function (acc, community) {
            const cleanName = (community || "").trim();
            if (cleanName.length < 3) {
              return acc;
            }

            const dedupeKey = cleanName.toLowerCase();
            if (seenCommunities.has(dedupeKey)) {
              return acc;
            }

            seenCommunities.add(dedupeKey);
            acc.push(cleanName);
            return acc;
          }, []);
          communitiesByZone[zone] = communities;
          renderCommunityOptions(communities);
        } catch (err) {
          communitiesByZone[zone] = [];
          renderCommunityOptions([]);
        }
      }

      function renderCommunityOptions(communities) {
        communityOptions.innerHTML = "";
        communitySelect.innerHTML = "";
        communityOptionsBottom.innerHTML = "";
        communitySelectBottom.innerHTML = "";

        const defaultOption = document.createElement("option");
        defaultOption.value = "";
        defaultOption.textContent = "All communities";
        communitySelect.appendChild(defaultOption);
        communitySelectBottom.appendChild(defaultOption.cloneNode(true));

        communities.forEach(function (community) {
          const option = document.createElement("option");
          option.value = community;
          communityOptions.appendChild(option);
          communityOptionsBottom.appendChild(option.cloneNode(true));

          const selectOption = document.createElement("option");
          selectOption.value = community;
          selectOption.textContent = community;
          communitySelect.appendChild(selectOption);
          communitySelectBottom.appendChild(selectOption.cloneNode(true));
        });
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

      async function fetchNewFarmers(route) {
        if (!activeStatsRequest) {
          activeStatsRequest = new AbortController();
        }

        const response = await fetch(route, { signal: activeStatsRequest.signal });
        if (!response.ok) {
          throw new Error("Request failed with status " + response.status);
        }

        return response.json();
      }

      async function loadZoneStats(zone) {
        selectedZoneTitle.textContent = zone;
        subtitleText.textContent = "Live zone-level stats";
        subtitleSyncs.style.display = "inline-flex";
        errorBox.style.display = "none";
        errorBox.textContent = "";
        setLoading();

        const dateQuery = getDateQuery();
        if (dateQuery === null) {
          return;
        }

        const route = zone === "General"
          ? "/api/farmers/stats" + dateQuery
          : "/api/zones/" + encodeURIComponent(zone) + "/farmers/stats" + dateQuery;
        const newFarmersRoute = zone === "General"
          ? "/api/farmers/new"
          : "/api/zones/" + encodeURIComponent(zone) + "/farmers/new";

        try {
          const data = await fetchStats(route);
          const newFarmersData = await fetchNewFarmers(newFarmersRoute);
          setStats(data);
          setNewFarmers(newFarmersData);
        } catch (err) {
          if (err.name === "AbortError") {
            return;
          }
          setStats({});
          setNewFarmers({});
          errorBox.textContent = "Could not load stats for " + zone + ".";
          errorBox.style.display = "block";
        }
      }

      async function loadCommunityStats(zone, community) {
        subtitleText.textContent = "Community: " + community;
        subtitleSyncs.style.display = "none";
        errorBox.style.display = "none";
        errorBox.textContent = "";
        setLoading();

        const dateQuery = getDateQuery();
        if (dateQuery === null) {
          return;
        }

        try {
          const data = await fetchStats(
            "/api/zones/" + encodeURIComponent(zone) + "/" + encodeURIComponent(community) + "/farmers/stats" + dateQuery
          );
          const newFarmersData = await fetchNewFarmers(
            "/api/zones/" + encodeURIComponent(zone) + "/" + encodeURIComponent(community) + "/farmers/new"
          );
          setStats(data);
          setNewFarmers(newFarmersData);
        } catch (err) {
          if (err.name === "AbortError") {
            return;
          }
          setStats({});
          setNewFarmers({});
          errorBox.textContent = "Could not load stats for " + community + ".";
          errorBox.style.display = "block";
        }
      }

      zoneButtons.forEach(function (button) {
        button.addEventListener("click", function () {
          zoneButtons.forEach(function (b) { b.classList.remove("active"); });
          button.classList.add("active");
          selectedZone = button.dataset.zone;
          zoneMobileSelect.value = selectedZone;
          loadCommunities(selectedZone);
          loadZoneStats(selectedZone);
        });
      });

      zoneMobileSelect.addEventListener("change", function () {
        const zone = zoneMobileSelect.value;
        selectedZone = zone;
        zoneButtons.forEach(function (b) {
          b.classList.toggle("active", b.dataset.zone === zone);
        });
        loadCommunities(selectedZone);
        loadZoneStats(selectedZone);
      });

      zoneMobileSelectBottom.addEventListener("change", function () {
        zoneMobileSelect.value = zoneMobileSelectBottom.value;
        zoneMobileSelect.dispatchEvent(new Event("change", { bubbles: true }));
      });

      communitySearch.addEventListener("change", function () {
        const value = communitySearch.value.trim();
        if (!value) {
          loadZoneStats(selectedZone);
          return;
        }

        const communities = communitiesByZone[selectedZone] || [];
        const match = communities.find(function (community) {
          return normalizeForCompare(community) === normalizeForCompare(value);
        });

        if (!match) {
          loadZoneStats(selectedZone);
          return;
        }

        loadCommunityStats(selectedZone, match);
      });

      communitySearchBottom.addEventListener("change", function () {
        communitySearch.value = communitySearchBottom.value;
        communitySearch.dispatchEvent(new Event("change", { bubbles: true }));
      });

      communitySelect.addEventListener("change", function () {
        const value = communitySelect.value.trim();
        if (!value) {
          loadZoneStats(selectedZone);
          return;
        }
        loadCommunityStats(selectedZone, value);
      });

      communitySelectBottom.addEventListener("change", function () {
        communitySelect.value = communitySelectBottom.value;
        communitySelect.dispatchEvent(new Event("change", { bubbles: true }));
      });

      function reloadForCurrentSelection() {
        const communityValue = (window.matchMedia("(max-width: 960px)").matches ? communitySelect.value : communitySearch.value).trim();
        if (!communityValue) {
          loadZoneStats(selectedZone);
          return;
        }
        loadCommunityStats(selectedZone, communityValue);
      }

      fromDate.addEventListener("change", reloadForCurrentSelection);
      toDate.addEventListener("change", reloadForCurrentSelection);

      fromDateBottom.addEventListener("change", function () {
        fromDate.value = fromDateBottom.value;
        fromDate.dispatchEvent(new Event("change", { bubbles: true }));
      });

      toDateBottom.addEventListener("change", function () {
        toDate.value = toDateBottom.value;
        toDate.dispatchEvent(new Event("change", { bubbles: true }));
      });

      loadCommunities("General");
      loadZoneStats("General");
    })();
  </script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/flowbite/2.5.2/flowbite.min.js"></script>
</body>
</html>`)
		return err
	})
}
