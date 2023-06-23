<script>
  import router from '../router.js'
  /**
   * @type {object[]}
   */
  export let files = []

  function fmtSize(size) {
    if (size > 1048576) {
      return Number(size / 1048576).toFixed(1) + 'M'
    }
    return Number(size / 1024).toFixed(1) + 'K'
  }

  function jump(file) {
    if (file.type == 'dir') {
      const nextPath = `${router.current()}/${file.name}/`.replace(/\/{2,}/, '/')
      router.push(nextPath)
    }
    if (file.type == 'file') {
      const nextPath = `${router.current()}/${file.name}`.replace(/\/{2,}/, '/')
      window.open(nextPath, '_blank')
    }
  }

  function handleEvent(ev, file) {
    switch (ev.type) {
      case 'keyup':
        if (ev.key == 'Enter') {
          jump(file)
        }
        break
      case 'keydown':
        if (ev.key === 'Enter') {
          ev.preventDefault()
        }
        break
      default:
        ev.preventDefault()
        jump(file)
    }
  }
</script>

<div class="files-wrap">
  <h1>~</h1>
  <div class="list">
    {#each files as file}
      <a href={file.name} on:click={e => handleEvent(e, file)} on:keydown={e => handleEvent(e, file)} on:keyup={e => handleEvent(e, file)}>
        {#if file.type == 'file'}
          <i class="ic ic-file" />
        {/if}

        {#if file.type == 'dir'}
          <i class="ic ic-dir" />
        {/if}

        <span class="name">{file.name}</span>

        {#if file.type == 'file'}
          <span class="size">{fmtSize(file.size)}</span>
        {/if}
      </a>
    {/each}
  </div>
</div>

<style lang="less">
  .files-wrap {
    width: 96%;
    max-width: 880px;
    margin: 0 auto;
    border: 1px solid #d0d7de;
    border-radius: 6px;

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

    .list {
      a {
        display: flex;
        flex-flow: row nowrap;
        align-items: center;
        padding: 8px 16px;
        text-decoration: none;
        border-top: 1px solid #d0d7de;

        .ic {
          margin-right: 14px;
          flex-shrink: 0;
          font-size: 15px;
        }

        .ic-file {
          color: #57606a;
        }

        .ic-dir {
          color: #54aeff;
        }

        .name {
          font-size: 14px;
          color: #333333;
          flex-grow: 1;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .size {
          flex-shrink: 0;
          white-space: nowrap;
          margin-left: 16px;
          color: #666666;
          font-size: 12px;
        }
      }
    }
  }
</style>
