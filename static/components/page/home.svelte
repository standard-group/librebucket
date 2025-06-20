<script>
  let isDarkMode = $state(true);

  function toggleTheme() {
    isDarkMode = !isDarkMode;
  }

  let logoSrc = isDarkMode
    ? "img/new-librebucket-logo-white.svg"
    : "img/new-librebucket-logo.svg";

  function goToLogin() {
    window.location.href = "/login";
  }
</script>

<svelte:head>
  <title>Librebucket</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
</svelte:head>

<div class="app" class:dark={isDarkMode} class:light={!isDarkMode}>
  <header class="header">
    <div class="header-content">
      <div class="logo">
        <img src={logoSrc || "/placeholder.svg"} alt="Librebucket logo" />
      </div>
      <button class="theme-toggle" onclick={toggleTheme}>
        {isDarkMode ? "☀️" : "🌙"}
      </button>
    </div>
  </header>

  <main class="main">
    <div class="main-content">
      <section class="hero-section">
        <div class="hero-content">
          <h1 class="main-title">Librebucket</h1>
          <p class="hero-subtitle">
            Create the future, <a href="https://example.com" class="hero-link">commit after commit</a>
          </p>
          <p class="hero-description">
            Librebucket is a self-hosted lightweight software forge.
          </p>
          <p class="hero-description">
            Easy to install and low maintenance, it just does the job.
          </p>
          <p class="instance-notice">
            You are currently visiting an instance, means that we do not owned it and, we are not responsible for any things which on this forge.
          </p>
          <div class="button-group">
            <button class="btn btn-primary" onclick={() => window.location.href = '/register'}>
              Join &amp; register
            </button>
            <button class="btn btn-secondary" onclick={goToLogin}>
              Sign In
            </button>
          </div>
        </div>
        <div class="hero-image">
          <img class="group-4" src="img/group-10-1.svg" alt="" />
          <img class="group-5" src="img/group-17-1.svg" alt="" />
        </div>
      </section>

      <section class="features-section">
        <h2 class="section-title">Why Librebucket?</h2>
        <div class="feature-grid">
          <div class="feature-item">
            <img src="img/feather.svg" alt="" />
            <div>
              <h3>Lightweight &amp; fast</h3>
              <p>
                Librebucket is less than 500MB and lightning fast, even Raspberry Pi can run it!
              </p>
            </div>
          </div>
          <div class="feature-item">
            <img src="img/coffee-to-go.svg" alt="" />
            <div>
              <h3>Coffee To Go, easy install &amp; maintain</h3>
              <p>
                Our "bucket" is packed with tools and stuff already, for developers, by developers.
              </p>
            </div>
          </div>
          <div class="feature-item">
            <img src="img/git.svg" alt="" />
            <div>
              <h3>Self-hosted alternative to GitHub and GitLab</h3>
              <p>
                GitHub and GitLab are full of telemetry. Librebucket is an awesome alternative!
              </p>
            </div>
          </div>
          <div class="feature-item">
            <img src="img/heart-cog.svg" alt="" />
            <div>
              <h3>100% Free and OSS</h3>
              <p>
                Standard Group loves open-source. Librebucket will be forever free and open-source!
              </p>
            </div>
          </div>
        </div>
      </section>

      <section class="contribute-section">
        <h2 class="section-title">Get Involved</h2>
        <p class="contribute-description">
          Librebucket is made by many motivated people, and we are looking forward to your contribution. Small or big contributions, you are still making big steps for Librebucket.
        </p>
        <div class="contribute-buttons">
          <button class="btn btn-contribute">
            <img src="img/github.svg" alt="" />
            Contribute on Github
          </button>
          <button class="btn btn-donate">
            <img src="img/heart-cog.svg" alt="" />
            Donate
          </button>
        </div>
      </section>
    </div>
  </main>

  <footer class="footer">
    <div class="footer-content">
      <div class="footer-left">
        <p>Version: 0.0.0-debug</p>
        <p>This home page is running on Librebucket</p>
        <p>
          Copyright © 2025 Librebucket authors. Content available under CC BY-SA 4.0.
        </p>
        <p>Credit to the Forgejo's contributors for design inspiration.</p>
      </div>
      <div class="footer-right">
        <a href="https://example.com">Licenses</a>
        <a href="https://example.com"><img src="img/api.svg" alt="" /></a>
        <a href="https://example.com"><img src="img/github.svg" alt="" /></a>
      </div>
    </div>
  </footer>
</div>

<style>
  :root {
    --primary-blue: #3487e7;
    --primary-blue-hover: #2976d6;
    --secondary-orange: #e77634;
    --secondary-orange-hover: #d66623;
    --danger-red: #e73434;
    --danger-red-hover: #d62323;

    /* Dark theme */
    --dark-bg: #171e26;
    --dark-secondary-bg: #10161d;
    --dark-card-bg: #242e38;
    --dark-text-primary: #ffffff;
    --dark-text-secondary: #ffffff66;
    --dark-border: #ffffff33;

    /* Light theme */
    --light-bg: #ffffff;
    --light-secondary-bg: #f8f9fa;
    --light-card-bg: #ffffff;
    --light-text-primary: #1a1a1a;
    --light-text-secondary: #666666;
    --light-border: #e0e0e0;

    --gradient-blue: linear-gradient(90deg, #3487e7 0%, #205fa8 100%);
    --gradient-orange: linear-gradient(90deg, #e77634 0%, #cdb623 100%);
    --gradient-red: linear-gradient(90deg, #e73434 0%, #811d1d 100%);
  }

  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  .app {
    min-height: 100vh;
    font-family: "Noto Sans", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    display: flex;
    flex-direction: column;
    transition: all 0.3s ease;
  }

  /* Dark theme styles */
  .app.dark {
    background-color: var(--dark-bg);
    color: var(--dark-text-primary);
  }

  .app.dark .header {
    background-color: var(--dark-secondary-bg);
    border-bottom: 1px solid var(--dark-border);
  }

  .app.dark .features-section {
    background-color: var(--dark-card-bg);
    border: 1px solid var(--dark-border);
  }

  .app.dark .footer {
    background-color: var(--dark-secondary-bg);
    border-top: 1px solid var(--dark-border);
  }

  /* Light theme styles */
  .app.light {
    background-color: var(--light-bg);
    color: var(--light-text-primary);
  }

  .app.light .header {
    background-color: var(--light-secondary-bg);
    border-bottom: 1px solid var(--light-border);
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  }

  .app.light .features-section {
    background-color: var(--light-card-bg);
    border: 1px solid var(--light-border);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  }

  .app.light .footer {
    background-color: var(--light-secondary-bg);
    border-top: 1px solid var(--light-border);
  }

  .app.light .hero-description,
  .app.light .instance-notice,
  .app.light .contribute-description,
  .app.light .feature-item p,
  .app.light .footer p,
  .app.light .footer a {
    color: var(--light-text-secondary);
  }

  .app.light .feature-item img {
    filter: invert(20%) sepia(92%) saturate(879%) hue-rotate(187deg) brightness(97%) contrast(87%);
  }

  .app.light .footer img {
    filter: invert(40%) sepia(0%) saturate(0%) hue-rotate(0deg) brightness(60%) contrast(100%);
  }

  .app.light .btn img {
    filter: invert(100%) sepia(0%) saturate(0%) hue-rotate(0deg) brightness(100%) contrast(100%);
  }

  .header {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: 70px;
    z-index: 100;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .header-content {
    width: 100%;
    max-width: 1280px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 1rem;
  }

  .logo {
    display: flex;
    align-items: center;
    height: 50px;
  }

  .logo img {
    height: 35px;
    transition: opacity 0.3s ease;
  }

  .theme-toggle {
    background: none;
    border: 2px solid var(--primary-blue);
    border-radius: 8px;
    padding: 8px 12px;
    font-size: 16px;
    cursor: pointer;
    transition: all 0.3s ease;
    color: inherit;
  }

  .theme-toggle:hover {
    background-color: var(--primary-blue);
    transform: scale(1.05);
  }

  .main {
    flex: 1;
    padding-top: 70px;
    padding-bottom: 160px;
    display: flex;
    justify-content: center;
    overflow-y: auto;
  }

  .main-content {
    width: 100%;
    max-width: 1280px;
    padding: 1rem;
  }

  .hero-section {
    display: grid;
    grid-template-columns: 1fr auto;
    gap: 4rem;
    align-items: center;
    padding: 2rem 0;
    min-height: 500px;
  }

  .hero-content {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .main-title {
    font-weight: 700;
    font-size: clamp(2.5rem, 8vw, 4rem);
    letter-spacing: -0.05em;
    margin-bottom: 1rem;
    line-height: 1.1;
  }

  .hero-subtitle {
    font-weight: 400;
    font-size: clamp(1.25rem, 4vw, 2rem);
    letter-spacing: -0.02em;
    line-height: 1.3;
  }

  .hero-link {
    background: var(--gradient-blue);
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
    text-decoration: none;
    font-weight: 500;
  }

  .hero-description {
    font-weight: 400;
    font-size: clamp(1rem, 3vw, 1.5rem);
    letter-spacing: -0.02em;
    margin-top: 0.5rem;
    line-height: 1.4;
  }

  .instance-notice {
    font-size: clamp(0.875rem, 2.5vw, 1rem);
    letter-spacing: -0.01em;
    margin-top: 1rem;
    line-height: 1.5;
  }

  .button-group {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    margin-top: 2rem;
  }

  .btn {
    padding: 16px 24px;
    border: none;
    border-radius: 10px;
    font-weight: 500;
    font-size: clamp(1rem, 2.5vw, 1.25rem);
    letter-spacing: -0.01em;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    color: white;
    text-decoration: none;
    min-height: 50px;
  }

  .btn:hover {
    transform: scale(1.03);
  }

  .btn-primary {
    background: var(--gradient-blue);
  }

  .btn-secondary {
    background: var(--gradient-orange);
  }

  .btn-contribute {
    background: var(--gradient-blue);
  }

  .btn-donate {
    background: var(--gradient-red);
  }

  .btn img {
    width: 20px;
    height: 20px;
    filter: invert(100%) sepia(0%) saturate(0%) hue-rotate(0deg) brightness(100%) contrast(100%);
  }

  .hero-image {
    position: relative;
    width: 225px;
    height: 303px;
    flex-shrink: 0;
  }

  .group-4,
  .group-5 {
    position: absolute;
    object-fit: cover;
  }

  .group-4 {
    width: 225px;
    height: 288px;
    top: 0;
    left: 0;
  }

  .group-5 {
    width: 128px;
    height: 125px;
    top: 178px;
    left: 0;
  }

  .features-section {
    padding: 3rem 1.5rem;
    border-radius: 1rem;
    margin: 2rem 0;
  }

  .section-title {
    text-align: center;
    background: var(--gradient-blue);
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
    font-weight: 700;
    font-size: clamp(2rem, 6vw, 4rem);
    letter-spacing: -0.05em;
    margin-bottom: 2rem;
    line-height: 1.1;
  }

  .feature-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
    align-items: start;
  }

  .feature-item {
    display: flex;
    align-items: flex-start;
    gap: 1rem;
  }

  .feature-item img {
    width: 50px;
    height: 50px;
    filter: invert(40%) sepia(92%) saturate(879%) hue-rotate(187deg) brightness(97%) contrast(87%);
    flex-shrink: 0;
  }

  .feature-item h3 {
    font-weight: 700;
    font-size: clamp(1.125rem, 3vw, 1.5rem);
    letter-spacing: -0.02em;
    margin-bottom: 0.5rem;
    line-height: 1.2;
  }

  .feature-item p {
    font-weight: 400;
    font-size: clamp(0.875rem, 2.5vw, 1rem);
    letter-spacing: -0.01em;
    line-height: 1.5;
  }

  .contribute-section {
    padding: 3rem 0;
    text-align: center;
  }

  .contribute-description {
    max-width: 800px;
    margin: 0 auto 2rem;
    font-size: clamp(1rem, 3vw, 1.5rem);
    letter-spacing: -0.02em;
    font-weight: 400;
    line-height: 1.6;
  }

  .contribute-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    justify-content: center;
  }

  .footer {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    min-height: 140px;
    z-index: 100;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .footer-content {
    width: 100%;
    max-width: 1280px;
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    padding: 1.5rem 1rem;
  }

  .footer-left {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    flex: 1;
  }

  .footer-left p {
    font-weight: 400;
    font-size: clamp(0.75rem, 2vw, 0.875rem);
    letter-spacing: -0.01em;
    line-height: 1.4;
  }

  .footer-right {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-shrink: 0;
  }

  .footer-right a {
    color: inherit;
    text-decoration: none;
    font-size: clamp(0.875rem, 2.5vw, 1rem);
    transition: opacity 0.2s;
  }

  .footer-right a:hover {
    opacity: 0.8;
  }

  .footer-right img {
    width: 28px;
    height: 28px;
    filter: invert(60%) sepia(0%) saturate(0%) hue-rotate(0deg) brightness(80%) contrast(100%);
  }

  /* Mobile-first responsive design */
  @media (max-width: 768px) {
    .header {
      height: 60px;
    }

    .header-content {
      padding: 0 1rem;
    }

    .logo img {
      height: 30px;
    }

    .theme-toggle {
      padding: 6px 10px;
      font-size: 14px;
    }

    .main {
      padding-top: 60px;
      padding-bottom: 180px;
    }

    .main-content {
      padding: 0.75rem;
    }

    .hero-section {
      grid-template-columns: 1fr;
      text-align: center;
      gap: 2rem;
      padding: 1.5rem 0;
      min-height: auto;
    }

    .hero-image {
      display: none;
    }

    .hero-content {
      gap: 0.75rem;
    }

    .main-title {
      margin-bottom: 0.75rem;
    }

    .button-group {
      flex-direction: column;
      gap: 0.75rem;
      margin-top: 1.5rem;
    }

    .btn {
      width: 100%;
      padding: 14px 20px;
    }

    .features-section {
      padding: 2rem 1rem;
      margin: 1.5rem 0;
    }

    .section-title {
      margin-bottom: 1.5rem;
    }

    .feature-grid {
      grid-template-columns: 1fr;
      gap: 1.5rem;
    }

    .feature-item {
      gap: 0.75rem;
    }

    .feature-item img {
      width: 40px;
      height: 40px;
    }

    .contribute-section {
      padding: 2rem 0;
    }

    .contribute-description {
      margin-bottom: 1.5rem;
    }

    .contribute-buttons {
      flex-direction: column;
      gap: 0.75rem;
      max-width: 300px;
      margin: 0 auto;
    }

    .footer {
      min-height: 160px;
      position: fixed;
    }

    .footer-content {
      flex-direction: column;
      gap: 1rem;
      padding: 1.5rem 1rem;
      align-items: center;
      text-align: center;
    }

    .footer-left {
      order: 2;
      align-items: center;
      gap: 0.25rem;
    }

    .footer-right {
      order: 1;
      gap: 1.5rem;
    }

    .footer-right img {
      width: 24px;
      height: 24px;
    }
  }

  @media (max-width: 480px) {
    .main-content {
      padding: 0.5rem;
    }

    .hero-section {
      padding: 1rem 0;
    }

    .hero-content {
      gap: 0.5rem;
    }

    .button-group {
      margin-top: 1rem;
    }

    .btn {
      padding: 12px 16px;
      min-height: 44px;
    }

    .features-section {
      padding: 1.5rem 0.75rem;
    }

    .feature-grid {
      gap: 1rem;
    }

    .contribute-section {
      padding: 1.5rem 0;
    }

    .footer-content {
      padding: 1rem 0.75rem;
    }

    .footer {
      min-height: 140px;
    }

    .main {
      padding-bottom: 150px;
    }
  }

  @media (max-width: 360px) {
    .main-content {
      padding: 0.25rem;
    }

    .features-section {
      padding: 1rem 0.5rem;
    }

    .footer-content {
      padding: 0.75rem 0.5rem;
    }
  }
</style>
