const fs = require('fs')
const axios = require('axios').default
const package = require('../package.json')

async function download(url) {
  console.log('[download]', url)
  const filename = url.split('/').pop()
  const writer = fs.createWriteStream('./bin/' + filename)
  return axios.request({ method: 'GET', url, responseType: 'stream' }).then(resp => resp.data.pipe(writer))
}

async function main() {
  fs.rmdirSync('bin', { recursive: true })
  fs.mkdirSync('bin')

  const version = package.version

  await download(`https://github.com/TMaize/dev-server/releases/download/v${version}/dev-server-linux-x64.tar.gz`)
  await download(`https://github.com/TMaize/dev-server/releases/download/v${version}/dev-server-mac-x64.tar.gz`)
  await download(`https://github.com/TMaize/dev-server/releases/download/v${version}/dev-server-windows-x64.zip`)
}

main()
