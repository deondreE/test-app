<script lang="ts">
  import { onMount } from "svelte";

  let message = $state("");
  let settings = $state();
  let server_message = $state("");
  let messages: any[] = $state([]);
  let messsage_output = $state("");
  let socket: WebSocket; 

  onMount(() => {
    let temp = localStorage.getItem("messages");
    let temp1 = localStorage.getItem("settings");
    if (temp1) settings = JSON.parse(temp1);
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

  const adjustHeight = (e: any) => {
    // @ts-ignore
    event.target.style.height = "4rem";
    // @ts-ignore 
    event.target.style.height = event.target.scrollHeight + "px";
  }

  const deleteMessage = (id: string) => {
    messages = messages.filter((msg) => msg.id !== id);
  }

  const update = (e: KeyboardEvent) => {
    if (e.key === "Enter") {
      send();
    }
  };
</script>

<div class="block">
  {#each messages as user}
    <div
      class="flex items-end space-x-4 p-8 rounded-lg bg-white overflow-hidden max-w-2xl"
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
      <button
        class="absolute top-2 right-2 bg-red-500 text-white p-2 rounded-full opacity-0 group-hover:opacity-100 transition-opacity duration-200"
        onclick={() => deleteMessage(user.id)}
      >
      ✕
    </button>
    </div>
  {/each}
  <span>
    {@html messsage_output}
  </span>
</div>
<div class="relative inset-x-0 bottom-0 p-4">
  <textarea
    bind:value={message}
    class="ml-[2rem] w-[36rem] min-h-[4rem] p-4 text-lg border border-gray-300 rounded-xl shadow-md focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white transition-all duration-200 ease-in-out resize-none"
    placeholder="Type your message here..."
    oninput={(e) => adjustHeight(e)}
    onkeydown={update}
  ></textarea>
</div>
