import { Pencil } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import type { Lead } from '../types'

type LeadsTableProps = {
  leads: Lead[]
  onEdit: (lead: Lead) => void
}

function formatDate(value?: string | null) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('pt-BR')
}

export function LeadsTable({ leads, onEdit }: LeadsTableProps) {
  if (leads.length === 0) {
    return (
      <p className='py-8 text-center text-sm text-muted-foreground'>
        Nenhum lead encontrado.
      </p>
    )
  }

  return (
    <div className='overflow-x-auto rounded-md border'>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className='min-w-16'>ID</TableHead>
            <TableHead className='min-w-40'>Nome</TableHead>
            <TableHead className='min-w-32'>Telefone</TableHead>
            <TableHead className='min-w-24'>Status</TableHead>
            <TableHead className='min-w-36'>Criado em</TableHead>
            <TableHead className='w-12' />
          </TableRow>
        </TableHeader>
        <TableBody>
          {leads.map((lead) => (
            <TableRow key={lead.id}>
              <TableCell>{lead.id}</TableCell>
              <TableCell className='font-medium'>{lead.nome}</TableCell>
              <TableCell>{lead.fone1}</TableCell>
              <TableCell>{lead.status_id}</TableCell>
              <TableCell>{formatDate(lead.criado_em)}</TableCell>
              <TableCell>
                <Button
                  variant='ghost'
                  size='icon'
                  aria-label='Editar lead'
                  onClick={() => onEdit(lead)}
                >
                  <Pencil className='h-4 w-4' />
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
