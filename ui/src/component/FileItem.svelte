<script>
  import router from '../router.js'
  /**
   * @type {{name: string, type: string}}
   */
  export let file = undefined

  function clickRow() {
    if (file.type == 'dir') {
      const nextPath = `${router.current()}/${file.name}/`.replace(/\/{2,}/, '/')
      router.push(nextPath)
    }
    if (file.type == 'file') {
      const nextPath = `${router.current()}/${file.name}`.replace(/\/{2,}/, '/')
      window.open(nextPath, '_blank')
    }
  }
</script>

<li on:click={clickRow}>
  {#if file.type == 'dir'}
    <img src="../assets/dir.svg" alt="" />
  {:else}
    <img src="../assets/file.svg" alt="" />
  {/if}
  <span>{file.name}</span>
</li>

<style lang="less">
  li {
    display: flex;
    flex-flow: row nowrap;
    align-items: center;
    padding: 8px 16px;
    text-decoration: none;
    border-top: 1px solid #d0d7de;
  }
</style>
