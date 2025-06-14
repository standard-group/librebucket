<script>
  // No script logic needed for this static page layout.
  // This is a Svelte component file (.svelte).
</script>

<div class="TV">
  <header class="overlap-5">
    <div class="header-content">
      <div class="new-librebucket-logo">
        <img
          src="img/new-librebucket-logo-white.svg"
          alt="Librebucket logo part"
        />
      </div>
    </div>
  </header>

  <main class="div">
    <div class="main-content-wrapper">
      <section class="top-section">
        <div class="text-wrapper-5">Librebucket</div>
        <p class="text-wrapper-9-1">
          Create the future, <a href="#" class="text-wrapper-9-1-2">
            commit after commit</a
          >
        </p>
        <p class="text-wrapper-9">
          Librebucket is a self-hosted lightweight software forge.
        </p>
        <p class="text-wrapper-10">
          Easy to install and low maintenance, it just does the job.
        </p>
        <p class="you-are-currently">
          You are currently visiting an instance, means that we do not owned it
          and, we are not <br />responsible for any things which on this forge.
        </p>
        <div class="button-group">
          <div class="group-2">
            <div class="div-wrapper register-button">
              <div class="join-register">Join &amp; register</div>
            </div>
          </div>
          <div class="group-3">
            <div class="div-wrapper login-button">
              <div class="text-wrapper-11">Log In</div>
            </div>
          </div>
        </div>
        <div class="overlap-8">
          <img class="group-4" src="img/group-10-1.svg" alt="Group image" />
          <img class="group-5" src="img/group-17-1.svg" alt="Group image" />
        </div>
      </section>

      <section class="overlap">
        <div class="text-wrapper">Why Librebucket?</div>

        <div class="feature-grid">
          <div class="feature-item">
            <img class="feather" src="img/feather.svg" alt="Feather icon" />
            <div>
              <div class="lightweight-fast">Lightweight &amp; fast</div>
              <p class="p">
                Librebucket is less than 500MB and lightning fast, even
                Raspberry Pi can run it!
              </p>
            </div>
          </div>
          <div class="feature-item">
            <img
              class="coffee-to-go"
              src="img/coffee-to-go.svg"
              alt="Coffee To Go icon"
            />
            <div>
              <div class="text-wrapper-2">
                Coffee To Go, easy install &amp; maintain
              </div>
              <p class="text-wrapper-3">
                Our "bucket" is packed with tools and stuff already, for
                developers, by developers.
              </p>
            </div>
          </div>
          <div class="feature-item">
            <img class="git" src="img/git.svg" alt="Git icon" />
            <div>
              <div class="text-wrapper-2">
                Self-hosted alternative to GitHub and GitLab
              </div>
              <p class="text-wrapper-3">
                GitHub and GitLab are full of telemetry. Librebucket is an
                awesome alternative!
              </p>
            </div>
          </div>
          <div class="feature-item">
            <img
              class="heart-cog"
              src="img/heart-cog.svg"
              alt="Heart Cog icon"
            />
            <div>
              <div class="text-wrapper-2">100% Free and OSS</div>
              <p class="text-wrapper-3">
                Standard Group loves open-source. Librebucket will be forever
                free and open-source!
              </p>
            </div>
          </div>
        </div>
      </section>

      <section class="overlap-6">
        <div class="text-wrapper-6">Get Involved</div>
        <p class="librebucket-is-made">
          Librebucket is made by many motivated people, and we are looking
          forward to your contribution. Small or big contributions, you are
          still making big steps for Librebucket.
        </p>
        <div class="contribution-buttons">
          <div class="overlap-wrapper">
            <div class="overlap-7 contribute-button">
              <img class="github" src="img/github.svg" alt="GitHub icon" />
              <div class="text-wrapper-7">Contribute on Github</div>
            </div>
          </div>
          <div class="overlap-group-wrapper">
            <div class="overlap-7 donate-button">
              <img
                class="heart-cog-2"
                src="img/heart-cog.svg"
                alt="Heart Cog icon"
              />
              <div class="text-wrapper-8">Donate</div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </main>

  <footer class="overlap-9">
    <div class="footer-left">
      <p class="text-wrapper-14">Version: 0.0.0-debug</p>
      <p class="text-wrapper-16">This home page is running on Librebucket</p>
      <p class="text-wrapper-12">
        Copyright © 2025 Librebucket authors. Content available under CC BY-SA
        4.0.
      </p>
      <p class="text-wrapper-13">
        Credit to the Forgejo's contributors for design inspiration.
      </p>
    </div>
    <div class="footer-right">
      <a href="#" class="text-wrapper-15">Licenses</a>
      <a href="#"><img class="api" src="img/api.svg" alt="API icon" /></a>
      <a href="#"
        ><img class="github-2" src="img/github.svg" alt="GitHub icon" /></a
      >
    </div>
  </footer>
</div>

<style>
  :root {
    --primary-blue: #3487e7;
    --dark-bg: #171e26;
    --dark-secondary-bg: #10161d;
    --dark-card-bg: #242e38;
    --text-primary: #ffffff;
    --text-secondary: #ffffff66;
    --gradient-blue: linear-gradient(90deg, #3487e7 0%, #205fa8 100%);
    --gradient-orange: linear-gradient(90deg, #e77634 0%, #cdb623 100%);
    --gradient-red: linear-gradient(90deg, #e73434 0%, #811d1d 100%);
  }

  .TV {
    background-color: var(--dark-bg);
    width: 100%;
    height: 100vh;
    display: flex;
    justify-content: center;
    position: relative;
    font-family: "Noto Sans", Helvetica, sans-serif;
  }

  .overlap-5 {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 80px;
    background-color: var(--dark-secondary-bg);
    z-index: 100;
    display: flex;
    justify-content: center;
    border-bottom: 1px solid #ffffff1a;
  }

  .header-content {
    width: 100%;
    max-width: 1280px;
    display: flex;
    align-items: center;
    padding: 0 2rem;
  }

  .new-librebucket-logo {
    display: flex;
    align-items: center;
    height: 60px;
    margin-top: 10px;
  }

  .div {
    width: 100%;
    max-width: 1280px;
    height: 100vh;
    position: relative;
    box-sizing: border-box;
    padding-top: 80px;
    padding-bottom: 165px;
    overflow-y: auto;
  }

  .main-content-wrapper {
    padding: 2rem;
  }

  main::-webkit-scrollbar {
    display: none;
  }

  .overlap-9 {
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 165px;
    background-color: var(--dark-secondary-bg);
    border-top: 1px solid #ffffff33;
    z-index: 100;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 2rem;
    box-sizing: border-box;
  }

  .footer-left {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .footer-right {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }

  .footer-right a,
  .footer-right img {
    position: static;
    filter: invert(90%) sepia(0%) saturate(7457%) hue-rotate(109deg)
      brightness(122%) contrast(103%);
  }

  .footer-right .api,
  .footer-right .github-2 {
    width: 36px;
    height: 36px;
  }

  .text-wrapper-15 {
    color: var(--text-secondary);
    font-size: 16px;
    text-decoration: none;
    transition: color 0.2s;
  }
  .text-wrapper-15:hover {
    color: var(--text-primary);
  }

  .top-section {
    position: relative;
    padding-top: 4rem;
    padding-bottom: 5rem;
  }

  .overlap {
    background-color: var(--dark-card-bg);
    border: 1px solid #ffffff33;
    padding: 4rem 2rem;
    border-radius: 1rem;
    margin-top: 2rem;
    position: relative;
  }

  .text-wrapper {
    text-align: center;
    background: var(--gradient-blue);
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
    color: transparent;
    font-weight: 700;
    font-size: 64px;
    letter-spacing: -3.2px;
    margin-bottom: 3rem;
  }

  .feature-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
    gap: 2.5rem;
    align-items: start;
  }

  .feature-item {
    display: flex;
    align-items: flex-start;
    gap: 1.5rem;
  }

  .feature-item img {
    width: 63px;
    height: 63px;
    filter: invert(40%) sepia(92%) saturate(879%) hue-rotate(187deg)
      brightness(97%) contrast(87%);
    flex-shrink: 0;
  }

  .lightweight-fast,
  .text-wrapper-2 {
    font-weight: 700;
    color: var(--text-primary);
    font-size: 24px;
    letter-spacing: -1.2px;
    margin-bottom: 0.5rem;
  }

  .p,
  .text-wrapper-3 {
    font-weight: 400;
    color: var(--text-secondary);
    font-size: 16px;
    letter-spacing: -0.8px;
    line-height: 1.5;
  }

  .text-wrapper-5 {
    font-weight: 700;
    color: var(--text-primary);
    font-size: 64px;
    letter-spacing: -3.2px;
  }

  .overlap-6 {
    padding: 4rem 0;
  }

  .text-wrapper-6 {
    font-weight: 700;
    color: var(--text-primary);
    font-size: 64px;
    letter-spacing: -3.2px;
    margin-bottom: 1rem;
  }

  .librebucket-is-made {
    width: 100%;
    max-width: 984px;
    font-size: 24px;
    letter-spacing: -1.2px;
    font-weight: 400;
    color: var(--text-secondary);
    line-height: 1.6;
    margin-bottom: 2rem;
  }

  .contribution-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
  }

  .overlap-7 {
    width: 296px;
    height: 66px;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    cursor: pointer;
    transition: transform 0.2s;
  }

  .overlap-7:hover {
    transform: scale(1.03);
  }

  .overlap-7.contribute-button {
    background: var(--gradient-blue);
  }

  .overlap-7.donate-button {
    background: var(--gradient-red);
  }

  .github,
  .heart-cog-2 {
    filter: invert(90%) sepia(0%) saturate(7457%) hue-rotate(109deg)
      brightness(122%) contrast(103%);
    width: 40px;
    height: 40px;
  }

  .text-wrapper-7,
  .text-wrapper-8 {
    font-weight: 500;
    color: var(--text-primary);
    font-size: 24px;
    letter-spacing: -1.2px;
  }

  .text-wrapper-9,
  .text-wrapper-10,
  .you-are-currently {
    font-weight: 400;
    color: var(--text-secondary);
    font-size: 24px;
    letter-spacing: -1.2px;
    margin-top: 1rem;
  }

  .text-wrapper-9-1 {
    font-weight: 400;
    color: var(--text-primary);
    font-size: 32px;
    letter-spacing: -1.6px;
  }

  .text-wrapper-9-1-2 {
    background: var(--gradient-blue);
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
    color: transparent;
    text-decoration: none;
  }

  .you-are-currently {
    font-size: 16px;
    letter-spacing: -0.8px;
  }

  .button-group {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    margin-top: 2rem;
  }

  .div-wrapper {
    width: 211px;
    height: 66px;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: transform 0.2s;
  }

  .div-wrapper:hover {
    transform: scale(1.03);
  }

  .div-wrapper.register-button {
    background: var(--gradient-blue);
  }

  .join-register {
    font-weight: 500;
    color: var(--text-primary);
    font-size: 24px;
    letter-spacing: -1.2px;
  }

  .div-wrapper.login-button {
    background: var(--gradient-orange);
  }

  .text-wrapper-11 {
    font-weight: 500;
    color: var(--text-primary);
    font-size: 24px;
    letter-spacing: -1.2px;
  }

  .overlap-8 {
    position: absolute;
    width: 225px;
    height: 303px;
    top: 5rem;
    right: 5rem;
    display: none;
  }

  @media (min-width: 1024px) {
    .overlap-8 {
      display: block;
    }
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

  .text-wrapper-12,
  .text-wrapper-13,
  .text-wrapper-14,
  .text-wrapper-16 {
    font-weight: 400;
    color: var(--text-secondary);
    font-size: 15px;
    letter-spacing: -0.75px;
  }

  @media (max-width: 900px) {
    .main-content-wrapper {
      padding: 1rem !important;
    }
    .text-wrapper, .text-wrapper-5, .text-wrapper-6 {
      font-size: 2rem !important;
      letter-spacing: -1px !important;
    }
    .feature-grid {
      grid-template-columns: 1fr !important;
      gap: 1rem !important;
    }
    .footer-left, .footer-right {
      flex-direction: column !important;
      align-items: flex-start !important;
      gap: 0.5rem !important;
    }
    .overlap-5, .overlap-9 {
      height: auto !important;
      padding: 0.5rem 1rem !important;
    }
    .div-wrapper, .overlap-7 {
      width: 100% !important;
      min-width: 0 !important;
      max-width: 100% !important;
    }
    .top-section {
      padding-top: 2rem !important;
      padding-bottom: 2rem !important;
    }
    .librebucket-is-made {
      font-size: 1rem !important;
    }
    .text-wrapper-9, .text-wrapper-10, .you-are-currently {
      font-size: 1rem !important;
    }
    .text-wrapper-9-1 {
      font-size: 1.2rem !important;
    }
    .footer-right .api, .footer-right .github-2 {
      width: 28px !important;
      height: 28px !important;
    }
    .github, .heart-cog-2 {
      width: 28px !important;
      height: 28px !important;
    }
    .overlap-8 {
      display: none !important;
    }
  }

  @media (max-width: 600px) {
    .text-wrapper, .text-wrapper-5, .text-wrapper-6 {
      font-size: 1.2rem !important;
      letter-spacing: -0.5px !important;
    }
    .feature-item .lightweight-fast, .feature-item .text-wrapper-2 {
      font-size: 1rem !important;
    }
    .librebucket-is-made {
      font-size: 0.9rem !important;
    }
    .text-wrapper-9, .text-wrapper-10, .you-are-currently {
      font-size: 0.9rem !important;
    }
    .text-wrapper-9-1 {
      font-size: 1rem !important;
    }
    .text-wrapper-7, .text-wrapper-8, .join-register, .text-wrapper-11 {
      font-size: 1rem !important;
    }
    .footer-left .text-wrapper-12, .footer-left .text-wrapper-13, .footer-left .text-wrapper-14, .footer-left .text-wrapper-16 {
      font-size: 0.8rem !important;
    }
  }
</style>
