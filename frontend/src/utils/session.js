const SESSION_TIMEOUT_MS = 30 * 60 * 1000;
const LAST_ACTIVITY_KEY = "lastActivityAt";

function now() {
  return Date.now();
}

function hasToken() {
  return Boolean(localStorage.getItem("token"));
}

export function touchSession() {
  if (!hasToken()) return;
  localStorage.setItem(LAST_ACTIVITY_KEY, String(now()));
}

export function clearAuthSession() {
  localStorage.removeItem("token");
  localStorage.removeItem("userRole");
  localStorage.removeItem(LAST_ACTIVITY_KEY);
}

export function isSessionExpired() {
  if (!hasToken()) return false;

  const raw = localStorage.getItem(LAST_ACTIVITY_KEY);
  if (!raw) {
    // Backward compatibility for existing sessions created before timeout feature.
    touchSession();
    return false;
  }

  const lastActivityAt = Number(raw);
  if (!Number.isFinite(lastActivityAt)) {
    touchSession();
    return false;
  }

  return now() - lastActivityAt > SESSION_TIMEOUT_MS;
}

export function startSessionTimeout(onTimeout) {
  if (typeof window === "undefined") {
    return () => {};
  }

  const globalKey = "__omahtaniSessionCleanup";
  if (typeof window[globalKey] === "function") {
    window[globalKey]();
  }

  let lastTouch = 0;
  const throttledTouch = () => {
    const current = now();
    if (current - lastTouch < 5000) return;
    lastTouch = current;
    touchSession();
  };

  const checkTimeout = () => {
    if (!hasToken()) return;
    if (!isSessionExpired()) return;
    clearAuthSession();
    onTimeout();
  };

  const events = [
    "click",
    "mousemove",
    "mousedown",
    "keydown",
    "scroll",
    "touchstart",
  ];

  events.forEach((eventName) => {
    window.addEventListener(eventName, throttledTouch, { passive: true });
  });

  const intervalId = window.setInterval(checkTimeout, 15000);
  checkTimeout();

  const cleanup = () => {
    events.forEach((eventName) => {
      window.removeEventListener(eventName, throttledTouch);
    });
    window.clearInterval(intervalId);
  };

  window[globalKey] = cleanup;
  return cleanup;
}

export { SESSION_TIMEOUT_MS, LAST_ACTIVITY_KEY };
