import type { UserModule } from '~/types'

export const install: UserModule = ({ app, router, routes, isClient, initialState }) => {
  const pinia = createPinia()
  app.use(pinia)

  // Sync Pinia state during SSG
  if (isClient) {
    pinia.state.value = typeof initialState.pinia !== 'undefined' ? initialState.pinia : {}
  }
  else {
    initialState.pinia = pinia.state.value
  }
}
