<script lang="ts">
  import { onMount } from "svelte";

  interface UserSettings {
    sessionId: string;
    fontSize: string;
    usernameColor: string;
    // TODO: user chat-log
  }

  // Initial settings, with dummy values as fallback
  let settings: UserSettings = { sessionId: '', fontSize: '16', usernameColor: '#ffffff' };
  
  // Control the visibility of the settings modal
  let showSettings = true;

  onMount(() => {
    let temp = localStorage.getItem("deondreE");
    if (temp) settings.sessionId = temp;
  });

  const updateSettings = (e: Event) => {
    e.preventDefault();
    localStorage.setItem("userSettings", JSON.stringify(settings));
    // Optionally close the modal after saving
    showSettings = false;
  }
</script>

<!-- Only render the settings modal if showSettings is true -->
{#if showSettings}
  <div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
    <div class="bg-[#1e1e1e] text-gray-200 p-6 rounded-lg shadow-lg w-[500px] border border-gray-700">
      <h2 class="text-lg font-semibold text-gray-100 border-b border-gray-700 pb-2 mb-4">
        User Settings
      </h2>

      <!-- Font Size Setting -->
      <div class="mb-4">
        <!-- svelte-ignore a11y_label_has_associated_control -->
        <label class="block text-gray-300 text-sm font-medium mb-1">Font Size</label>
        <input
          type="number"
          bind:value={settings.fontSize}
          class="w-full bg-[#252526] border border-gray-600 rounded-md p-2 text-gray-100 focus:ring focus:ring-blue-500 focus:border-blue-400 outline-none"
          min="12"
          max="30"
        />
      </div>

      <!-- Username Color Setting -->
      <div class="mb-4">
        <!-- svelte-ignore a11y_label_has_associated_control -->
        <label class="block text-gray-300 text-sm font-medium mb-1">Username Color</label>
        <input
          type="color"
          bind:value={settings.usernameColor}
          class="w-full h-10 border border-gray-600 rounded-md bg-[#252526] cursor-pointer"
        />
      </div>

      <!-- Buttons -->
      <div class="flex justify-end space-x-2 border-t border-gray-700 pt-4 mt-4">
        <button
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-500 transition"
          on:click={updateSettings}
        >
          Save
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Main content: Only rendered when settings modal is closed -->
{#if !showSettings}
  <div class="p-4">
    <!-- Your main application content goes here -->
    <p>This is the main content of the app.</p>
  </div>
{/if}
