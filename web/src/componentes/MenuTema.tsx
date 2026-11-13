// Troca de tema: igual ao sistema, claro, escuro ou sol (seção 3.3 do
// documento 7). Fica na barra superior até o menu do usuário existir.

import { Check, Contrast, Monitor, Moon, Sun, type LucideIcon } from 'lucide-react';
import { DropdownMenu } from 'radix-ui';
import { useTranslation } from 'react-i18next';

import { PREFERENCIAS_DE_TEMA, useTema, type PreferenciaTema } from '../lib/tema';

import { Botao } from './Botao';

const ICONES: Record<PreferenciaTema, LucideIcon> = {
  sistema: Monitor,
  claro: Sun,
  escuro: Moon,
  sol: Contrast,
};

function ehPreferencia(valor: string): valor is PreferenciaTema {
  return (PREFERENCIAS_DE_TEMA as readonly string[]).includes(valor);
}

export function MenuTema() {
  const { t } = useTranslation();
  const { preferencia, definir } = useTema();
  const Icone = ICONES[preferencia];

  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger asChild>
        <Botao variante="icone" aria-label={t('tema.botao', { tema: t(`tema.${preferencia}`) })}>
          <Icone size={20} aria-hidden="true" />
        </Botao>
      </DropdownMenu.Trigger>
      {/* Sem portal: o menu fica dentro do cabeçalho (marco banner), na ordem
          de leitura logo após o botão que o abre. */}
      <DropdownMenu.Content
        align="end"
        sideOffset={4}
        className="z-(--tz-z-suspenso) min-w-56 rounded-md border border-borda bg-superficie-elevada p-1 text-texto shadow-3"
      >
        <DropdownMenu.Label className="px-3 py-2 texto-sobrelinha text-texto-secundario">
          {t('tema.rotulo')}
        </DropdownMenu.Label>
        <DropdownMenu.RadioGroup
          value={preferencia}
          onValueChange={(valor) => {
            if (ehPreferencia(valor)) definir(valor);
          }}
        >
          {PREFERENCIAS_DE_TEMA.map((opcao) => (
            <DropdownMenu.RadioItem
              key={opcao}
              value={opcao}
              className="flex h-10 cursor-default items-center gap-2 rounded-xs px-3 texto-corpo-p data-highlighted:bg-superficie-hover"
            >
              <span className="inline-flex size-4 items-center justify-center">
                <DropdownMenu.ItemIndicator>
                  <Check size={16} aria-hidden="true" />
                </DropdownMenu.ItemIndicator>
              </span>
              {t(`tema.${opcao}`)}
            </DropdownMenu.RadioItem>
          ))}
        </DropdownMenu.RadioGroup>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  );
}
