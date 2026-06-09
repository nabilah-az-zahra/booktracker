import type { Config } from 'prettier'

const config: Config = {
    semi: false,
    singleQuote: true,
    tabWidth: 4,
    printWidth: 100,
    plugins: ['prettier-plugin-tailwindcss'],
}

export default config