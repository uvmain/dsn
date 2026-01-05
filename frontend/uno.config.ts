import { defineConfig, presetWebFonts, presetWind3, transformerDirectives, transformerVariantGroup } from 'unocss'

export default defineConfig({
  shortcuts: {
    'btn': 'px-4 py-2 rounded inline-block bg-primary-600 text-white cursor-pointer hover:bg-primary-700 disabled:cursor-default disabled:bg-gray-600 disabled:opacity-50',
    'btn-orange': 'btn bg-orange-500 hover:bg-orange-600',
    'icon-btn': 'text-[0.9em] inline-block cursor-pointer select-none opacity-75 transition duration-200 ease-in-out hover:opacity-100 hover:text-teal-600',
  },
  theme: {
    colors: {
      primary: {
        50: '#f0fdfa',
        100: '#ccfbf1',
        200: '#99f6e4',
        300: '#5eead4',
        400: '#2dd4bf',
        500: '#14b8a6',
        600: '#0d9488',
        700: '#0f766e',
        800: '#115e59',
        900: '#134e4a',
      },
    },
  },
  presets: [
    presetWind3(),
    presetWebFonts({
      provider: 'google',
      fonts: {
        sans: 'Inter:400,500,600,700',
        mono: 'Fira Code:400,500',
      },
    }),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
})
