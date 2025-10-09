<script lang="ts">
  import { page } from "$app/state";
  import Menu from "$lib/components/Menu";
  import type { App } from "$lib/types/App";

  let menu: Menu;
  let menuData: App;

  function openMenu(event: MouseEvent, data: App) {
    menuData = data;
    menu.open(event, data);
  }
</script>

<div class="apps">
  {#each page.data.apps as app}
    <div class="app" id={app.id.toString()}>
      <div class="content">
        <div class="top">
          <button class="menu-opener" onclick={(e) => openMenu(e, app)}>
            +
          </button>
        </div>
        <a href={`//${app.name}.${page.data.domain}`} target="_blank">
          <img
            src={app.icon
              ? `/api/icons/apps/${app.icon}`
              : "https://placehold.co/60"}
            alt="An icon"
          />
          <span>{app.name}</span>
        </a>
      </div>
    </div>
  {/each}
</div>

<Menu
  bind:this={menu}
  data={menuData}
  items={[
    {
      title: "Open",
      onclick: () => {
        window.open(`//${menuData.name}.${page.data.domain}`);
        menu.close();
      },
    },
    {
      title: "Uninstall",
      onclick: async () => {
        if (confirm("Are you sure?")) {
          const req = await fetch(`/api/apps/${menuData.id.toString()}`, {
            method: "DELETE",
            headers: { Authorization: page.data.apiKey },
          });
          if (req.ok) {
            const appElement = document.getElementById(menuData.id.toString());
            if (appElement) {
              appElement.remove();
            }
          } else {
            alert("Failed to uninstall the app. Please try again.");
          }
        }
        menu.close();
      },
    },
  ]}
/>

<style>
  .apps {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(131px, 1fr));
    gap: 20px;
  }

  .app {
    background: white;
    border-radius: 10px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    padding: 20px;
    user-select: none;
    height: 80px;
    width: 92px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .app > .content {
    position: relative;
    bottom: 5px;
  }

  .app > .content > .top > .menu-opener {
    margin-left: 90px;
    border-radius: 50px;
    border-color: #419fff;
    font-size: 14px;
    width: 25px;
    height: 25px;
  }

  .app > .content > a {
    text-decoration: none;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .apps > .app > .content > a > img {
    width: 60px;
    height: 60px;
  }

  .apps > .app > .content > a > span {
    font-size: 1.2em;
    display: block;
    color: #333;
    overflow: hidden;
    text-overflow: ellipsis;
    text-wrap: nowrap;
    width: 95px;
  }

  .apps > .app > .content > a:hover > span {
    color: green;
  }
</style>
