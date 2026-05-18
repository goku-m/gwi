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
  <style>
    * {
      box-sizing: border-box;
    }
    body {
      margin: 0;
      font-family: "Segoe UI", sans-serif;
      background: #f5f7f9;
      color: #102127;
    }
    .wrap {
      max-width: 960px;
      margin: 0 auto;
      padding: 24px 16px;
    }
    .card {
      background: #fff;
      border: 1px solid #d8e5e3;
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
      border: 1px solid #d8e5e3;
      border-radius: 10px;
      padding: 10px 12px;
      font-size: 0.95rem;
      color: #102127;
      background: #f8fbfa;
    }
    .date-input:focus {
      outline: none;
      border-color: #8ad2ca;
      box-shadow: 0 0 0 3px rgba(138, 210, 202, 0.25);
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
      border: 1px solid #d8e5e3;
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
      color: #15803d;
      font-weight: 600;
    }
    a {
      color: #0f766e;
      text-decoration: none;
      font-weight: 600;
    }
    a:hover {
      text-decoration: underline;
    }
    @media (max-width: 640px) {
      .wrap {
        padding: 16px 12px;
      }
      .card {
        padding: 16px;
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
      <a href="/">Back to dashboard</a>
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
