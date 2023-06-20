import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const dirname = path.resolve(fileURLToPath(import.meta.url), '..')

const files = fs.readdirSync(dirname)

let content = fs.readFileSync(path.resolve(dirname, 'base.css'), 'utf-8')

for (let i = 0; i < files.length; i++) {
  const file = files[i]
  if (!file.endsWith('.svg')) continue

  const buffer = fs.readFileSync(path.resolve(dirname, file))
  const dataUri = `data:image/svg+xml;base64,${buffer.toString('base64')}`

  const name = file.replace(/.svg$/, '')

  content += `\n.icon-${name} {\n  --data: url(${dataUri});\n  mask-image: var(--data);\n  -webkit-mask-image: var(--data);\n}\n`
}

fs.writeFileSync(path.resolve(dirname, 'index.css'), content)
