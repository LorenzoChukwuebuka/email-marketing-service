import type { Config } from 'tailwindcss'
import daisyui from 'daisyui'

const config: Config = {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {}
  },
  daisyui: {
    themes: []
  },
  plugins: [daisyui]
}

export default config