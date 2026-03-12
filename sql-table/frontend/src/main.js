import './style.css';
import './app.css';

window.connect = async function() {
  const host = document.getElementById('host').value;
  const port = document.getElementById('port').value;
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;
  const statusEl = document.getElementById('connection-status');

  try {
    await window.go.main.App.Connect({ Host: host, Port: port, Username: username, Password: password });
    statusEl.innerText = 'Connected';
    statusEl.className = 'success';
    
    const dbs = await window.go.main.App.ListDatabases();
    console.log('Databases:', dbs);
  } catch (err) {
    statusEl.innerText = 'Error: ' + err;
    statusEl.className = 'error';
  }
};

window.executeQuery = async function() {
  const query = document.getElementById('query').value;
  const resultsEl = document.getElementById('results');

  try {
    const result = await window.go.main.App.ExecuteQuery(query);
    
    if (!result || !result.Columns || result.Columns.length === 0) {
      resultsEl.innerHTML = '<p>Query executed successfully (no results)</p>';
      return;
    }

    let html = '<table><thead><tr>';
    for (const col of result.Columns) {
      html += `<th>${col}</th>`;
    }
    html += '</tr></thead><tbody>';
    
    for (const row of result.Rows) {
      html += '<tr>';
      for (const cell of row) {
        html += `<td>${cell === null ? 'NULL' : cell}</td>`;
      }
      html += '</tr>';
    }
    html += '</tbody></table>';
    
    resultsEl.innerHTML = html;
  } catch (err) {
    resultsEl.innerHTML = `<p class="error">Error: ${err}</p>`;
  }
};
