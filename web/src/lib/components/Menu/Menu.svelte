<script lang="ts">
  import type { Props } from "./types";
  import { onMount, tick } from "svelte";

  let { id = "menu", data, items }: Props = $props();

  let isVisible = $state(false);
  let posX = $state(0);
  let posY = $state(0);

  // svelte-ignore non_reactive_update
  let menu: HTMLElement;

  export async function open(event: MouseEvent, menuData: object) {
    event.stopPropagation();

    const rect = (event.target as HTMLElement).getBoundingClientRect();

    data = menuData;
    isVisible = true;

    await tick();

    const menuRect = menu.getBoundingClientRect();

    posX = event.clientX;
    posY = rect.bottom + window.scrollY;

    if (posX + menuRect.width > window.innerWidth) {
      posX = window.innerWidth - menuRect.width - 10;
    }

    if (posY + menuRect.height > window.innerHeight + window.scrollY) {
      posY = rect.top + window.scrollY - menuRect.height;
    }
  }

  export function close() {
    isVisible = false;
  }

  onMount(() => {
    const handler = (event: MouseEvent) => {
      if (isVisible && !menu.contains(event.target as Node)) {
        close();
      }
    };

    window.addEventListener("click", handler);
    return () => window.removeEventListener("click", handler);
  });
</script>

{#if isVisible}
  <div bind:this={menu} {id} style="left: {posX}px; top: {posY}px;">
    {#each items as item}
      <button onclick={item.onclick}>{item.title}</button>
    {/each}
  </div>
{/if}

<style>
  div {
    position: absolute;
    background: white;
    border: 1px solid #ccc;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    min-width: 150px;
    border-radius: 8px;
    user-select: none;
    z-index: 1000;
  }

  div button {
    cursor: pointer;
    display: flex;
    background: none;
    border: none;
    width: 100%;
    padding: 10px 16px;
    font-size: 1em;
  }

  div button:hover {
    background-color: lightgrey;
  }
</style>
