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
      window.location.replace('/index.html');
      throw new Error('Unauthorized');
    }
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      throw new Error(data.message || `HTTP ${res.status}`);
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
  money: (n) => '₽\u00a0' + Number(n ?? 0).toLocaleString('ru-RU'),
  date:  (d) => d ? new Date(d).toLocaleDateString('ru-RU') : '—',
  pct:   (n) => (n >= 0 ? '▲ ' : '▼ ') + Math.abs(n).toFixed(1) + '%',
};
