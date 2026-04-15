/**
 * Базовый HTTP-клиент для работы с /api/*
 * Автоматически редиректит на /index.html при 401.
 */
const api = (() => {
  async function request(method, path, body) {
    const opts = {
      method,
      credentials: 'include',
      headers: { Accept: 'application/json' },
    };
    if (body !== undefined) {
      opts.headers['Content-Type'] = 'application/json';
      opts.body = JSON.stringify(body);
    }

    const res = await fetch('/api' + path, opts);

    if (res.status === 401) {
      window.location.replace('/');
      throw new Error('Unauthorized');
    }
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      const err = new Error(data.message || `HTTP ${res.status}`);
      err.data = data;
      err.status = res.status;
      throw err;
    }
    if (res.status === 204) return null;
    return res.json();
  }

  return {
    get:    (path)       => request('GET',    path),
    post:   (path, body) => request('POST',   path, body),
    put:    (path, body) => request('PUT',    path, body),
    patch:  (path, body) => request('PATCH',  path, body),
    delete: (path)       => request('DELETE', path),
  };
})();

/** Утилиты форматирования, доступны глобально */
const fmt = {
  money(v) {
    if (v == null) return '₽ 0';
    return '₽ ' + Number(v).toLocaleString('ru-RU', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
  },
  date(v) {
    if (!v) return '—';
    // Date-only strings (YYYY-MM-DD) must be parsed manually to avoid UTC→local shift
    if (typeof v === 'string' && /^\d{4}-\d{2}-\d{2}$/.test(v)) {
      const [y, mo, d] = v.split('-').map(Number);
      return new Date(y, mo - 1, d).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' });
    }
    return new Date(v).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' });
  },
  pct(v) {
    if (v == null) return '0%';
    const arrow = v >= 0 ? '▲' : '▼';
    return arrow + ' ' + Math.abs(v).toFixed(1) + '%';
  },
};
