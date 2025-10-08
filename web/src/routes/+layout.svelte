<header>
  <nav>
    <div class="left">
      <a href="/" class="logo">localapps</a>
      <a href="https://github.com/CyberL1/localapps" target="_blank">
        Github repository
      </a>
    </div>
  </nav>
</header>

<main>
  <slot />
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
      "Helvetica Neue", Arial, sans-serif;
    background: #f0f0f0;
  }

  main {
    margin-top: 20px;
    margin-left: 20px;
  }

  header {
    background-color: #24292e;
    padding: 10px 20px;
    user-select: none;
  }

  nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  nav a {
    color: #ffffff;
    text-decoration: none;
    margin-right: 20px;
    font-size: 14px;
    line-height: 1.5;
  }

  nav > .left {
    display: flex;
    align-items: center;
  }

  nav > .left > a.logo {
    font-weight: bold;
    font-size: 16px;
  }

  nav a:hover {
    text-decoration: underline;
  }
</style>
