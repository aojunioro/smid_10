import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, type RenderResult } from 'vitest-browser-react'
import { type Locator, userEvent } from 'vitest/browser'
import { UserAuthForm } from './user-auth-form'

const FORM_MESSAGES = {
  loginEmpty: 'Informe o login.',
  passwordEmpty: 'Informe a senha.',
} as const

const {
  navigate,
  setUserMock,
  setAccessTokenMock,
  setAuthTokenMock,
  loginMock,
} = vi.hoisted(() => ({
  navigate: vi.fn(),
  setUserMock: vi.fn(),
  setAccessTokenMock: vi.fn(),
  setAuthTokenMock: vi.fn(),
  loginMock: vi.fn(),
}))

vi.mock('@/features/auth/api', () => ({
  login: loginMock,
}))

vi.mock('@/lib/api/client', () => ({
  apiClient: { setAuthToken: setAuthTokenMock },
  ApiError: class ApiError extends Error {},
}))

vi.mock('@/stores/auth-store', () => ({
  useAuthStore: () => ({
    auth: {
      setUser: setUserMock,
      setAccessToken: setAccessTokenMock,
    },
  }),
}))

vi.mock('@tanstack/react-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@tanstack/react-router')>()
  return {
    ...actual,
    useNavigate: () => navigate,
    Link: ({
      children,
      to,
      className,
      ...rest
    }: {
      children?: React.ReactNode
      to: string
      className?: string
    }) => (
      <a href={to} className={className} {...rest}>
        {children}
      </a>
    ),
  }
})

describe('UserAuthForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    loginMock.mockResolvedValue({
      token: 'jwt-token',
      expires_at: '2030-01-01T00:00:00Z',
      user: {
        id: 1,
        login: 'admin',
        name: 'Admin',
        email: 'admin@example.com',
        system_unit_id: 1,
        active: true,
      },
    })
  })

  describe('Rendering without redirectTo', () => {
    let screen: RenderResult
    let loginInput: Locator
    let passwordInput: Locator
    let signInButton: Locator

    beforeEach(async () => {
      screen = await render(<UserAuthForm />)
      loginInput = screen.getByRole('textbox', { name: /^Login$/i })
      passwordInput = screen.getByLabelText(/^Senha$/i)
      signInButton = screen.getByRole('button', { name: /^Entrar$/i })
    })

    it('renders fields and submit button', async () => {
      await expect.element(loginInput).toBeInTheDocument()
      await expect.element(passwordInput).toBeInTheDocument()
      await expect.element(signInButton).toBeInTheDocument()
    })

    it('shows validation messages when submitting empty form', async () => {
      await userEvent.click(signInButton)

      await expect
        .element(screen.getByText(FORM_MESSAGES.loginEmpty))
        .toBeInTheDocument()
      await expect
        .element(screen.getByText(FORM_MESSAGES.passwordEmpty))
        .toBeInTheDocument()
    })

    it('authenticates and navigates to default route on success', async () => {
      await userEvent.fill(loginInput, 'admin')
      await userEvent.fill(passwordInput, 'secret')

      await userEvent.click(signInButton)

      await vi.waitFor(() => expect(loginMock).toHaveBeenCalledOnce())
      expect(loginMock).toHaveBeenCalledWith({
        login: 'admin',
        password: 'secret',
      })

      await vi.waitFor(() => expect(setUserMock).toHaveBeenCalledOnce())
      expect(setUserMock).toHaveBeenCalledWith(
        expect.objectContaining({
          id: 1,
          login: 'admin',
          name: 'Admin',
          email: 'admin@example.com',
          systemUnitId: 1,
          exp: expect.any(Number),
        })
      )
      expect(setAccessTokenMock).toHaveBeenCalledWith('jwt-token')
      expect(setAuthTokenMock).toHaveBeenCalledWith('jwt-token')

      await vi.waitFor(() =>
        expect(navigate).toHaveBeenCalledWith({ to: '/', replace: true })
      )
    })
  })

  it('navigates to redirectTo when provided', async () => {
    const { getByRole, getByLabelText } = await render(
      <UserAuthForm redirectTo='/settings' />
    )

    await userEvent.fill(getByRole('textbox', { name: /Login/i }), 'admin')
    await userEvent.fill(getByLabelText('Senha'), 'secret')

    await userEvent.click(getByRole('button', { name: /Entrar/i }))

    await vi.waitFor(() => expect(setUserMock).toHaveBeenCalledOnce())
    expect(setAccessTokenMock).toHaveBeenCalledWith('jwt-token')

    await vi.waitFor(() =>
      expect(navigate).toHaveBeenCalledWith({
        to: '/settings',
        replace: true,
      })
    )
  })
})
