<script lang="ts">
  import { onMount } from "svelte";

  let message = $state("");
  let server_message = $state("");
  let messages: any[] = $state([]);
  let messsage_output = $state("");
  let socket: WebSocket;

  onMount(() => {
    let temp = localStorage.getItem("messages");
    if (temp) messages = JSON.parse(temp);
    socket = new WebSocket("ws://localhost:8080/message");

    socket.addEventListener("open", () => {
      console.log("WebSocket connected!");
      server_message += "Status: Connected\n";
    });

    socket.addEventListener("message", (event) => {
      let d = JSON.parse(event.data);
      const dtRegx = /(\d{2}:\d{2})/;
      const match = d.sent_time.match(dtRegx);
      if (match) console.log(match[1]);
      let n = {
        userId: d.id,
        message: `${d.msg}`,
        sessionId: d.session,
        time: match[1],
        username: d.username,
      };
      messages = [...messages, n];
      localStorage.setItem("messages", JSON.stringify(messages));
      localStorage.setItem(`${n.username}`, n.sessionId);
      messsage_output += `UserId: ${d.id}:Message: ${d.msg}:Session Id: ${d.session}<br>`;
    });

    socket.addEventListener("error", (event) => {
      messsage_output += "Error: Connection error\n";
      console.error("WebSocket error:", event);
    });

    socket.addEventListener("close", () => {
      messsage_output += "Status: Disconnected\n";
    });
  });

  const send = () => {
    if (message.trim() !== "") {
      let m = {
        message: message,
        username: "deondreE",
      };
      socket.send(JSON.stringify(m));
      message = "";
    } else {
      server_message += "Please enter a message before sending.\n";
    }
  };

  const update = (e) => {
    if (e.key === "Enter") {
      send();
    }
  };
</script>

<div class="block">
  {#each messages as user}
    <div
      class="flex items-start space-x-4 p-4 rounded-lg bg-white mb-4 overflow-hidden max-w-3xl"
    >
      <div class="flex flex-col flex-grow">
        <div class="flex justify-between w-full">
          <span class="font-semibold text-purple-400 text-lg"
            >{user.username}</span
          >
          <span class="text-sm text-gray-500">{user.time}</span>
        </div>
        <span class="whitespace-pre-line text-md font-sans font-medium text-gray-800 mt-2"
          >{user.message}</span
        >
      </div>
    </div>
  {/each}
  <!-- New Messages -->
  <span>
    {@html messsage_output}
  </span>
</div>
<div class="absolute inset-x-0 bottom-0 p-4">
  <textarea
    bind:value={message}
    class="ml-[180px] w-[36rem] h-32 p-4 text-lg border-2 border-gray-300 rounded-md shadow-md focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
    placeholder="Type your message here..."
    onkeydown={update}
  ></textarea>
</div>
