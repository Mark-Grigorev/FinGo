/**
 * Auth guard — подключается на всех защищённых страницах.
 * При загрузке проверяет сессию через GET /api/auth/me.
 * 401 → редирект на /index.html
 */
(async function authGuard() {
  // Скрываем контент до проверки (избегаем flash незащищённого контента)
  document.documentElement.style.visibility = 'hidden';

  try {
    const res = await fetch('/api/auth/me', {
      method: 'GET',
      credentials: 'include',           // отправляем httpOnly cookie
      headers: { 'Accept': 'application/json' },
    });

    if (res.status === 401 || res.status === 403) {
      window.location.replace('/index.html');
      return;
    }

    // Можно использовать данные пользователя прямо здесь
    if (res.ok) {
      const user = await res.json();
      // Сохраняем в window для использования на странице
      window.currentUser = user;
    }
  } catch {
    // Нет соединения с бэкендом — разрешаем просмотр (dev-режим)
    // В продакшне можно раскомментировать редирект:
    // window.location.replace('/index.html');
  }

  document.documentElement.style.visibility = '';
})();


/**
 * Выход из аккаунта.
 * Вызывается кнопкой "Выйти" в сайдбаре.
 */
async function logout() {
  try {
    await fetch('/api/auth/logout', {
      method: 'POST',
      credentials: 'include',
    });
  } catch { /* игнорируем */ }
  window.location.replace('/index.html');
}
