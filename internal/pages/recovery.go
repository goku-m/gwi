package pages

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func Recovery() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Recovery</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Manrope:wght@400;500;600;700;800&display=swap" rel="stylesheet">
  <script src="https://cdn.tailwindcss.com"></script>
  <link href="https://cdnjs.cloudflare.com/ajax/libs/flowbite/2.5.2/flowbite.min.css" rel="stylesheet" />
  <script src="https://unpkg.com/htmx.org@2.0.4"></script>
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
    .nav-link.active {
      background: var(--accent-soft);
      border-color: var(--accent);
      color: var(--accent);
      font-weight: 800;
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
      display: grid;
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
      color: var(--accent);
      font-size: 0.76rem;
      font-weight: 700;
    }
    .community-picker {
      display: flex;
      align-items: end;
      justify-content: flex-end;
      gap: 10px;
      flex-wrap: wrap;
      min-width: 220px;
    }
    .community-zone-picker,
    .community-community-picker {
      display: grid;
      gap: 6px;
      min-width: 180px;
      flex: 1 1 180px;
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
      border-color: var(--accent);
      box-shadow: 0 0 0 3px rgba(22, 101, 52, 0.18);
    }
    .recovery-dashboard {
      display: grid;
      gap: 12px;
      margin-bottom: 14px;
    }
    .dashboard-head {
      display: flex;
      align-items: flex-end;
      justify-content: space-between;
      gap: 10px;
      flex-wrap: wrap;
    }
    .dashboard-kicker {
      margin: 0;
      color: var(--muted);
      font-size: 0.78rem;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    .dashboard-title {
      margin: 3px 0 0;
      font-size: 1rem;
      font-weight: 800;
      letter-spacing: -0.01em;
    }
    .dashboard-note {
      margin: 0;
      color: var(--muted);
      font-size: 0.86rem;
      font-weight: 600;
    }
    .dashboard-grid {
      display: grid;
      gap: 10px;
      grid-template-columns: repeat(5, minmax(0, 1fr));
    }
    .dashboard-card {
      background: linear-gradient(180deg, #ffffff 0%, #f9fcfb 100%);
      border: 1px solid var(--border);
      border-radius: var(--radius);
      box-shadow: var(--shadow);
      padding: 12px;
      min-height: 108px;
      display: grid;
      gap: 6px;
    }
    .dashboard-label {
      color: var(--muted);
      font-size: 0.76rem;
      font-weight: 800;
      letter-spacing: 0.06em;
      text-transform: uppercase;
    }
    .dashboard-value {
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
    .dashboard-value.good { color: var(--accent); }
    .dashboard-value.warn { color: #b45309; }
    .dashboard-value.accent { color: var(--accent); }
    .dashboard-subvalue {
      color: var(--muted);
      font-size: 0.82rem;
      line-height: 1.25;
    }
    .dashboard-progress {
      height: 9px;
      background: #e7efed;
      border-radius: 999px;
      overflow: hidden;
      margin-top: 2px;
    }
    .dashboard-progress span {
      display: block;
      height: 100%;
      width: 0%;
      border-radius: inherit;
      background: linear-gradient(90deg, #166534 0%, #22c55e 100%);
      transition: width 220ms ease;
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
    @media (pointer: coarse) and (hover: none) {
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
        justify-content: stretch;
      }
      .community-zone-picker,
      .community-community-picker {
        min-width: 0;
      }
      .dashboard-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
    }
    @media (max-width: 560px) {
      .main { padding: 14px; }
      .title { font-size: 0.98rem; }
      .main-mobile-brand .brand-logo { width: 84px; }
      .dashboard-grid { grid-template-columns: 1fr; }
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
          <a class="nav-link active" href="/recovery">Recovery</a>
        </div>
      </div>
    </aside>
    <main class="main" id="recoveryMain">
      <div class="main-mobile-brand">
        <img class="brand-logo" src="/static/images/tk.png" alt="TEBMA KANDU logo" />
      </div>
      <div class="main-mobile-actions">
        <a class="nav-link" href="/">Home</a>
        <a class="nav-link" href="/logs">Daily Logs</a>
        <a class="nav-link active" href="/recovery">Recovery</a>
      </div>
      <section class="header">
        <div class="header-left">
          <h1 class="title" id="selectedZoneTitle">Recovery</h1>
          <div class="subtitle">
            <span id="subtitleText">Overall zone recovery overview</span>
            <span class="subtitle-syncs" id="subtitleSyncs">Total Syncs: <span id="totalSyncs">0</span></span>
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
          <div class="community-community-picker">
            <label class="community-label" for="communitySelect">Community</label>
            <select id="communitySelect" class="zone-mobile-select" aria-label="Community select">
              <option value="">All communities</option>
            </select>
          </div>
        </div>
      </section>
      <div class="loading-indicator" role="status" aria-live="polite">
        <span class="spinner" aria-hidden="true"></span>
        <span>Fetching latest data...</span>
      </div>
      <section class="recovery-dashboard">
     
        <div class="dashboard-grid">
          <article class="dashboard-card">
            <span class="dashboard-label">Recovery Rate</span>
            <p class="dashboard-value accent" id="recoveryRateValue">0%</p>
           
            <div class="dashboard-progress" aria-hidden="true"><span id="recoveryProgressBar"></span></div>
          </article>
          <article class="dashboard-card">
            <span class="dashboard-label">Amount Payable</span>
            <p class="dashboard-value good" id="amountPayableValue">GH&#8373; 0.00</p>
            <span class="dashboard-subvalue">Amount paid to farmers for nuts</span>
          </article>
          <article class="dashboard-card">
            <span class="dashboard-label">Recovered</span>
            <p class="dashboard-value good" id="recoveredPrefinance">GH&#8373; 0.00</p>
            <span class="dashboard-subvalue">Prefinance already recovered</span>
          </article>
          <article class="dashboard-card">
            <span class="dashboard-label">Outstanding</span>
            <p class="dashboard-value warn" id="outstandingBalance">GH&#8373; 0.00</p>
            <span class="dashboard-subvalue">Amount still to recover</span>
          </article>
          <article class="dashboard-card">
            <span class="dashboard-label">Total Prefinance</span>
            <p class="dashboard-value" id="totalPrefinanceValue">GH&#8373; 0.00</p>
            <span class="dashboard-subvalue">Current zone total</span>
          </article>
        </div>
      </section>
    </main>
  </div>
  <script>
    (function () {
      const zoneButtons = Array.from(document.querySelectorAll(".zone-btn"));
      const selectedZoneTitle = document.getElementById("selectedZoneTitle");
      const subtitleText = document.getElementById("subtitleText");
      const subtitleSyncs = document.getElementById("subtitleSyncs");
      const totalSyncs = document.getElementById("totalSyncs");
      const recoveryRateValue = document.getElementById("recoveryRateValue");
      const recoveryRateSummary = document.getElementById("recoveryRateSummary");
      const recoveryProgressBar = document.getElementById("recoveryProgressBar");
      const amountPayableValue = document.getElementById("amountPayableValue");
      const recoveredPrefinance = document.getElementById("recoveredPrefinance");
      const outstandingBalance = document.getElementById("outstandingBalance");
      const totalPrefinanceValue = document.getElementById("totalPrefinanceValue");
      const recoveryMain = document.getElementById("recoveryMain");
      const zoneMobileSelect = document.getElementById("zoneMobileSelect");
      const communitySelect = document.getElementById("communitySelect");

      let selectedZone = "General";
      let selectedCommunity = "";
      let activeRequestController = null;
      const communitiesByZone = {};

      function formatNumber(value, maxFractionDigits) {
        const n = Number(value || 0);
        return n.toLocaleString(undefined, { maximumFractionDigits: maxFractionDigits });
      }

      function formatCurrency(value) {
        return "<span class=\"currency-symbol\">GH\u20B5</span><span class=\"currency-amount\">" + formatNumber(value, 2) + "</span>";
      }

      function normalizeForCompare(value) {
        return (value || "").trim().toLowerCase();
      }

      function setLoading() {
        recoveryMain.classList.add("loading");
        totalSyncs.textContent = "...";
        if (recoveryRateValue) recoveryRateValue.textContent = "...";
        if (recoveryRateSummary) recoveryRateSummary.textContent = "Loading recovery snapshot...";
        if (recoveryProgressBar) recoveryProgressBar.style.width = "0%";
        if (amountPayableValue) amountPayableValue.textContent = "...";
        if (recoveredPrefinance) recoveredPrefinance.textContent = "...";
        if (outstandingBalance) outstandingBalance.textContent = "...";
        if (totalPrefinanceValue) totalPrefinanceValue.textContent = "...";
      }

      function renderRecoveryDashboard(data) {
        const prefinance = Number(data.totalPrefinance || 0);
        const balance = Number(data.totalBalance || 0);
        const amountPayable = Number(data.totalPayable || 0);
        const recoveredPrefinanceValue = Math.max(0, prefinance - balance);
        const recoveryPercent = prefinance > 0 ? (recoveredPrefinanceValue / prefinance) * 100 : 0;

        if (recoveryRateValue) {
          recoveryRateValue.textContent = recoveryPercent.toFixed(1) + "%";
        }
        if (recoveryRateSummary) {
          recoveryRateSummary.innerHTML = formatCurrency(recoveredPrefinanceValue) + " recovered of " + formatCurrency(prefinance);
        }
        if (recoveryProgressBar) {
          recoveryProgressBar.style.width = Math.max(0, Math.min(100, recoveryPercent)).toFixed(1) + "%";
        }
        if (amountPayableValue) amountPayableValue.innerHTML = formatCurrency(amountPayable);
        if (recoveredPrefinance) recoveredPrefinance.innerHTML = formatCurrency(recoveredPrefinanceValue);
        if (outstandingBalance) outstandingBalance.innerHTML = formatCurrency(balance);
        if (totalPrefinanceValue) totalPrefinanceValue.innerHTML = formatCurrency(prefinance);
      }

      function renderCommunityOptions(communities) {
        communitySelect.innerHTML = "";

        const defaultOption = document.createElement("option");
        defaultOption.value = "";
        defaultOption.textContent = "All communities";
        communitySelect.appendChild(defaultOption);

        communities.forEach(function (community) {
          const option = document.createElement("option");
          option.value = community;
          option.textContent = community;
          communitySelect.appendChild(option);
        });
      }

      async function loadCommunities(zone) {
        selectedCommunity = "";
        communitySelect.value = "";
        communitySelect.disabled = true;
        renderCommunityOptions([]);

        if (zone === "General") {
          communitySelect.disabled = true;
          return;
        }

        communitySelect.disabled = false;

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

      async function loadRecoveryData(zone, community) {
        if (activeRequestController) {
          activeRequestController.abort();
        }

        activeRequestController = new AbortController();
        const controller = activeRequestController;

        selectedZoneTitle.textContent = zone;
        subtitleText.textContent = community ? "Community: " + community : "Overall zone recovery overview";
        subtitleSyncs.style.display = "inline-flex";
        setLoading();

        const route = community
          ? "/api/zones/" + encodeURIComponent(zone) + "/" + encodeURIComponent(community) + "/farmers/stats"
          : zone === "General"
            ? "/api/farmers/overview"
            : "/api/zones/" + encodeURIComponent(zone) + "/farmers/overview";

        try {
          const data = await fetchJson(route, controller.signal);

          if (controller.signal.aborted) {
            return;
          }

          const syncCount = typeof data.totalSyncs !== "undefined" ? data.totalSyncs : data.dailySyncs;
          totalSyncs.textContent = formatNumber(syncCount, 0);
          renderRecoveryDashboard(data);
          recoveryMain.classList.remove("loading");
        } catch (err) {
          if (err.name === "AbortError") {
            return;
          }
          recoveryMain.classList.remove("loading");
        } finally {
          if (activeRequestController === controller) {
            activeRequestController = null;
          }
        }
      }

      async function fetchJson(route, signal) {
        const response = await fetch(route, { signal: signal });
        if (!response.ok) {
          throw new Error("Request failed with status " + response.status);
        }

        return response.json();
      }

      zoneButtons.forEach(function (button) {
        button.addEventListener("click", function () {
          zoneButtons.forEach(function (b) { b.classList.remove("active"); });
          button.classList.add("active");
          selectedZone = button.dataset.zone;
          zoneMobileSelect.value = selectedZone;
          loadCommunities(selectedZone);
          loadRecoveryData(selectedZone, "");
        });
      });

      zoneMobileSelect.addEventListener("change", function () {
        const zone = zoneMobileSelect.value;
        selectedZone = zone;
        zoneButtons.forEach(function (b) {
          b.classList.toggle("active", b.dataset.zone === zone);
        });
        loadCommunities(selectedZone);
        loadRecoveryData(selectedZone, "");
      });

      communitySelect.addEventListener("change", function () {
        selectedCommunity = communitySelect.value.trim();
        if (!selectedCommunity) {
          loadRecoveryData(selectedZone, "");
          return;
        }

        const communities = communitiesByZone[selectedZone] || [];
        const match = communities.find(function (community) {
          return normalizeForCompare(community) === normalizeForCompare(selectedCommunity);
        });

        if (!match) {
          loadRecoveryData(selectedZone, "");
          return;
        }

        loadRecoveryData(selectedZone, match);
      });

      loadCommunities("General");
      loadRecoveryData("General", "");
    })();
  </script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/flowbite/2.5.2/flowbite.min.js"></script>
</body>
</html>`)
		return err
	})
}
