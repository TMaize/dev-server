<script>
  import router from '../router.js'
  /**
   * @type {{name: string, type: string}}
   */
  export let file = undefined

  function clickRow(ev) {
    if (ev.type === 'keydown' && ev.key !== 'Enter') return
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

<li on:click={clickRow} on:keydown={clickRow}>
  {#if file.type == 'dir'}
    <div class="icon-base icon-dir" />
  {:else}
    <div class="icon-base icon-file" />
  {/if}
  <span>{file.name}</span>
</li>

<style lang="less">
  li {
    // display: flex;
    // flex-flow: row nowrap;
    // align-items: center;
    // padding: 8px 16px;
    // text-decoration: none;
    // border-top: 1px solid #d0d7de;
    // font-size: 16px;
  }
</style>
