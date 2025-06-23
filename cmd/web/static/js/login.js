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

  // Password toggle
  const togglePasswordBtn = document.getElementById("toggle-password");
  const passwordInput = document.getElementById("password");
  const eyeIcon = document.getElementById("eye-icon");

  if (togglePasswordBtn && passwordInput && eyeIcon) {
    togglePasswordBtn.addEventListener("click", () => {
      const isPassword = passwordInput.type === "password";
      passwordInput.type = isPassword ? "text" : "password";
      eyeIcon.src = isPassword
        ? "/static/img/eye-off.svg"
        : "/static/img/eye.svg";
      eyeIcon.alt = isPassword ? "Hide password" : "Show password";
    });
  }
});