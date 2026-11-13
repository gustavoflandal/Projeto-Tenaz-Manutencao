// Menu do usuário na barra superior (seção 5.3 do documento 7): quem está
// logado e em qual empresa, minha conta, troca de empresa (quando há mais de
// uma) e sair.

import { LogOut, UserRound } from 'lucide-react';
import { DropdownMenu } from 'radix-ui';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router';

import { Botao } from '../componentes/Botao';
import type { Sessao } from '../lib/tipos';
import { useSair, useTrocarEmpresa } from '../modulos/adm/sessao';

const ITEM =
  'flex h-10 cursor-default items-center gap-2 rounded-xs px-3 texto-corpo-p data-highlighted:bg-superficie-hover';

function iniciais(nome: string): string {
  const partes = nome.trim().split(/\s+/).filter(Boolean);
  const primeira = partes[0]?.[0] ?? '';
  const ultima = partes.length > 1 ? (partes[partes.length - 1]?.[0] ?? '') : '';
  return `${primeira}${ultima}`.toUpperCase();
}

export function MenuUsuario({ sessao }: { readonly sessao: Sessao }) {
  const { t } = useTranslation();
  const navegar = useNavigate();
  const sair = useSair();
  const trocar = useTrocarEmpresa();
  const outras = (sessao.empresas ?? []).filter((e) => !e.atual);

  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger asChild>
        <Botao variante="icone" aria-label={t('usuarioMenu.botao', { nome: sessao.usuario.nome })}>
          <span aria-hidden="true" className="inline-flex size-8 items-center justify-center rounded-total bg-borda texto-rotulo text-texto">
            {iniciais(sessao.usuario.nome) || <UserRound size={16} />}
          </span>
        </Botao>
      </DropdownMenu.Trigger>
      <DropdownMenu.Content
        align="end"
        sideOffset={4}
        className="z-(--tz-z-suspenso) min-w-64 rounded-md border border-borda bg-superficie-elevada p-1 text-texto shadow-3"
      >
        <div className="px-3 py-2">
          <p className="texto-rotulo">{sessao.usuario.nome}</p>
          <p className="texto-legenda text-texto-secundario">{sessao.usuario.email}</p>
          <p className="texto-legenda text-texto-secundario">{sessao.empresa.nome}</p>
        </div>
        <DropdownMenu.Separator className="my-1 h-px bg-borda" />
        <DropdownMenu.Item className={ITEM} onSelect={() => void navegar('/conta')}>
          {t('usuarioMenu.minhaConta')}
        </DropdownMenu.Item>
        {outras.length > 0 && (
          <>
            <DropdownMenu.Separator className="my-1 h-px bg-borda" />
            <DropdownMenu.Label className="px-3 py-2 texto-sobrelinha text-texto-secundario">
              {t('usuarioMenu.trocarEmpresa')}
            </DropdownMenu.Label>
            {outras.map((empresa) => (
              <DropdownMenu.Item
                key={empresa.id}
                className={ITEM}
                disabled={trocar.isPending}
                onSelect={() => trocar.mutate(empresa.id, { onSuccess: () => void navegar('/') })}
              >
                {empresa.nome}
              </DropdownMenu.Item>
            ))}
          </>
        )}
        <DropdownMenu.Separator className="my-1 h-px bg-borda" />
        <DropdownMenu.Item
          className={ITEM}
          onSelect={() => {
            // Navega antes de esvaziar a sessão, para a guarda não guardar este endereço como retorno.
            void navegar('/entrar', { replace: true });
            sair.mutate();
          }}
        >
          <LogOut size={16} aria-hidden="true" />
          {t('usuarioMenu.sair')}
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  );
}
