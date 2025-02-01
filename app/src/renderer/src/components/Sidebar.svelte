<script lang="ts">
  import { onMount } from "svelte";
  import Icon from '@iconify/svelte';
  import { goto } from "$app/navigation";
  import SettingsPopup from "./SettingsPopup.svelte";
  let showSettings = $state(false);

  let activeServer = {
    id: 0,
    name: 'Test Server',
    channels: [
      {
        id: '0',
        name: 'general'
      }
    ]
  };

  let user: any = $state({});

  onMount(() => {
    let temp = localStorage.getItem("deondreE");
    if (temp) user = { sessionId: temp };
  });

  const handleClick = async (channelId: string, serverId: number) => {
    goto(`/server/${serverId}/channel/${channelId}`);
  }
  
  $effect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key === ",") {
        event.preventDefault();
        showSettings = !showSettings;
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  })
</script>

<div class="fixed top-0 left-0 w-[20%] h-screen bg-gray-800 text-white p-4">
  <h2 class="text-2xl font-semibold">{activeServer.name}</h2>
  <hr class="border-gray-600 my-2" />
  
  {#each activeServer.channels as channel}
    <button class="text-left w-full py-2 px-4 rounded-lg hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500" onclick={() => handleClick(channel.id, activeServer.id)}>
      {channel.name}
    </button>
  {/each}
  {#if !showSettings}
  <div class="flex">
    <button class="absolute bottom-0 left-0 text-gray-400 border-none p-4 w-auto h-auto active:bg-green-100">
      deondreE
      <Icon icon="material-symbols:settings-account-box-sharp" width="24" height="24" />
    </button>
  </div>
  {:else}
  <SettingsPopup />
  {/if}
  </div>