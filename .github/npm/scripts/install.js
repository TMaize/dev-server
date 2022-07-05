const fs = require('fs')
const os = require('os')
const path = require('path')
const cp = require('child_process')

const platform = os.platform() //win32 linux darwin ...
const arch = os.arch() // arm64 x64 ...

function getFileName() {
  const prefix = `dev-server-${platform}-${arch}`
  if (platform === 'win32') {
    return `${prefix}.zip`
  }
  return `${prefix}.tar.gz`
}

const filename = getFileName()
const file = path.resolve(__dirname, '..', 'bin', filename)

if (!fs.existsSync(file)) {
  console.log(`${filename} not found`)
  process.exit(1)
}

if (platform == 'win32') {
  cp.execSync(`Expand-Archive -Force ${filename} .`, { stdio: 'inherit', shell: 'powershell.exe', cwd: path.resolve(__dirname, '..', 'bin') })
} else {
  cp.execSync(`tar xf ${filename}`, { stdio: 'inherit', cwd: path.resolve(__dirname, '..', 'bin') })
}
