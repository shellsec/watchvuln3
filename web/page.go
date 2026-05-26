package web

const dashboardHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>漏洞情报看板 · WatchVuln</title>
  <style>
    :root { --bg:#0f1419; --card:#1a2332; --border:#2d3a4f; --text:#e7ecf3; --muted:#8b9cb3;
      --crit:#ff4d4f; --high:#ff7a45; --med:#faad14; --low:#52c41a; --accent:#3b82f6; }
    * { box-sizing: border-box; }
    body { margin:0; font-family: "Segoe UI", system-ui, sans-serif; background:var(--bg); color:var(--text); }
    header { padding:20px 24px; border-bottom:1px solid var(--border); display:flex; flex-wrap:wrap; gap:16px; align-items:center; justify-content:space-between; }
    .brand { display:flex; align-items:center; gap:12px; flex-wrap:wrap; }
    h1 { margin:0; font-size:1.35rem; font-weight:600; }
    .repo-link { color:var(--muted); font-size:.85rem; text-decoration:none; border:1px solid var(--border); border-radius:6px; padding:4px 10px; }
    .repo-link:hover { color:var(--accent); border-color:var(--accent); }
    .stats { display:flex; gap:12px; flex-wrap:wrap; }
    .stat { background:var(--card); border:1px solid var(--border); border-radius:8px; padding:10px 14px; min-width:88px; }
    .stat b { display:block; font-size:1.25rem; }
    .stat span { color:var(--muted); font-size:.75rem; }
    main { padding:20px 24px; max-width:1400px; margin:0 auto; }
    .toolbar { display:flex; flex-wrap:wrap; gap:10px; margin-bottom:16px; }
    input, select, button { background:var(--card); border:1px solid var(--border); color:var(--text); border-radius:6px; padding:8px 12px; font-size:.9rem; }
    button { cursor:pointer; background:var(--accent); border-color:var(--accent); }
    button:hover { filter:brightness(1.1); }
    table { width:100%; border-collapse:collapse; background:var(--card); border-radius:10px; overflow:hidden; border:1px solid var(--border); }
    th, td { padding:10px 12px; text-align:left; border-bottom:1px solid var(--border); font-size:.88rem; }
    th { color:var(--muted); font-weight:500; background:#151d28; }
    tr:hover { background:#1f2a3d; cursor:pointer; }
    .sev { display:inline-block; padding:2px 8px; border-radius:4px; font-size:.75rem; font-weight:600; }
    .sev-严重 { background:rgba(255,77,79,.2); color:var(--crit); }
    .sev-高危 { background:rgba(255,122,69,.2); color:var(--high); }
    .sev-中危 { background:rgba(250,173,20,.2); color:var(--med); }
    .sev-低危 { background:rgba(82,196,26,.2); color:var(--low); }
    .pager { margin-top:14px; display:flex; gap:8px; align-items:center; color:var(--muted); }
    .modal { display:none; position:fixed; inset:0; background:rgba(0,0,0,.65); z-index:10; align-items:center; justify-content:center; padding:20px; }
    .modal.open { display:flex; }
    .modal-box { background:var(--card); border:1px solid var(--border); border-radius:12px; max-width:720px; width:100%; max-height:85vh; overflow:auto; padding:20px; }
    .modal-box h2 { margin:0 0 12px; font-size:1.1rem; }
    .modal-box pre { white-space:pre-wrap; word-break:break-word; color:var(--muted); font-size:.85rem; }
    .links a { color:var(--accent); margin-right:10px; font-size:.85rem; }
    .empty { text-align:center; padding:40px; color:var(--muted); }
    .hint { margin:-8px 0 14px; color:var(--muted); font-size:.82rem; }
  </style>
</head>
<body>
  <header>
    <div class="brand">
      <h1>漏洞情报看板</h1>
      <a class="repo-link" href="https://github.com/shellsec/watchvuln3" target="_blank" rel="noopener noreferrer">GitHub</a>
    </div>
    <div class="stats" id="stats"></div>
  </header>
  <main>
    <div class="toolbar">
      <input type="search" id="q" placeholder="搜索标题 / CVE / 描述…" style="min-width:220px" />
      <select id="severity"><option value="">全部等级</option><option>严重</option><option>高危</option><option>中危</option><option>低危</option></select>
      <select id="source"><option value="">全部来源</option></select>
      <select id="sort">
        <option value="disclosure">按披露日期</option>
        <option value="update">按入库更新</option>
      </select>
      <button type="button" id="searchBtn">查询</button>
    </div>
    <p class="hint" id="sortHint">按披露日期排序，最近公开的 CVE 排在前面。</p>
    <table>
      <thead><tr><th>等级</th><th>标题</th><th>CVE</th><th>披露日期</th><th>入库更新</th></tr></thead>
      <tbody id="tbody"></tbody>
    </table>
    <div class="pager">
      <button type="button" id="prevBtn">上一页</button>
      <span id="pageInfo"></span>
      <button type="button" id="nextBtn">下一页</button>
    </div>
  </main>
  <div class="modal" id="modal"><div class="modal-box" id="modalBody"></div></div>
  <script>
    let page = 1, total = 0, limit = 30, items = [];
    const sortHints = {
      disclosure: '按披露日期排序，最近公开的 CVE 排在前面。',
      update: '按入库更新时间排序，最近被程序同步或变更的记录排在前面。'
    };
    const sevClass = s => 'sev sev-' + (s || '');
    function updateSortHint() {
      const sort = document.getElementById('sort').value;
      document.getElementById('sortHint').textContent = sortHints[sort] || sortHints.disclosure;
    }
    async function loadStats() {
      const r = await fetch('/api/stats'); const d = await r.json();
      const el = document.getElementById('stats');
      el.innerHTML = '<div class="stat"><b>' + d.total + '</b><span>漏洞总数</span></div>';
      for (const [k,v] of Object.entries(d.by_severity || {})) {
        el.innerHTML += '<div class="stat"><b>' + v + '</b><span>' + k + '</span></div>';
      }
    }
    async function loadSources() {
      const r = await fetch('/api/sources'); const list = await r.json();
      const sel = document.getElementById('source');
      list.forEach(s => {
        const o = document.createElement('option');
        o.value = s.id; o.textContent = s.name;
        sel.appendChild(o);
      });
    }
    async function loadVulns() {
      const params = new URLSearchParams({ page, limit, q: document.getElementById('q').value,
        severity: document.getElementById('severity').value, source: document.getElementById('source').value,
        sort: document.getElementById('sort').value });
      const r = await fetch('/api/vulns?' + params); const d = await r.json();
      total = d.total; items = d.items || []; page = d.page;
      const tb = document.getElementById('tbody');
      if (!items.length) { tb.innerHTML = '<tr><td colspan="5" class="empty">暂无数据</td></tr>'; }
      else {
        tb.innerHTML = items.map((v,i) => '<tr data-i="'+i+'"><td><span class="'+sevClass(v.severity)+'">'+(v.severity||'-')+'</span></td><td>'+esc(v.title)+'</td><td>'+esc(v.cve||'-')+'</td><td>'+esc(v.disclosure||'-')+'</td><td>'+esc((v.update_time||'').slice(0,10))+'</td></tr>').join('');
        tb.querySelectorAll('tr[data-i]').forEach(tr => tr.onclick = () => showDetail(items[+tr.dataset.i]));
      }
      document.getElementById('pageInfo').textContent = '第 '+page+' 页 / 共 '+total+' 条';
      document.getElementById('prevBtn').disabled = page <= 1;
      document.getElementById('nextBtn').disabled = page * limit >= total;
    }
    function esc(s) { const d=document.createElement('div'); d.textContent=s; return d.innerHTML; }
    function showDetail(v) {
      const refs = (v.references||[]).map(u=>'<a href="'+u+'" target="_blank" rel="noopener">'+u+'</a>').join('<br>');
      const tags = (v.tags||[]).join(', ') || '-';
      document.getElementById('modalBody').innerHTML = '<h2>'+esc(v.title)+'</h2><p class="'+sevClass(v.severity)+'">'+esc(v.severity)+' · '+esc(v.cve||'无CVE')+'</p><p><a href="'+esc(v.from)+'" target="_blank" rel="noopener">来源链接</a></p><h3>描述</h3><pre>'+esc(v.description||'')+'</pre><h3>标签</h3><pre>'+esc(tags)+'</pre><h3>修复建议</h3><pre>'+esc(v.solutions||'')+'</pre><h3>参考</h3><div class="links">'+refs+'</div>';
      document.getElementById('modal').classList.add('open');
    }
    document.getElementById('modal').onclick = e => { if (e.target.id==='modal') e.target.classList.remove('open'); };
    document.getElementById('searchBtn').onclick = () => { page=1; loadVulns(); };
    document.getElementById('sort').onchange = () => { updateSortHint(); page=1; loadVulns(); };
    document.getElementById('q').onkeydown = e => { if (e.key==='Enter') { page=1; loadVulns(); }};
    document.getElementById('prevBtn').onclick = () => { if (page>1) { page--; loadVulns(); }};
    document.getElementById('nextBtn').onclick = () => { if (page*limit<total) { page++; loadVulns(); }};
    updateSortHint(); loadStats(); loadSources(); loadVulns();
  </script>
</body>
</html>`
