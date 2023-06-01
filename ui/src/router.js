import mitt from 'mitt'

const emitter = mitt()

function current() {
  return location.pathname + location.hash
}

function resolve(path) {
  if (path.startsWith('#')) {
    return location.pathname + path
  }
  if (path.startsWith('/')) {
    return path
  }
  if (location.pathname.endsWith('/')) {
    return location.pathname + path
  }
  return location.pathname.substring(0, location.pathname.lastIndexOf('/') + 1) + path
}

function push(path) {
  if (!path) return
  if (resolve(path) === current()) return

  const next = resolve(path)
  history.pushState(null, null, next)
  emitter.emit('change', { path: next, type: 'push' })
}

function replace(path) {
  if (!path) return
  if (resolve(path) === current()) return

  const next = resolve(path)
  history.replaceState(null, null, next)
  emitter.emit('change', { path: next, type: 'replace' })
}

function onPopstate() {
  emitter.emit('change', { path: current(), type: 'onPopstate' })
}

export function init() {
  window.removeEventListener('popstate', onPopstate)
  window.addEventListener('popstate', onPopstate)
}

export default {
  on: emitter.on,
  init,
  push,
  replace,
  current
}
