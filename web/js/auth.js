(async function authGuard() {
  document.documentElement.style.visibility = 'hidden';

  try {
    const res = await fetch('/api/auth/me', {
      method: 'GET',
      credentials: 'include',
      headers: { Accept: 'application/json' },
    });

    if (res.status === 401 || res.status === 403) {
      window.location.replace('/index.html');
      return;
    }

    if (res.ok) {
      const user = await res.json();
      window.currentUser = user;
    }
  } catch {
    // Бэкенд недоступен — показываем страницу (dev-режим)
  }

  document.documentElement.style.visibility = '';
})();

async function logout() {
  try {
    await fetch('/api/auth/logout', { method: 'POST', credentials: 'include' });
  } catch { /* игнорируем */ }
  window.location.replace('/index.html');
}
