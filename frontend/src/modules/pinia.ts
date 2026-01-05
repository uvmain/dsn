import type { UserModule } from '~/types'

export const install: UserModule = ({ app, isClient, initialState }) => {
  const pinia = createPinia()
  app.use(pinia)

  if (isClient) {
    pinia.state.value = typeof initialState.pinia !== 'undefined' ? initialState.pinia : {}
  }
  else {
    initialState.pinia = pinia.state.value
  }
}
