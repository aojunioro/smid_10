import { useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Link, useNavigate } from '@tanstack/react-router'
import { Loader2, LogIn } from 'lucide-react'
import { toast } from 'sonner'
import { IconFacebook, IconGithub } from '@/assets/brand-icons'
import { login } from '@/features/auth/api'
import { ApiError, apiClient } from '@/lib/api/client'
import { useAuthStore } from '@/stores/auth-store'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { PasswordInput } from '@/components/password-input'

const formSchema = z.object({
  login: z.string().min(1, 'Informe o login.'),
  password: z.string().min(1, 'Informe a senha.'),
})

interface UserAuthFormProps extends React.HTMLAttributes<HTMLFormElement> {
  redirectTo?: string
}

export function UserAuthForm({
  className,
  redirectTo,
  ...props
}: UserAuthFormProps) {
  const [isLoading, setIsLoading] = useState(false)
  const navigate = useNavigate()
  const { auth } = useAuthStore()

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      login: '',
      password: '',
    },
  })

  async function onSubmit(data: z.infer<typeof formSchema>) {
    setIsLoading(true)
    try {
      const response = await login({
        login: data.login,
        password: data.password,
      })

      const expiresAt = new Date(response.expires_at).getTime()
      auth.setUser({
        id: response.user.id,
        login: response.user.login,
        name: response.user.name,
        email: response.user.email,
        systemUnitId: response.user.system_unit_id,
        exp: expiresAt,
      })
      auth.setAccessToken(response.token)
      apiClient.setAuthToken(response.token)

      const targetPath = redirectTo || '/'
      navigate({ to: targetPath, replace: true })
      toast.success(`Bem-vindo, ${response.user.name}!`)
    } catch (error) {
      const message =
        error instanceof ApiError
          ? error.message
          : 'Nao foi possivel entrar. Tente novamente.'
      toast.error(message)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className={cn('grid gap-3', className)}
        {...props}
      >
        <FormField
          control={form.control}
          name='login'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Login</FormLabel>
              <FormControl>
                <Input placeholder='seu.login' autoComplete='username' {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name='password'
          render={({ field }) => (
            <FormItem className='relative'>
              <FormLabel>Senha</FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder='********'
                  autoComplete='current-password'
                  {...field}
                />
              </FormControl>
              <FormMessage />
              <Link
                to='/forgot-password'
                className='absolute inset-e-0 -top-0.5 text-sm font-medium text-muted-foreground hover:opacity-75'
              >
                Esqueceu a senha?
              </Link>
            </FormItem>
          )}
        />
        <Button className='mt-2' disabled={isLoading} type='submit'>
          {isLoading ? <Loader2 className='animate-spin' /> : <LogIn />}
          Entrar
        </Button>

        <motionless-divider />

        <div className='grid grid-cols-2 gap-2'>
          <Button variant='outline' type='button' disabled={isLoading}>
            <IconGithub className='h-4 w-4' /> GitHub
          </Button>
          <Button variant='outline' type='button' disabled={isLoading}>
            <IconFacebook className='h-4 w-4' /> Facebook
          </Button>
        </div>
      </form>
    </Form>
  )
}

function motionless-divider() {
  return (
    <div className='relative my-2'>
      <div className='absolute inset-0 flex items-center'>
        <span className='w-full border-t' />
      </div>
      <motionless-divider-label />
    </motionless-divider>
  )
}

function motionless-divider-label() {
  return (
    <div className='relative flex justify-center text-xs uppercase'>
      <span className='bg-background px-2 text-muted-foreground'>
        Ou continue com
      </span>
    </motionless-divider-label>
  )
}
