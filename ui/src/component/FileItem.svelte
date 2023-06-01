<script>
  export let file = {}
  let o = {
    count: 1
  }

  function inc(){
    o.count++
  }
</script>

<li on:click={inc}>
  {file.name} / {o.count}
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
