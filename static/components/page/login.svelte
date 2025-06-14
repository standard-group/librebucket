<script>
    let isDarkMode = true;

    let passwordFieldType = "password";

    function toggleTheme() {
        isDarkMode = !isDarkMode;
    }

    function togglePassword() {
        passwordFieldType =
            passwordFieldType === "password" ? "text" : "password";
    }
</script>

<div class="app" class:dark={isDarkMode} class:light={!isDarkMode}>
    <header class="header">
        <div class="header-content">
            <div class="logo">
                <img
                    src={isDarkMode
                        ? "img/new-librebucket-logo-white.svg"
                        : "img/new-librebucket-logo.svg"}
                    alt="Librebucket logo"
                />
            </div>
            <button class="theme-toggle" on:click={toggleTheme}>
                <span>
                    {#if isDarkMode}
                        ☀️
                    {:else}
                        🌙
                    {/if}
                </span>
            </button>
        </div>
    </header>

    <main class="main">
        <div class="login-container">
            <div class="login-card">
                <div class="login-header">
                    <img
                        src={isDarkMode
                            ? "img/new-librebucket-logo-white.svg"
                            : "img/new-librebucket-logo.svg"}
                        alt="Librebucket logo"
                        class="login-logo"
                    />
                    <h1 class="login-title">Librebucket</h1>
                </div>

                <form class="login-form" on:submit|preventDefault>
                    <div class="form-group">
                        <label for="username" class="form-label"
                            >Username or email</label
                        >
                        <input
                            type="text"
                            id="username"
                            name="username"
                            class="form-input"
                            required
                        />
                    </div>

                    <div class="form-group">
                        <label for="password" class="form-label">Password</label
                        >
                        <div class="password-container">
                            <input
                                type={passwordFieldType}
                                id="password"
                                name="password"
                                class="form-input"
                                required
                            />
                            <button
                                type="button"
                                class="password-toggle"
                                on:click={togglePassword}
                            >
                                <span>
                                    {#if passwordFieldType === "password"}
                                        👁️
                                    {:else}
                                        🙈
                                    {/if}
                                </span>
                            </button>
                        </div>
                        <a href="#" class="forgot-password"
                            >Forgot your password?</a
                        >
                    </div>

                    <button type="submit" class="btn btn-primary btn-signin">
                        Sign in
                    </button>

                    <div class="signin-footer">
                        <p class="register-text">
                            Don't have an account?
                            <a href="#" class="register-link">Register now</a>
                        </p>
                    </div>
                </form>
            </div>
        </div>
    </main>
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
        --dark-input-bg: #1a222a;
        --dark-input-border: #ffffff22;
        --dark-input-focus: #3487e7;

        /* Light theme */
        --light-bg: #ffffff;
        --light-secondary-bg: #f8f9fa;
        --light-card-bg: #ffffff;
        --light-text-primary: #1a1a1a;
        --light-text-secondary: #666666;
        --light-border: #e0e0e0;
        --light-input-bg: #ffffff;
        --light-input-border: #d0d7de;
        --light-input-focus: #3487e7;

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
        font-family:
            "Noto Sans",
            -apple-system,
            BlinkMacSystemFont,
            "Segoe UI",
            Roboto,
            sans-serif;
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

    .app.dark .login-card {
        background-color: var(--dark-card-bg);
        border: 1px solid var(--dark-border);
    }

    .app.dark .form-input {
        background-color: var(--dark-input-bg);
        border-color: var(--dark-input-border);
        color: var(--dark-text-primary);
    }

    .app.dark .form-input:focus {
        border-color: var(--dark-input-focus);
        box-shadow: 0 0 0 3px rgba(52, 135, 231, 0.1);
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

    .app.light .login-card {
        background-color: var(--light-card-bg);
        border: 1px solid var(--light-border);
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
    }

    .app.light .form-input {
        background-color: var(--light-input-bg);
        border-color: var(--light-input-border);
        color: var(--light-text-primary);
    }

    .app.light .form-input:focus {
        border-color: var(--light-input-focus);
        box-shadow: 0 0 0 3px rgba(52, 135, 231, 0.1);
    }

    .app.light .form-label,
    .app.light .register-text {
        color: var(--light-text-secondary);
    }

    .header {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        height: 80px;
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
        padding: 0 2rem;
    }

    .logo {
        display: flex;
        align-items: center;
        height: 60px;
    }

    .logo img {
        height: 40px;
        transition: opacity 0.3s ease;
    }

    .theme-toggle {
        background: none;
        border: 2px solid var(--primary-blue);
        border-radius: 8px;
        padding: 8px 12px;
        font-size: 18px;
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
        padding-top: 80px;
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: calc(100vh - 80px);
    }

    .login-container {
        width: 100%;
        max-width: 400px;
        padding: 2rem;
    }

    .login-card {
        border-radius: 12px;
        padding: 3rem 2.5rem;
        transition: all 0.3s ease;
    }

    .login-header {
        text-align: center;
        margin-bottom: 2rem;
    }

    .login-logo {
        height: 50px;
        margin-bottom: 1rem;
        transition: opacity 0.3s ease;
    }

    .login-title {
        font-size: 2rem;
        font-weight: 700;
        letter-spacing: -1px;
        background: var(--gradient-blue);
        -webkit-background-clip: text;
        background-clip: text;
        -webkit-text-fill-color: transparent;
    }

    .login-form {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .form-label {
        font-size: 14px;
        font-weight: 500;
        letter-spacing: -0.2px;
    }

    .form-input {
        width: 100%;
        padding: 12px 16px;
        border: 2px solid;
        border-radius: 8px;
        font-size: 16px;
        transition: all 0.2s ease;
        font-family: inherit;
    }

    .form-input:focus {
        outline: none;
    }

    .password-container {
        position: relative;
        display: flex;
        align-items: center;
    }

    /* Making room for the icon inside the input */
    .password-container .form-input {
        padding-right: 40px;
    }

    .password-toggle {
        position: absolute;
        right: 1px;
        top: 1px;
        bottom: 1px;
        width: 38px;
        background: none;
        border: none;
        cursor: pointer;
        font-size: 16px;
        color: inherit;
        opacity: 0.7;
        transition: opacity 0.2s ease;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .password-toggle:hover {
        opacity: 1;
    }

    .forgot-password {
        color: var(--primary-blue);
        text-decoration: none;
        font-size: 14px;
        align-self: flex-end;
        margin-top: 0.25rem;
        transition: opacity 0.2s ease;
    }

    .forgot-password:hover {
        opacity: 0.8;
    }

    .btn {
        padding: 14px 24px;
        border: none;
        border-radius: 8px;
        font-weight: 600;
        font-size: 16px;
        cursor: pointer;
        transition: all 0.2s ease;
        color: white;
        font-family: inherit;
    }

    .btn:hover {
        transform: translateY(-1px);
    }

    .btn-primary {
        background: var(--gradient-blue);
    }

    .btn-signin {
        margin-top: 0.5rem;
    }

    .register-link {
        color: var(--primary-blue);
        text-decoration: none;
        transition: opacity 0.2s ease;
    }

    .signin-footer {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        margin-top: 1.5rem;
        text-align: center;
    }

    .register-text {
        font-size: 14px;
        line-height: 1.5;
    }

    .register-link {
        color: var(--primary-blue);
        text-decoration: none;
        transition: opacity 0.2s ease;
    }

    .register-link:hover {
        opacity: 0.8;
    }

    @media (max-width: 480px) {
        .login-container {
            padding: 1rem;
        }

        .login-card {
            padding: 2rem 1.5rem;
        }

        .header-content {
            padding: 0 1rem;
        }

        .theme-toggle {
            padding: 6px 10px;
            font-size: 16px;
        }
    }
</style>
