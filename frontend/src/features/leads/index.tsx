import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Plus, RefreshCw } from 'lucide-react'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { ApiError } from '@/lib/api/client'
import { fetchLeads } from './api'
import { LeadSheet } from './components/lead-sheet'
import { LeadsTable } from './components/leads-table'
import type { Lead } from './types'

export function Leads() {
  const [sheetOpen, setSheetOpen] = useState(false)
  const [selectedLead, setSelectedLead] = useState<Lead | null>(null)

  const { data, isLoading, isError, error, refetch, isFetching } = useQuery({
    queryKey: ['leads'],
    queryFn: () => fetchLeads(),
  })

  const openCreate = () => {
    setSelectedLead(null)
    setSheetOpen(true)
  }

  const openEdit = (lead: Lead) => {
    setSelectedLead(lead)
    setSheetOpen(true)
  }

  return (
    <>
      <Header fixed>
        <Search className='me-auto' />
        <ThemeSwitch />
        <ConfigDrawer />
        <ProfileDropdown />
      </Header>

      <Main className='flex flex-1 flex-col gap-4 sm:gap-6'>
        <div className='flex flex-wrap items-end justify-between gap-2'>
          <div>
            <h2 className='text-2xl font-bold tracking-tight'>Leads</h2>
            <p className='text-muted-foreground'>
              Potenciais clientes do funil comercial.
            </p>
          </div>
          <div className='flex gap-2'>
            <Button
              variant='outline'
              size='sm'
              onClick={() => refetch()}
              disabled={isFetching}
            >
              <RefreshCw className={isFetching ? 'animate-spin' : ''} />
              Atualizar
            </Button>
            <Button size='sm' onClick={openCreate}>
              <Plus />
              Novo lead
            </Button>
          </div>
        </div>

        {isLoading ? (
          <p className='text-sm text-muted-foreground'>Carregando leads...</p>
        ) : null}

        {isError ? (
          <p className='text-sm text-destructive'>
            {error instanceof ApiError
              ? error.message
              : 'Erro ao carregar leads.'}
          </p>
        ) : null}

        {data ? <LeadsTable leads={data.leads} onEdit={openEdit} /> : null}
      </Main>

      <LeadSheet
        open={sheetOpen}
        onOpenChange={setSheetOpen}
        lead={selectedLead}
      />
    </>
  )
}
