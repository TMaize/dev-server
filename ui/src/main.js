import App from './App.svelte'
import router from './router'
import './style.css'
import './icons/index.css'


router.init()

const app = new App({
  target: document.getElementById('app')
})

export default app
