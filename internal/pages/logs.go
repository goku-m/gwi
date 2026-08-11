package pages

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func Logs() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Logs</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Manrope:wght@400;500;600;700;800&display=swap" rel="stylesheet">
  <script src="https://cdn.tailwindcss.com"></script>
  <link href="https://cdnjs.cloudflare.com/ajax/libs/flowbite/2.5.2/flowbite.min.css" rel="stylesheet" />
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
    * {
      box-sizing: border-box;
    }
    body {
      margin: 0;
      font-family: "Manrope", "Segoe UI", sans-serif;
      color: var(--text);
      background: var(--bg);
    }
    .wrap {
      max-width: 960px;
      margin: 0 auto;
      padding: 24px 16px;
    }
    .header {
      background: var(--surface);
      border: 1px solid var(--border);
      border-radius: var(--radius);
      box-shadow: var(--shadow);
      padding: 7px 10px;
      margin-bottom: 8px;
    }
    .header-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      flex-wrap: wrap;
    }
    .header-title {
      margin: 0;
      font-size: 1rem;
      font-weight: 800;
      letter-spacing: -0.01em;
    }
    .header-subtitle {
      margin: 2px 0 0;
      color: var(--muted);
      font-size: 0.76rem;
    }
    .header-actions {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
    }
    .nav-link {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      border: 1px solid var(--border);
      background: var(--surface);
      border-radius: 12px;
      padding: 7px 9px;
      color: var(--text);
      text-decoration: none;
      font-size: 0.86rem;
      font-weight: 600;
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
    .header-actions.mobile {
      display: none;
      width: 100%;
    }
    .header-actions.mobile .nav-link {
      flex: 1 1 0;
    }
    .main-mobile-brand {
      display: none;
    }
    .main-mobile-actions {
      display: none;
    }
    .card {
      background: #fff;
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 20px;
    }
    .toolbar {
      display: flex;
      align-items: end;
      justify-content: space-between;
      gap: 12px;
      margin-bottom: 16px;
      flex-wrap: wrap;
    }
    .date-group {
      display: grid;
      gap: 6px;
      min-width: 240px;
      flex: 1 1 280px;
    }
    .date-label {
      font-size: 0.85rem;
      font-weight: 700;
      color: #56707d;
    }
    .date-input {
      width: 100%;
      max-width: 100%;
      border: 1px solid var(--border);
      border-radius: 10px;
      padding: 10px 12px;
      font-size: 0.95rem;
      color: #102127;
      background: #f8fbfa;
    }
    .date-input:focus {
      outline: none;
      border-color: var(--accent);
      box-shadow: 0 0 0 3px rgba(22, 101, 52, 0.18);
    }
    h1 {
      margin: 0 0 10px;
      font-size: 1.5rem;
    }
    p {
      margin: 0 0 14px;
      color: #56707d;
    }
    .status {
      margin: 0 0 12px;
      color: #56707d;
      font-size: 0.92rem;
    }
    .names-list {
      list-style: none;
      margin: 0 0 16px;
      padding: 0;
      display: grid;
      gap: 8px;
    }
    .names-list li {
      border: 1px solid var(--border);
      border-radius: 10px;
      background: #f8fbfa;
      padding: 10px 12px;
      font-weight: 400;
      line-height: 1.45;
    }
    .log-strong {
      font-weight: 700;
    }
    .log-datetime {
      color: var(--accent);
      font-weight: 600;
    }
    a {
      color: var(--accent);
      text-decoration: none;
      font-weight: 600;
    }
    a:hover {
      text-decoration: underline;
    }
    @media (pointer: coarse) and (hover: none) {
      .wrap {
        padding: 16px 12px;
      }
      .card {
        padding: 16px;
      }
      .header {
        padding: 7px 8px;
      }
      .header-actions.desktop {
        display: none;
      }
      .header-actions.mobile {
        display: flex;
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
      .toolbar {
        align-items: stretch;
      }
      .date-group {
        min-width: 0;
        flex: 1 1 100%;
      }
    }
  </style>
</head>
<body>
  <main class="wrap">
    <div class="main-mobile-brand">
      <img class="brand-logo" src="/static/images/tk.png" alt="TEBMA KANDU logo" />
    </div>
    <div class="main-mobile-actions">
      <a class="nav-link" href="/">Home</a>
      <a class="nav-link active" href="/logs">Daily Logs</a>
      <a class="nav-link" href="/recovery">Recovery</a>
    </div>
    <header class="header">
      <div class="header-row">
        <div>
          <h1 class="header-title">Logs</h1>
          <p class="header-subtitle">Browse daily farmer activity.</p>
        </div>
        <nav class="header-actions desktop" aria-label="Primary">
          <a class="nav-link" href="/">Home</a>
          <a class="nav-link active" href="/logs">Daily Logs</a>
          <a class="nav-link" href="/recovery">Recovery</a>
        </nav>
      </div>
  
    </header>
    <section class="card">
      <div class="toolbar">
        <div class="date-group">
          <label class="date-label" for="logDate">Select Date</label>
          <input id="logDate" class="date-input" type="date" />
        </div>
      </div>
      <h1>Logs</h1>
      <p class="status" id="logStatus">Choose a date to load logs.</p>
      <ul class="names-list" id="namesList"></ul>
    </section>
  </main>
  <script>
    (function () {
      const logDate = document.getElementById("logDate");
      const logStatus = document.getElementById("logStatus");
      const namesList = document.getElementById("namesList");

      function formatDateForInput(date) {
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, "0");
        const day = String(date.getDate()).padStart(2, "0");
        return year + "-" + month + "-" + day;
      }

      function renderLogs(logs) {
        namesList.innerHTML = "";
        function toPascalText(value) {
          return String(value || "")
            .toLowerCase()
            .split(/\s+/)
            .filter(Boolean)
            .map(function (part) {
              return part.charAt(0).toUpperCase() + part.slice(1);
            })
            .join(" ");
        }
        logs.forEach(function (log) {
          const item = document.createElement("li");
          const communityNames = String(log.communityNames || "").trim();
          const communitiesText = communityNames ? toPascalText(communityNames) : "N/A";
          const dateTime = "<span class=\"log-datetime\">" + log.date + ": " + log.time + "</span>";
          if (log.action === "updated") {
            const amount = Number(log.amount || 0).toLocaleString(undefined, { maximumFractionDigits: 2 });
            item.innerHTML = dateTime + " <span class=\"log-strong\">" + toPascalText(log.updatedBy) + "</span> from <span class=\"log-strong\">" + toPascalText(log.zoneName) + "</span> issued <span class=\"log-strong\">GHc</span> <span class=\"log-strong\">" + amount + "</span> in <span class=\"log-strong\">Pre-Finance</span> in communities: <span class=\"log-strong\">" + communitiesText + "</span>.";
          } else if (log.action === "weighed") {
            const weight = Number(log.weightKg || 0).toLocaleString(undefined, { maximumFractionDigits: 2 });
            const amount = Number(log.amount || 0).toLocaleString(undefined, { maximumFractionDigits: 2 });
            item.innerHTML = dateTime + " <span class=\"log-strong\">" + toPascalText(log.updatedBy) + "</span> from <span class=\"log-strong\">" + toPascalText(log.zoneName) + "</span> weighed <span class=\"log-strong\">" + weight + "kg</span> of nuts at a total value of GHc <span class=\"log-strong\">" + amount + "</span> in communities: <span class=\"log-strong\">" + communitiesText + "</span>.";
          } else {
            const countLabel = Number(log.count) === 1 ? "farmer" : "farmers";
            item.innerHTML = dateTime + " <span class=\"log-strong\">" + toPascalText(log.createdBy) + "</span> from <span class=\"log-strong\">" + toPascalText(log.zoneName) + "</span> added <span class=\"log-strong\">" + log.count + "</span> " + countLabel + " to communities: <span class=\"log-strong\">" + communitiesText + "</span>.";
          }
          namesList.appendChild(item);
        });
      }

      async function loadLogs(dateValue) {
        if (!dateValue) {
          logStatus.textContent = "Choose a date to load logs.";
          namesList.innerHTML = "";
          return;
        }

        logStatus.textContent = "Loading...";
        namesList.innerHTML = "";

        try {
          const response = await fetch("/api/farmers/logs?date=" + encodeURIComponent(dateValue));
          if (!response.ok) {
            throw new Error("Request failed");
          }

          const payload = await response.json();
          const logs = Array.isArray(payload.logs) ? payload.logs : [];

          if (logs.length === 0) {
            logStatus.textContent = "No logs found for " + dateValue + ".";
            return;
          }

          logStatus.textContent = "Logs for " + dateValue + ":";
          renderLogs(logs);
        } catch (err) {
          logStatus.textContent = "Could not load logs for " + dateValue + ".";
        }
      }

      const today = formatDateForInput(new Date());
      logDate.value = today;
      loadLogs(today);

      logDate.addEventListener("change", function () {
        loadLogs(logDate.value);
      });
    })();
  </script>
</body>
</html>`)
		return err
	})
}
