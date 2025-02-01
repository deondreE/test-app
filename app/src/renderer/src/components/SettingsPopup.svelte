<script lang="ts">
  import { onMount } from "svelte";

  interface UserSettings {
    sessionId: string;
    fontSize: string;
    usernameColor: string;
    // TODO: user chat-log
  }

  let user: { sessionId: string } = $state();

  let settings: UserSettings | undefined = { sessionId: '', usernameColor: '', fontSize: ''};
  onMount(() => {
    let temp = localStorage.getItem("deondreE");
    if (temp) user = { sessionId: temp };
  });
</script>

<div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
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
      <button class="px-4 py-2 bg-gray-700 text-gray-300 rounded-md hover:bg-gray-600 transition">
        Cancel
      </button>
      <button class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-500 transition">
        Save
      </button>
    </div>
  </div>
</div>
