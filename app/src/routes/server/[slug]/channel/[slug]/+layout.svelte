<script>
  import "../../../../../renderer/src/assets/main.css";
  import Sidebar from "../../../../../renderer/src/components/Sidebar.svelte";
</script>

<main class="">
  <aside class="">
    <Sidebar />
  </aside>

  <section class="">
    <slot></slot>
  </section>
</main>

<style>
</style>
