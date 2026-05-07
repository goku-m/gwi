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
    .zone-mobile-picker {
      display: none;
      margin: 0 10px 10px;
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
    .zone-btn:focus-visible {
      outline: 2px solid #78c9bf;
      outline-offset: 1px;
    }
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
      padding: 18px 20px;
      box-shadow: var(--shadow);
      margin-bottom: 18px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
    }
    .header-left { min-width: 0; }
    .title {
      margin: 0;
      font-size: clamp(1.25rem, 1rem + 1vw, 1.8rem);
      line-height: 1.18;
      font-weight: 800;
      letter-spacing: -0.01em;
    }
    .subtitle {
      margin-top: 7px;
      color: var(--muted);
      font-size: 0.93rem;
      min-height: 1.2rem;
    }
    .header-metrics {
      margin-top: 12px;
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
    }
    .header-metric {
      background: var(--surface-soft);
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 8px 10px;
      min-width: 138px;
    }
    .header-metric-label {
      display: block;
      color: var(--muted);
      font-size: 0.74rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.04em;
      margin-bottom: 4px;
    }
    .header-metric-value {
      margin: 0;
      font-size: 1.05rem;
      font-weight: 800;
      color: #0d3a36;
    }
    .metric-increment {
      color: #15803d;
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
      border-color: #8ad2ca;
      box-shadow: 0 0 0 3px rgba(138, 210, 202, 0.25);
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
      gap: 14px;
      grid-template-columns: repeat(2, minmax(190px, 1fr));
    }
    .card {
      background: linear-gradient(180deg, #ffffff 0%, #f9fcfb 100%);
      border: 1px solid var(--border);
      border-radius: var(--radius);
      padding: 16px;
      box-shadow: var(--shadow);
      min-height: 116px;
      transition: transform 140ms ease, box-shadow 140ms ease;
    }
    .card:hover {
      transform: translateY(-2px);
      box-shadow: 0 16px 30px rgba(16, 33, 39, 0.11);
    }
    .card-label {
      margin: 0 0 10px;
      color: var(--muted);
      font-size: 0.84rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.04em;
    }
    .card-value {
      margin: 0;
      font-size: clamp(1.14rem, 1rem + 0.8vw, 1.7rem);
      font-weight: 800;
      letter-spacing: -0.01em;
      color: #0d3a36;
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
      color: #0a5d56;
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
      .sidebar {
        height: auto;
        position: static;
        border-right: 0;
        border-bottom: 1px solid var(--border);
      }
      .zone-list {
        display: none;
      }
      .zone-mobile-picker { display: block; }
      .cards { grid-template-columns: repeat(2, minmax(150px, 1fr)); }
      .header { align-items: flex-start; flex-direction: column; }
      .community-picker { width: 100%; min-width: 0; }
      .community-input { display: none; }
      .community-select { display: block; }
      .filters-row { grid-template-columns: 1fr 1fr; }
      .charts { grid-template-columns: 1fr; }
    }
    @media (max-width: 560px) {
      .main { padding: 14px; }
      .cards { grid-template-columns: repeat(2, minmax(130px, 1fr)); }
      .title { font-size: 1.2rem; }
      .brand-logo { width: 108px; }
    }
  </style>
</head>
<body>
  <div class="layout">
    <aside class="sidebar">
      <div class="brand">
        <img class="brand-logo" src="/static/images/tk.png" alt="TEBMA KANDU logo" />
      </div>
      <div class="zone-mobile-picker">
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
    </aside>
    <main class="main" id="dashboardMain">
      <section class="header">
        <div class="header-left">
          <h1 class="title" id="selectedZoneTitle">General</h1>
          <div class="subtitle" id="subtitleText">Zone-level Stats</div>
          <div class="header-metrics">
            <div class="header-metric">
              <span class="header-metric-label">Farmers</span>
              <p class="header-metric-value" id="totalFarmers">0</p>
            </div>
            <div class="header-metric">
              <span class="header-metric-label">Communities</span>
              <p class="header-metric-value" id="totalCommunities">0</p>
            </div>
            <div class="header-metric">
              <span class="header-metric-label">Syncs Today</span>
              <p class="header-metric-value" id="dailySyncs">0</p>
            </div>
          </div>
        </div>
        <div class="community-picker">
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

      <section class="cards">        <article class="card"><p class="card-label">Total Nuts (kg)</p><p class="card-value" id="totalKgBrought">0</p></article>
        <article class="card"><p class="card-label">Nuts Value (GH₵)</p><p class="card-value" id="totalAmount">0</p></article>
        <article class="card"><p class="card-label">Prefinance Given (GH₵)</p><p class="card-value" id="totalPrefinance">0</p></article>
        <article class="card"><p class="card-label">Recovery Balance (GH₵)</p><p class="card-value" id="totalBalance">0</p></article>
      </section>
      <div class="loading-indicator" role="status" aria-live="polite">
        <span class="spinner" aria-hidden="true"></span>
        <span>Fetching latest data...</span>
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
      <div class="error" id="errorBox" role="status" aria-live="polite"></div>
    </main>
  </div>

  <script>
    (function () {
      const zoneButtons = Array.from(document.querySelectorAll(".zone-btn"));
      const selectedZoneTitle = document.getElementById("selectedZoneTitle");
      const subtitleText = document.getElementById("subtitleText");
      const errorBox = document.getElementById("errorBox");
      const communitySearch = document.getElementById("communitySearch");
      const communitySelect = document.getElementById("communitySelect");
      const zoneMobileSelect = document.getElementById("zoneMobileSelect");
      const communityOptions = document.getElementById("communityOptions");
      const fromDate = document.getElementById("fromDate");
      const toDate = document.getElementById("toDate");
      const dashboardMain = document.getElementById("dashboardMain");

      const totalFarmers = document.getElementById("totalFarmers");
      const totalCommunities = document.getElementById("totalCommunities");
      const dailySyncs = document.getElementById("dailySyncs");
      const totalKgBrought = document.getElementById("totalKgBrought");
      const totalAmount = document.getElementById("totalAmount");
      const totalPrefinance = document.getElementById("totalPrefinance");
      const totalBalance = document.getElementById("totalBalance");
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

      function setLoading() {
        dashboardMain.classList.add("loading");
        totalFarmers.textContent = "...";
        totalCommunities.textContent = "...";
        dailySyncs.textContent = "...";
        totalKgBrought.textContent = "...";
        totalAmount.textContent = "...";
        totalPrefinance.textContent = "...";
        totalBalance.textContent = "...";
      }

      function setStats(data) {
        dashboardMain.classList.remove("loading");
        latestTotalFarmers = Number(data.totalFarmers || 0);
        totalFarmers.innerHTML = formatNumber(latestTotalFarmers, 0) + " <span class=\"metric-increment\">↑" + formatNumber(latestNewFarmers, 0) + "</span>";
        totalCommunities.textContent = formatNumber(data.totalCommunities, 0);
        dailySyncs.textContent = formatNumber(data.dailySyncs, 0);
        totalKgBrought.textContent = formatNumber(data.totalKgBrought, 2);
        totalAmount.textContent = formatNumber(data.totalAmount, 2);
        totalPrefinance.textContent = formatNumber(data.totalPrefinance, 2);
        totalBalance.textContent = formatNumber(data.totalBalance, 2);
        renderCharts(data);
      }

      function setNewFarmers(data) {
        const count = (data && typeof data.newFarmers !== "undefined") ? data.newFarmers : 0;
        latestNewFarmers = Number(count || 0);
        totalFarmers.innerHTML = formatNumber(latestTotalFarmers, 0) + " <span class=\"metric-increment\">↑" + formatNumber(latestNewFarmers, 0) + "</span>";
      }

      function renderCharts(data) {
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
          recoveryChart = new ApexCharts(document.querySelector("#recoveryChart"), recoveryOptions);
          recoveryChart.render();
        } else {
          recoveryChart.updateOptions(recoveryOptions);
        }

        if (!financeChart) {
          financeChart = new ApexCharts(document.querySelector("#financeChart"), financeOptions);
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
        communityOptions.innerHTML = "";

        if (zone === "General") {
          communitySearch.disabled = true;
          communitySelect.disabled = true;
          communitySearch.placeholder = "All communities";
          communitySelect.innerHTML = "<option value=\"\">All communities</option>";
          return;
        }

        communitySearch.disabled = false;
        communitySelect.disabled = false;
        communitySearch.placeholder = "Search communities";

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

        const defaultOption = document.createElement("option");
        defaultOption.value = "";
        defaultOption.textContent = "All communities";
        communitySelect.appendChild(defaultOption);

        communities.forEach(function (community) {
          const option = document.createElement("option");
          option.value = community;
          communityOptions.appendChild(option);

          const selectOption = document.createElement("option");
          selectOption.value = community;
          selectOption.textContent = community;
          communitySelect.appendChild(selectOption);
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

      communitySelect.addEventListener("change", function () {
        const value = communitySelect.value.trim();
        if (!value) {
          loadZoneStats(selectedZone);
          return;
        }
        loadCommunityStats(selectedZone, value);
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
