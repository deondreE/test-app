<script>
  import Sidebar from "../../../../../renderer/src/components/Sidebar.svelte";
  import "../../../../../renderer/src/assets/main.css";
</script>

<main class="flex">
  <aside class="">
    <Sidebar />
  </aside>

  <section class="content">
    <slot></slot>
  </section>
</main>

<style>
  .content {
    flex: 1;
    padding: 20px;
    margin-left: 150px; /* Push content away from sidebar */
    transition: margin-left 0.3s ease;
  }
</style>
