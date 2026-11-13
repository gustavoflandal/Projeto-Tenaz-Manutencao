// Campos de formulário (componente 6.2 do documento 7): rótulo sempre visível
// acima do campo, texto de apoio ou erro abaixo, "(opcional)" nos que podem
// ficar vazios. Erro com borda de 2 px, ícone e mensagem em `texto`, ligado
// ao campo por `aria-invalid` e `aria-describedby`.

import { ChevronDown, CircleX, Eye, EyeOff } from 'lucide-react';
import { Switch } from 'radix-ui';
import { createContext, useContext, useId, useState, type ComponentProps, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

import { TIPO_ERRO_DO_SERVIDOR } from '../lib/erros';
import { traduzirChave } from '../lib/i18n';

import { Botao } from './Botao';

/** Erro de um campo do formulário: do esquema (chave de texto) ou do servidor (texto pronto). */
export interface ErroDeCampo {
  readonly type?: string;
  readonly message?: string;
}

export function textoDoErro(erro: ErroDeCampo | undefined): string | undefined {
  if (!erro?.message) return undefined;
  return erro.type === TIPO_ERRO_DO_SERVIDOR ? erro.message : traduzirChave(erro.message);
}

// ---------------------------------------------------------------- identidade

const ContextoDoFormulario = createContext<string | null>(null);

export const ProvedorDoFormulario = ContextoDoFormulario.Provider;

/** Identificador do campo: estável por formulário, para o resumo de erros apontar para ele. */
export function idDoCampo(idDoFormulario: string, nome: string): string {
  return `${idDoFormulario}-${nome.replaceAll('.', '-')}`;
}

function useIdDoCampo(nome: string | undefined): string {
  const formulario = useContext(ContextoDoFormulario);
  const aleatorio = useId();
  return formulario && nome ? idDoCampo(formulario, nome) : aleatorio;
}

// --------------------------------------------------------------------- base

export const CLASSE_CAMPO =
  'h-10 w-full rounded-md border border-borda-controle bg-superficie px-3 texto-corpo text-texto hover:border-borda-controle-hover focus-visible:border-foco disabled:cursor-not-allowed disabled:border-borda disabled:bg-fundo disabled:text-texto-desabilitado aria-invalid:border-2 aria-invalid:border-status-parado [&[readonly]]:border-transparent [&[readonly]]:bg-transparent [&[readonly]]:hover:border-transparent [&[readonly]]:focus-visible:border-foco';

interface RotuloEMensagensProps {
  readonly id: string;
  readonly rotulo: string;
  readonly opcional?: boolean;
  readonly apoio?: string;
  readonly erro?: ErroDeCampo;
  readonly children: (ligacao: {
    readonly id: string;
    readonly 'aria-invalid': true | undefined;
    readonly 'aria-describedby': string | undefined;
  }) => ReactNode;
}

function RotuloEMensagens({ id, rotulo, opcional, apoio, erro, children }: RotuloEMensagensProps) {
  const { t } = useTranslation();
  const mensagem = textoDoErro(erro);
  const idMensagem = `${id}-mensagem`;
  const idApoio = `${id}-apoio`;
  const descricao = [apoio ? idApoio : null, mensagem ? idMensagem : null].filter(Boolean).join(' ');

  return (
    <div className="flex flex-col gap-2">
      <label htmlFor={id} className="texto-rotulo">
        {rotulo}
        {opcional && ' '}
        {opcional && <span className="ml-1 texto-legenda text-texto-secundario">{t('comum.opcional')}</span>}
      </label>
      {children({
        id,
        'aria-invalid': mensagem ? true : undefined,
        'aria-describedby': descricao || undefined,
      })}
      {apoio && (
        <p id={idApoio} className="texto-legenda text-texto-secundario">
          {apoio}
        </p>
      )}
      {mensagem && (
        <p id={idMensagem} className="flex items-start gap-2 texto-legenda text-texto">
          <CircleX size={16} aria-hidden="true" className="mt-0.5 shrink-0 text-status-parado" />
          <span>{mensagem}</span>
        </p>
      )}
    </div>
  );
}

// ------------------------------------------------------------------- campos

export interface CampoTextoProps extends Omit<ComponentProps<'input'>, 'id'> {
  readonly rotulo: string;
  readonly erro?: ErroDeCampo | undefined;
  readonly apoio?: string;
  readonly opcional?: boolean;
  /** Valor de dado (TAG, CPF, CNPJ, hora): JetBrains Mono. */
  readonly dado?: boolean;
}

export function CampoTexto({ rotulo, erro, apoio, opcional, dado, className, name, ...resto }: CampoTextoProps) {
  const id = useIdDoCampo(name);
  return (
    <RotuloEMensagens id={id} rotulo={rotulo} apoio={apoio} opcional={opcional} erro={erro}>
      {(ligacao) => (
        <input {...ligacao} name={name} className={[CLASSE_CAMPO, dado && 'texto-dado', className].filter(Boolean).join(' ')} {...resto} />
      )}
    </RotuloEMensagens>
  );
}

export function CampoSenha({ rotulo, erro, apoio, opcional, className, name, ...resto }: Omit<CampoTextoProps, 'type' | 'dado'>) {
  const { t } = useTranslation();
  const id = useIdDoCampo(name);
  const [visivel, setVisivel] = useState(false);
  return (
    <RotuloEMensagens id={id} rotulo={rotulo} apoio={apoio} opcional={opcional} erro={erro}>
      {(ligacao) => (
        <div className="relative">
          <input
            {...ligacao}
            name={name}
            type={visivel ? 'text' : 'password'}
            className={[CLASSE_CAMPO, 'pr-12', className].filter(Boolean).join(' ')}
            {...resto}
          />
          <Botao
            variante="icone"
            aria-label={visivel ? t('comum.ocultarSenha') : t('comum.mostrarSenha')}
            aria-pressed={visivel}
            aria-controls={id}
            onClick={() => setVisivel((atual) => !atual)}
            className="absolute top-0 right-0"
          >
            {visivel ? <EyeOff size={20} aria-hidden="true" /> : <Eye size={20} aria-hidden="true" />}
          </Botao>
        </div>
      )}
    </RotuloEMensagens>
  );
}

export interface CampoSelecaoProps extends Omit<ComponentProps<'select'>, 'id'> {
  readonly rotulo: string;
  readonly erro?: ErroDeCampo | undefined;
  readonly apoio?: string;
  readonly opcional?: boolean;
}

export function CampoSelecao({ rotulo, erro, apoio, opcional, className, name, children, ...resto }: CampoSelecaoProps) {
  const id = useIdDoCampo(name);
  return (
    <RotuloEMensagens id={id} rotulo={rotulo} apoio={apoio} opcional={opcional} erro={erro}>
      {(ligacao) => (
        <div className="relative">
          <select {...ligacao} name={name} className={[CLASSE_CAMPO, 'appearance-none pr-10', className].filter(Boolean).join(' ')} {...resto}>
            {children}
          </select>
          <ChevronDown
            size={20}
            aria-hidden="true"
            className="pointer-events-none absolute top-1/2 right-3 -translate-y-1/2 text-texto-secundario"
          />
        </div>
      )}
    </RotuloEMensagens>
  );
}

export interface CampoCaixaProps extends Omit<ComponentProps<'input'>, 'id' | 'type'> {
  readonly rotulo: string;
  readonly apoio?: string;
}

/** Caixa de marcação de 20 px com área clicável de 32 px (6.18). */
export function CampoCaixa({ rotulo, apoio, name, className, ...resto }: CampoCaixaProps) {
  const id = useIdDoCampo(name);
  const idApoio = `${id}-apoio`;
  return (
    <div className="flex items-start gap-1">
      <span className="flex size-8 shrink-0 items-center justify-center">
        <input
          id={id}
          name={name}
          type="checkbox"
          aria-describedby={apoio ? idApoio : undefined}
          className={['size-5 accent-texto', className].filter(Boolean).join(' ')}
          {...resto}
        />
      </span>
      <div className="flex flex-col pt-1">
        <label htmlFor={id} className="texto-corpo">
          {rotulo}
        </label>
        {apoio && (
          <p id={idApoio} className="texto-legenda text-texto-secundario">
            {apoio}
          </p>
        )}
      </div>
    </div>
  );
}

export interface InterruptorProps {
  readonly rotulo: string;
  readonly apoio?: string;
  readonly marcado: boolean;
  readonly aoMudar: (marcado: boolean) => void;
  readonly desabilitado?: boolean;
}

/** Liga/desliga de efeito imediato ou de formulário (6.3): o rótulo descreve o estado ligado. */
export function Interruptor({ rotulo, apoio, marcado, aoMudar, desabilitado }: InterruptorProps) {
  const id = useId();
  const idApoio = `${id}-apoio`;
  return (
    <div className="flex items-start gap-3">
      <Switch.Root
        id={id}
        checked={marcado}
        onCheckedChange={aoMudar}
        disabled={desabilitado}
        aria-describedby={apoio ? idApoio : undefined}
        className="group mt-0.5 inline-flex h-6 w-11 shrink-0 items-center rounded-total border border-borda-controle bg-superficie p-0.5 hover:enabled:border-texto-secundario aria-checked:border-texto aria-checked:bg-texto aria-checked:hover:enabled:bg-texto-secundario disabled:cursor-not-allowed disabled:opacity-60"
      >
        <Switch.Thumb className="block size-4 rounded-total bg-texto-secundario transition-transform duration-(--tz-duracao-rapida) group-aria-checked:translate-x-5 group-aria-checked:bg-superficie" />
      </Switch.Root>
      <div className="flex flex-col">
        <label htmlFor={id} className="texto-corpo">
          {rotulo}
        </label>
        {apoio && (
          <p id={idApoio} className="texto-legenda text-texto-secundario">
            {apoio}
          </p>
        )}
      </div>
    </div>
  );
}

export interface InterruptorDeTabelaProps {
  /** Nome acessível completo, com a tela e a captura: o interruptor não tem rótulo visível. */
  readonly rotulo: string;
  readonly marcado: boolean;
  readonly aoMudar: (marcado: boolean) => void;
  readonly desabilitado?: boolean;
}

/** Interruptor sem rótulo visível, para célula de tabela (o cabeçalho da coluna diz o que é). */
export function InterruptorDeTabela({ rotulo, marcado, aoMudar, desabilitado }: InterruptorDeTabelaProps) {
  return (
    <Switch.Root
      checked={marcado}
      onCheckedChange={aoMudar}
      disabled={desabilitado}
      aria-label={rotulo}
      className="group inline-flex h-6 w-11 shrink-0 items-center rounded-total border border-borda-controle bg-superficie p-0.5 hover:enabled:border-texto-secundario aria-checked:border-texto aria-checked:bg-texto aria-checked:hover:enabled:bg-texto-secundario disabled:cursor-not-allowed disabled:opacity-60"
    >
      <Switch.Thumb className="block size-4 rounded-total bg-texto-secundario transition-transform duration-(--tz-duracao-rapida) group-aria-checked:translate-x-5 group-aria-checked:bg-superficie" />
    </Switch.Root>
  );
}
