document.addEventListener("DOMContentLoaded", function () {
  const app = document.querySelector(".app");
  const themeToggle = document.querySelector(".theme-toggle");

  // Defensive check: ensure the root element exists
  if (!app) {
    console.warn("App element not found");
    return;
  }

  // Safely read the saved theme (might throw in private mode)
  let preferred;
  try {
    preferred = localStorage.getItem("theme");
  } catch (e) {
    console.warn("localStorage not available:", e);
    preferred = null;
  }

  if (preferred === "light") {
    app.classList.remove("dark");
    app.classList.add("light");
  } else {
    app.classList.add("dark");
    app.classList.remove("light");
  }

  if (themeToggle) {
    themeToggle.addEventListener("click", () => {
      if (app.classList.contains("dark")) {
        app.classList.remove("dark");
        app.classList.add("light");
        try {
          localStorage.setItem("theme", "light");
        } catch (e) {
          console.warn("Failed to save theme preference:", e);
        }
      } else {
        app.classList.remove("light");
        app.classList.add("dark");
        try {
          localStorage.setItem("theme", "dark");
        } catch (e) {
          console.warn("Failed to save theme preference:", e);
        }
      }
    });
  }
});

function goToLogin() {
  window.location.href = "/login";
}