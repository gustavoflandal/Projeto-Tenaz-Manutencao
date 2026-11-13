// Convite de usuário (RF-ADM-004): nome, e-mail, perfis e, se quiser, a pessoa.
// A pessoa pode não existir (o usuário nasce sem pessoa), ser criada junto ou
// ser uma pessoa que já está no cadastro, sem usuário.

import { useMemo, useState } from 'react';

import type { Moeda } from '@tenaz/dominio';
import { Controller } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { BuscaPorTexto } from '../../../componentes/BuscaPorTexto';
import { CampoSelecao, CampoTexto } from '../../../componentes/Campo';
import { Formulario } from '../../../componentes/Formulario';
import { tratarErroDeEnvio } from '../../../lib/erros';
import {
  custoHoraOpcional,
  emailObrigatorio,
  precoUnitarioParaOContrato,
  textoObrigatorio,
  useFormulario,
  vazioParaIndefinido,
  z,
} from '../../../lib/formularios';
import type { EntradaUsuario } from '../../../lib/tipos';
import { usePessoas } from '../pessoas';
import { useUnidades } from '../unidades';
import { useMoedaDeCusto, useRegiao } from '../regiao';
import { useCriarUsuario, usePerfis } from '../usuarios';

import { SeletorDePerfis } from './SeletorDePerfis';

const esquemaDoFormulario = (moeda: Moeda | null) =>
  z
    .object({
      nome: textoObrigatorio('validacao.nomeObrigatorio').max(200, 'validacao.textoLongo'),
      email: emailObrigatorio,
      telefone: z.string().trim().max(30, 'validacao.textoLongo'),
      perfis: z.array(z.string()).min(1, 'validacao.perfilObrigatorio'),
      custo_hora: custoHoraOpcional(moeda),
      unidade_padrao_id: z.string(),
      pessoa: z.enum(['nenhuma', 'nova', 'existente']),
      pessoa_id: z.string(),
    })
    .refine((v) => v.pessoa !== 'existente' || v.pessoa_id !== '', { path: ['pessoa_id'], message: 'validacao.pessoaObrigatoria' });

export function FormularioDeUsuario() {
  const { t } = useTranslation();
  const custo = useMoedaDeCusto();
  const { moeda } = useRegiao();
  const esquema = useMemo(() => esquemaDoFormulario(moeda), [moeda]);
  const navegar = useNavigate();
  const criar = useCriarUsuario();
  const perfis = usePerfis();
  const unidades = useUnidades();
  const [mensagem, setMensagem] = useState<string | null>(null);
  const [buscaPessoa, setBuscaPessoa] = useState('');
  const formulario = useFormulario(esquema, {
    defaultValues: {
      nome: '',
      email: '',
      telefone: '',
      perfis: [],
      custo_hora: '',
      unidade_padrao_id: '',
      pessoa: 'nenhuma',
      pessoa_id: '',
    },
  });
  const { errors } = formulario.formState;
  const modoPessoa = formulario.watch('pessoa');
  const pessoas = usePessoas({ com_usuario: 'false', situacao: 'ativa', ...(buscaPessoa ? { busca: buscaPessoa } : {}) });

  const enviar = formulario.handleSubmit((v) => {
    setMensagem(null);
    const custoParaEnvio = vazioParaIndefinido(v.custo_hora) && precoUnitarioParaOContrato(v.custo_hora, moeda);
    const corpo: EntradaUsuario = {
      nome: v.nome,
      email: v.email,
      perfis: v.perfis,
      ...(vazioParaIndefinido(v.telefone) ? { telefone: v.telefone } : {}),
      ...(custoParaEnvio ? { custo_hora: custoParaEnvio } : {}),
      ...(v.unidade_padrao_id ? { unidade_padrao_id: v.unidade_padrao_id } : {}),
      ...(v.pessoa === 'nova' ? { nova_pessoa: { nome: v.nome, email: v.email } } : {}),
      ...(v.pessoa === 'existente' ? { pessoa_id: v.pessoa_id } : {}),
    };
    criar.mutate(corpo, {
      onSuccess: (usuario) => void navegar(`/administracao/usuarios/${usuario.id}`, { state: { conviteEnviado: usuario.convite_enviado } }),
      onError: (erro) =>
        tratarErroDeEnvio(erro, {
          setError: formulario.setError,
          definirMensagem: setMensagem,
          opcoes: { transformar: (campo) => campo.replace(/^nova_pessoa\./, '') },
        }),
    });
  });


  // O resumo de erros diz de qual campo é cada mensagem.
  const rotulosDoResumo: Record<string, string> = {
    nome: t('usuario.nome'),
    email: t('usuario.email'),
    telefone: t('usuario.telefone'),
    perfis: t('usuario.perfis'),
    custo_hora: t('usuario.custoHora', { moeda: custo.simbolo }),
    unidade_padrao_id: t('usuario.unidadePadrao'),
    pessoa_id: t('usuario.pessoaExistente'),
  };

  return (
    <Formulario erros={errors} rotulos={rotulosDoResumo} onSubmit={(e) => void enviar(e)} className="flex max-w-180 flex-col gap-6">
      {mensagem && <Alerta tipo="erro">{mensagem}</Alerta>}
      <p className="texto-corpo text-texto-secundario">{t('usuario.conviteTexto')}</p>
      <div className="grid gap-4 md:grid-cols-2">
        <CampoTexto rotulo={t('usuario.nome')} autoComplete="off" erro={errors.nome} {...formulario.register('nome')} />
        <CampoTexto rotulo={t('usuario.email')} type="email" autoComplete="off" erro={errors.email} {...formulario.register('email')} />
        <CampoTexto rotulo={t('usuario.telefone')} opcional type="tel" autoComplete="off" erro={errors.telefone} {...formulario.register('telefone')} />
        <CampoTexto rotulo={t('usuario.custoHora', { moeda: custo.simbolo })} opcional dado inputMode="decimal" apoio={custo.semCentavos ? t('usuario.custoHoraApoioSemCentavos') : t('usuario.custoHoraApoio')} erro={errors.custo_hora} {...formulario.register('custo_hora')} />
        <CampoSelecao rotulo={t('usuario.unidadePadrao')} opcional erro={errors.unidade_padrao_id} {...formulario.register('unidade_padrao_id')}>
          <option value="">{t('pessoa.semUnidade')}</option>
          {(unidades.data?.dados ?? []).map((u) => (
            <option key={u.id} value={u.id}>
              {u.nome}
            </option>
          ))}
        </CampoSelecao>
      </div>

      <Controller
        control={formulario.control}
        name="perfis"
        render={({ field, fieldState }) => (
          <SeletorDePerfis perfis={perfis.data?.dados ?? []} valor={field.value} aoMudar={field.onChange} erro={fieldState.error} />
        )}
      />
      {perfis.isError && <Alerta tipo="erro">{t('usuario.perfisFalhou')}</Alerta>}

      <fieldset className="flex flex-col gap-3">
        <legend className="mb-1 texto-rotulo">{t('usuario.pessoa')}</legend>
        <p className="texto-legenda text-texto-secundario">{t('usuario.pessoaApoio')}</p>
        {(['nenhuma', 'nova', 'existente'] as const).map((modo) => (
          <label key={modo} className="flex items-center gap-3 texto-corpo">
            <input type="radio" value={modo} className="size-5 accent-texto" {...formulario.register('pessoa')} />
            {t(`usuario.pessoaModo.${modo}`)}
          </label>
        ))}
        {modoPessoa === 'existente' && (
          <div className="flex flex-col gap-3 pl-8">
            <BuscaPorTexto rotulo={t('usuario.buscarPessoa')} valor={buscaPessoa} aoBuscar={setBuscaPessoa} />
            <CampoSelecao rotulo={t('usuario.pessoaExistente')} erro={errors.pessoa_id} {...formulario.register('pessoa_id')}>
              <option value="">{t('usuario.escolhaPessoa')}</option>
              {pessoas.itens.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.nome}
                </option>
              ))}
            </CampoSelecao>
            {pessoas.temMais && <p className="texto-legenda text-texto-secundario">{t('usuario.pessoasMaisResultados')}</p>}
          </div>
        )}
      </fieldset>

      <div className="flex justify-end gap-2">
        <Botao onClick={() => void navegar('/administracao/usuarios')}>{t('comum.voltar')}</Botao>
        <Botao type="submit" variante="primario" carregando={criar.isPending}>
          {criar.isPending ? t('usuario.convidando') : t('usuario.convidar')}
        </Botao>
      </div>
    </Formulario>
  );
}
