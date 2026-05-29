import { useEffect } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { SidePanel } from '@/components/smid/side-panel'
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
import { ApiError } from '@/lib/api/client'
import { createLead, updateLead } from '../api'
import type { Lead } from '../types'

const formSchema = z.object({
  nome: z.string().min(1, 'Informe o nome.'),
  fone1: z.string().min(1, 'Informe o telefone principal.'),
  status_id: z.string().min(1, 'Informe o status.'),
  email: z.string().optional(),
})

type LeadFormValues = z.infer<typeof formSchema>

type LeadSheetProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  lead?: Lead | null
}

export function LeadSheet({ open, onOpenChange, lead }: LeadSheetProps) {
  const queryClient = useQueryClient()
  const isUpdate = Boolean(lead)

  const form = useForm<LeadFormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      nome: lead?.nome ?? '',
      fone1: lead?.fone1 ?? '',
      status_id: lead ? String(lead.status_id) : '1',
      email: lead?.email ?? '',
    },
  })

  useEffect(() => {
    if (!open) return
    form.reset({
      nome: lead?.nome ?? '',
      fone1: lead?.fone1 ?? '',
      status_id: lead ? String(lead.status_id) : '1',
      email: lead?.email ?? '',
    })
  }, [open, lead, form])

  const mutation = useMutation({
    mutationFn: async (values: LeadFormValues) => {
      const email = values.email?.trim() || undefined
      const statusId = Number(values.status_id)
      if (isUpdate && lead) {
        return updateLead(lead.id, {
          nome: values.nome,
          fone1: values.fone1,
          status_id: statusId,
          email,
        })
      }
      return createLead({
        nome: values.nome,
        fone1: values.fone1,
        status_id: statusId,
        email,
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['leads'] })
      toast.success(isUpdate ? 'Lead atualizado.' : 'Lead criado.')
      onOpenChange(false)
      form.reset()
    },
    onError: (error) => {
      const message =
        error instanceof ApiError ? error.message : 'Nao foi possivel salvar.'
      toast.error(message)
    },
  })

  return (
    <SidePanel
      open={open}
      onOpenChange={(next) => {
        onOpenChange(next)
        if (!next) form.reset()
      }}
      title={isUpdate ? 'Editar lead' : 'Novo lead'}
      description='Dados basicos do lead conforme SPEC_LEADS.'
    >
      <Form {...form}>
        <form
          className='flex flex-col gap-4'
          onSubmit={form.handleSubmit((values) => mutation.mutate(values))}
        >
          <FormField
            control={form.control}
            name='nome'
            render={({ field }) => (
              <FormItem>
                <FormLabel>Nome</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name='fone1'
            render={({ field }) => (
              <FormItem>
                <FormLabel>Telefone</FormLabel>
                <FormControl>
                  <Input inputMode='tel' {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name='status_id'
            render={({ field }) => (
              <FormItem>
                <FormLabel>Status (ID)</FormLabel>
                <FormControl>
                  <Input type='number' min={1} inputMode='numeric' {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name='email'
            render={({ field }) => (
              <FormItem>
                <FormLabel>E-mail</FormLabel>
                <FormControl>
                  <Input type='email' {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type='submit' disabled={mutation.isPending} className='mt-2'>
            {mutation.isPending ? 'Salvando...' : 'Salvar'}
          </Button>
        </form>
      </Form>
    </SidePanel>
  )
}
