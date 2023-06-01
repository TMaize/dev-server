import App from './App.svelte'
import router from './router'
import './style.css'


router.init()

const app = new App({
  target: document.getElementById('app')
})

export default app
