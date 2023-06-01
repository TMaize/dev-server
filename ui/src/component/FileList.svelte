<script>
  import FileItem from './FileItem.svelte'
  /**
   * @type {object[]}
   */
  export let files = []
</script>

<div class="files-wrap">
  <h1>~</h1>
  <ul>
    {#each files as file}
      <FileItem {file} />
    {/each}
  </ul>
</div>

<style lang="less">
  .files-wrap {
    width: 96%;
    max-width: 880px;
    margin: 0 auto;
    border: 1px solid #d0d7de;
    border-radius: 6px;
    overflow: hidden;

    h1 {
      height: 55px;
      line-height: 55px;
      color: #24292f;
      font-size: 14px;
      padding: 0 16px;
      margin: 0;
      background-color: #f6f8fa;
      box-sizing: border-box;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    ul {
      list-style: none;
      padding: 0;
      margin: 0;
    }
  }
</style>
