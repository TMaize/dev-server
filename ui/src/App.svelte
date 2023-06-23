<script>
  import { onMount } from 'svelte'
  import router from './router'
  import FileList from './component/FileList.svelte'
  import PageFooter from './component/PageFooter.svelte'
  let files = []

  async function loadData(ev) {
    const path = ev && ev.path ? ev.path : router.current()
    const data = await fetch('/:' + path).then(resp => resp.json())

    if (Array.isArray(data.data)) {
      const list = data.data
      list.sort(function (a, b) {
        return a.type.charCodeAt(0) - b.type.charCodeAt(0)
      })
      files = list
    }
  }

  router.on('change', loadData)

  onMount(loadData)
</script>

<FileList {files} />
<PageFooter />
