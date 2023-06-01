<script>
  import { onMount } from 'svelte'
  import router from './router'
  import FileList from './component/FileList.svelte'

  let files = [{ name: 1 }, { name: 'aaa' }]

  function loadData(ev) {
    const path = ev && ev.path ? ev.path : router.current()
    console.log(ev, path)
  }

  router.on('change', loadData)

  onMount(loadData)
</script>

<div>
  <FileList {files} />
</div>
