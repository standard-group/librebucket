document.addEventListener("DOMContentLoaded", function () {
  const app = document.querySelector(".app");
  const themeToggle = document.querySelector(".theme-toggle");

  const preferred = localStorage.getItem("theme");
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
        localStorage.setItem("theme", "light");
      } else {
        app.classList.remove("light");
        app.classList.add("dark");
        localStorage.setItem("theme", "dark");
      }
    });
  }
});

function goToLogin() {
  window.location.href = "/login";
}
